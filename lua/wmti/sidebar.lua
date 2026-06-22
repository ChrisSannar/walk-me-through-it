-- wmti sidebar: the right (or left) split that lists steps and shows the
-- current step's `why`. Pure rendering + buffer-local maps; no resolve logic.
local M = {}

local config = require("wmti.config")

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
  local lines, cur_line = {}, nil

  table.insert(lines, s.wt.title or "Walkthrough")
  table.insert(lines, rule)
  if s.stale then
    table.insert(lines, "⚠ authored against older code")
  end
  table.insert(lines, "")

  for i, step in ipairs(s.wt.steps) do
    local marker = (i == s.index) and "▶ " or "  "
    table.insert(lines, marker .. i .. ". " .. step.title)
    if i == s.index then
      cur_line = #lines - 1
      local loc = (resolved and resolved.line_start) or step.line_start
      table.insert(lines, "    " .. step.file .. ":" .. loc)
      if resolved and resolved.status == "moved" then
        table.insert(lines, "    ⚠ may have moved")
      elseif resolved and resolved.status == "missing" then
        table.insert(lines, "    ⚠ file missing")
      end
      table.insert(lines, "")
      for _, wl in ipairs(wrap(step.why, width - 4)) do
        table.insert(lines, "  " .. wl)
      end
      table.insert(lines, "")
    end
  end

  table.insert(lines, rule)
  table.insert(lines, "n next · p prev · ⏎ dive · q close")

  vim.bo[buf].modifiable = true
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
  vim.bo[buf].modifiable = false

  local ns = s.side_ns
  vim.api.nvim_buf_clear_namespace(buf, ns, 0, -1)
  vim.api.nvim_buf_set_extmark(buf, ns, 0, 0, { line_hl_group = "Title", hl_eol = true })
  if cur_line then
    vim.api.nvim_buf_set_extmark(buf, ns, cur_line, 0, { line_hl_group = "WmtiCurrent", hl_eol = true })
    if s.side_win and vim.api.nvim_win_is_valid(s.side_win) then
      pcall(vim.api.nvim_win_set_cursor, s.side_win, { cur_line + 1, 0 })
    end
  end
end

return M
