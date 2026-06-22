-- wmti nav: the session orchestrator. Holds the open walkthrough, drives the
-- code window via core.resolve, and keeps the sidebar in sync. All Neovim
-- window/git I/O lives here; the pure decisions live in core.
local M = {}

local core = require("wmti.core")
local sidebar = require("wmti.sidebar")
local walkthrough = require("wmti.walkthrough")

M.session = nil

local function ensure_highlights()
  vim.api.nvim_set_hl(0, "WmtiStep", { link = "Visual", default = true })
  vim.api.nvim_set_hl(0, "WmtiCurrent", { link = "PmenuSel", default = true })
end

-- A "normal" editing window: not floating, showing a real file buffer
-- (buftype ""). Excludes neo-tree, dashboards, terminals, our sidebar, etc.
local function is_normal_win(win)
  if not vim.api.nvim_win_is_valid(win) then
    return false
  end
  if vim.api.nvim_win_get_config(win).relative ~= "" then
    return false -- floating
  end
  local buf = vim.api.nvim_win_get_buf(win)
  return vim.bo[buf].buftype == ""
end

-- Pick the window the code should jump into: the current one if it's a normal
-- editing window, else the first normal window on the tab, else the current.
local function pick_code_win()
  local cur = vim.api.nvim_get_current_win()
  if is_normal_win(cur) then
    return cur
  end
  for _, w in ipairs(vim.api.nvim_tabpage_list_wins(0)) do
    if is_normal_win(w) then
      return w
    end
  end
  return cur
end

-- HEAD sha of the repo at `root`, or nil when not a git repo / git fails.
local function head_sha(root)
  local ok, res = pcall(function()
    return vim.system({ "git", "-C", root, "rev-parse", "HEAD" }, { text = true }):wait()
  end)
  if not ok or res.code ~= 0 then
    return nil
  end
  return vim.trim(res.stdout)
end

local function highlight_range(s, buf, ls, le)
  if s.hl_buf and vim.api.nvim_buf_is_valid(s.hl_buf) then
    vim.api.nvim_buf_clear_namespace(s.hl_buf, s.ns, 0, -1)
  end
  vim.api.nvim_buf_clear_namespace(buf, s.ns, 0, -1)
  local last = vim.api.nvim_buf_line_count(buf)
  for l = ls, math.min(le, last) do
    if l >= 1 then
      vim.api.nvim_buf_set_extmark(buf, s.ns, l - 1, 0, { line_hl_group = "WmtiStep", hl_eol = true })
    end
  end
  s.hl_buf = buf
end

-- Jump to step `index` (clamped): load its file into the code window, resolve
-- its real location, center + highlight it, then re-render the sidebar.
-- Focus ends in the sidebar (sidebar-driven preview).
function M.goto_step(index)
  local s = M.session
  if not s then
    return
  end
  index = math.max(1, math.min(index, #s.wt.steps))
  s.index = index
  local step = s.wt.steps[index]
  local full = s.root .. "/" .. step.file

  -- Load the file's buffer and place it directly into the code window. We use
  -- nvim_win_set_buf (not :edit) so 'switchbuf' and autocmds can't redirect the
  -- code into some other / extra window — the jump stays in the main window.
  local lines, cbuf = nil, nil
  local readable = vim.fn.filereadable(full) == 1
  if readable and s.code_win and vim.api.nvim_win_is_valid(s.code_win) then
    -- Guard bufload/set_buf: a swap file or read error must not crash the
    -- viewer. Suppress the swap-attention prompt ('A') during the load — this
    -- is a read-only viewer, so an existing swap (same file open elsewhere)
    -- shouldn't interrupt — then restore shortmess.
    local save_sm = vim.o.shortmess
    if not save_sm:find("A", 1, true) then
      vim.o.shortmess = save_sm .. "A"
    end
    local ok, buf = pcall(function()
      local b = vim.fn.bufadd(full)
      vim.fn.bufload(b)
      vim.bo[b].buflisted = true
      vim.api.nvim_win_set_buf(s.code_win, b)
      return b
    end)
    vim.o.shortmess = save_sm
    if ok then
      cbuf = buf
      lines = vim.api.nvim_buf_get_lines(cbuf, 0, -1, false)
    end
  end

  local r = core.resolve(step, lines)

  if cbuf then
    local total = vim.api.nvim_buf_line_count(cbuf)
    local target = math.max(1, math.min(r.line_start, total))
    vim.api.nvim_win_set_cursor(s.code_win, { target, 0 })
    vim.api.nvim_win_call(s.code_win, function()
      vim.cmd("normal! zz")
    end)
    highlight_range(s, cbuf, r.line_start, r.line_end)
  end

  -- Focus stays in the sidebar (we never moved it) so stepping keeps working.
  if s.side_win and vim.api.nvim_win_is_valid(s.side_win) then
    vim.api.nvim_set_current_win(s.side_win)
  end
  sidebar.render(s, r)
end

-- Open a walkthrough from a JSON path: parse, set up the session + sidebar,
-- and land on step 1. Replaces any currently open session.
function M.open(path)
  if M.session then
    M.close()
  end
  local wt, err = walkthrough.load_path(path)
  if not wt then
    vim.notify("wmti: " .. tostring(err), vim.log.levels.ERROR)
    return
  end
  ensure_highlights()
  local root = vim.fn.getcwd()
  local s = {
    path = path,
    wt = wt,
    index = 1,
    code_win = pick_code_win(),
    root = root,
    ns = vim.api.nvim_create_namespace("wmti_code"),
    stale = false,
  }
  if wt.base_sha then
    s.stale = core.staleness(wt.base_sha, head_sha(root))
  end
  M.session = s
  sidebar.open(s)
  M.goto_step(1)
end

function M.next()
  local s = M.session
  if not s then
    return
  end
  if s.index >= #s.wt.steps then
    return -- at the end; the sidebar already shows where you are
  end
  M.goto_step(s.index + 1)
end

function M.prev()
  local s = M.session
  if not s then
    return
  end
  if s.index <= 1 then
    return -- at the start; the sidebar already shows where you are
  end
  M.goto_step(s.index - 1)
end

-- Move focus into the code window at the current step.
function M.dive()
  local s = M.session
  if not s then
    return
  end
  if s.code_win and vim.api.nvim_win_is_valid(s.code_win) then
    vim.api.nvim_set_current_win(s.code_win)
  end
end

function M.close()
  local s = M.session
  if not s then
    return
  end
  if s.hl_buf and vim.api.nvim_buf_is_valid(s.hl_buf) then
    vim.api.nvim_buf_clear_namespace(s.hl_buf, s.ns, 0, -1)
  end
  if s.side_win and vim.api.nvim_win_is_valid(s.side_win) then
    vim.api.nvim_win_close(s.side_win, true)
  end
  M.session = nil
end

return M
