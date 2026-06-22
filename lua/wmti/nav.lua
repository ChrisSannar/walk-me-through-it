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

  local lines = nil
  local readable = vim.fn.filereadable(full) == 1
  if readable and s.code_win and vim.api.nvim_win_is_valid(s.code_win) then
    vim.api.nvim_set_current_win(s.code_win)
    vim.cmd("edit " .. vim.fn.fnameescape(full))
    s.code_win = vim.api.nvim_get_current_win()
    local cbuf = vim.api.nvim_win_get_buf(s.code_win)
    lines = vim.api.nvim_buf_get_lines(cbuf, 0, -1, false)
  end

  local r = core.resolve(step, lines)

  if readable and s.code_win and vim.api.nvim_win_is_valid(s.code_win) then
    local cbuf = vim.api.nvim_win_get_buf(s.code_win)
    local total = vim.api.nvim_buf_line_count(cbuf)
    local target = math.max(1, math.min(r.line_start, total))
    vim.api.nvim_win_set_cursor(s.code_win, { target, 0 })
    vim.api.nvim_win_call(s.code_win, function()
      vim.cmd("normal! zz")
    end)
    highlight_range(s, cbuf, r.line_start, r.line_end)
  end

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
    code_win = vim.api.nvim_get_current_win(),
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
    vim.notify("wmti: end of walkthrough", vim.log.levels.INFO)
    return
  end
  M.goto_step(s.index + 1)
end

function M.prev()
  local s = M.session
  if not s then
    return
  end
  if s.index <= 1 then
    vim.notify("wmti: start of walkthrough", vim.log.levels.INFO)
    return
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
