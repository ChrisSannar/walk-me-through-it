-- Headless smoke test for the Neovim shell. Run from repo root:
--   nvim --headless -l tests/shell_smoke.lua
package.path = "./lua/?.lua;./lua/?/init.lua;" .. package.path

local function check(cond, msg)
  if not cond then
    print("FAIL - " .. msg)
    os.exit(1)
  end
  print("ok   - " .. msg)
end

-- A real file in the current window stands in for the "code window".
vim.cmd("edit lua/wmti/core.lua")

local wmti = require("wmti")
wmti.setup({})
wmti.open() -- discovers .wmti/sample.json and auto-opens (single walkthrough)

local nav = require("wmti.nav")
check(nav.session ~= nil, "session created by open")
check(nav.session.side_win and vim.api.nvim_win_is_valid(nav.session.side_win), "sidebar window is valid")
check(nav.session.index == 1, "lands on step 1")

local c1 = vim.api.nvim_win_get_cursor(nav.session.code_win)[1]
nav.next()
check(nav.session.index == 2, "next advances the index")
local c2 = vim.api.nvim_win_get_cursor(nav.session.code_win)[1]
check(c1 ~= c2, "code window cursor moved between steps (" .. c1 .. " -> " .. c2 .. ")")

local sb = vim.api.nvim_buf_get_lines(nav.session.side_buf, 0, -1, false)
check(sb[1] == "wmti walks itself", "sidebar shows the walkthrough title")

-- Step 4 anchors on core.resolve; confirm it resolves onto that function line.
nav.goto_step(4)
local cbuf = vim.api.nvim_win_get_buf(nav.session.code_win)
local at = vim.api.nvim_win_get_cursor(nav.session.code_win)[1]
local line = vim.api.nvim_buf_get_lines(cbuf, at - 1, at, false)[1]
check(line and line:find("function M.resolve", 1, true) ~= nil, "step 4 lands on core.resolve (line " .. at .. ")")

-- Prev/next bounds: from step 1, prev stays at 1.
nav.goto_step(1)
nav.prev()
check(nav.session.index == 1, "prev stops at the start")

wmti.close()
check(nav.session == nil, "close tears down the session")

print("\nSMOKE OK")
