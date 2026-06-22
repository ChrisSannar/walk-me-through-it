-- wmti sidebar: the right (or left) split that lists steps and shows the
-- current step's `why`. Pure rendering + buffer-local maps; no resolve logic.
local M = {}

local config = require("wmti.config")

-- Color coding for the explanation window. All link to standard groups with
-- default=true, so a user's colorscheme / overrides win.
local function ensure_highlights()
  local links = {
    WmtiTitle = "Title",
    WmtiRule = "NonText",
    WmtiStale = "WarningMsg",
    WmtiStepCurrent = "Title",
    WmtiStepIdle = "Comment",
    WmtiLocation = "Directory",
    WmtiFlag = "WarningMsg",
    WmtiWhy = "Normal",
    WmtiHelp = "Comment",
    WmtiCurrentBg = "CursorLine",
  }
  for group, target in pairs(links) do
    vim.api.nvim_set_hl(0, group, { link = target, default = true })
  end
end

-- Greedy word-wrap to `width` columns.
local function wrap(text, width)
  local out, line = {}, ""
  for word in tostring(text):gmatch("%S+") do
    if #line > 0 and #line + 1 + #word > width then
      table.insert(out, line)
      line = word
    else
      line = (#line == 0) and word or (line .. " " .. word)
    end
  end
  if #line > 0 then
    table.insert(out, line)
  end
  return out
end

local function map(buf, lhs, fn)
  vim.keymap.set("n", lhs, fn, { buffer = buf, nowait = true, silent = true })
end

-- Open the sidebar split, storing window/buffer/namespace on the session.
function M.open(s)
  ensure_highlights()

  local buf = vim.api.nvim_create_buf(false, true)
  vim.bo[buf].bufhidden = "wipe"
  vim.bo[buf].buftype = "nofile"
  vim.bo[buf].swapfile = false
  vim.bo[buf].filetype = "wmti"

  vim.cmd(config.options.side == "left" and "topleft vsplit" or "botright vsplit")
  local win = vim.api.nvim_get_current_win()
  vim.api.nvim_win_set_buf(win, buf)
  vim.api.nvim_win_set_width(win, config.options.width)

  vim.wo[win].number = false
  vim.wo[win].relativenumber = false
  vim.wo[win].wrap = false
  vim.wo[win].cursorline = false
  vim.wo[win].signcolumn = "no"
  vim.wo[win].winfixwidth = true
  vim.wo[win].list = false

  s.side_win = win
  s.side_buf = buf
  s.side_ns = vim.api.nvim_create_namespace("wmti_side")

  local nav = require("wmti.nav")
  map(buf, "n", function() nav.next() end)
  map(buf, "p", function() nav.prev() end)
  map(buf, "<Tab>", function() nav.next() end)
  map(buf, "<S-Tab>", function() nav.prev() end)
  map(buf, "<CR>", function() nav.dive() end)
  map(buf, "q", function() nav.close() end)
end

-- Re-render the whole sidebar for the current step. `resolved` is the
-- core.resolve result for the current step (for the location + status flag).
function M.render(s, resolved)
  local buf = s.side_buf
  if not (buf and vim.api.nvim_buf_is_valid(buf)) then
    return
  end
  local width = config.options.width
  local rule = string.rep("─", math.max(1, width - 2))

  -- rows: { text, group?, bg? }. One highlight group per line keeps multibyte
  -- markers (▶ ─ ⚠ ⏎) off the column math.
  local rows, cur_line = {}, nil
  local function add(text, group, bg)
    rows[#rows + 1] = { text = text, group = group, bg = bg }
  end

  add(s.wt.title or "Walkthrough", "WmtiTitle")
  add(rule, "WmtiRule")
  if s.stale then
    add("⚠ authored against older code", "WmtiStale")
  end
  add("")

  for i, step in ipairs(s.wt.steps) do
    local current = (i == s.index)
    local marker = current and "▶ " or "  "
    add(marker .. i .. ". " .. step.title, current and "WmtiStepCurrent" or "WmtiStepIdle", current and "WmtiCurrentBg")
    if current then
      cur_line = #rows - 1
      local loc = (resolved and resolved.line_start) or step.line_start
      add("    " .. step.file .. ":" .. loc, "WmtiLocation")
      if resolved and resolved.status == "moved" then
        add("    ⚠ may have moved", "WmtiFlag")
      elseif resolved and resolved.status == "missing" then
        add("    ⚠ file missing", "WmtiFlag")
      end
      add("")
      for _, wl in ipairs(wrap(step.why, width - 4)) do
        add("  " .. wl, "WmtiWhy")
      end
      add("")
    end
  end

  add(rule, "WmtiRule")
  add("n next · p prev · ⏎ dive · q close", "WmtiHelp")

  local text = {}
  for _, r in ipairs(rows) do
    text[#text + 1] = r.text
  end

  vim.bo[buf].modifiable = true
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, text)
  vim.bo[buf].modifiable = false

  local ns = s.side_ns
  vim.api.nvim_buf_clear_namespace(buf, ns, 0, -1)
  for i, r in ipairs(rows) do
    local line0 = i - 1
    if r.bg then
      vim.api.nvim_buf_set_extmark(buf, ns, line0, 0, { line_hl_group = r.bg, hl_eol = true })
    end
    if r.group then
      pcall(vim.api.nvim_buf_add_highlight, buf, ns, r.group, line0, 0, -1)
    end
  end

  if cur_line and s.side_win and vim.api.nvim_win_is_valid(s.side_win) then
    pcall(vim.api.nvim_win_set_cursor, s.side_win, { cur_line + 1, 0 })
  end
end

return M
