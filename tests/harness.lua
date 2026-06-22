-- Minimal zero-dependency test harness.
-- Run specs with: nvim --headless -l tests/<name>_spec.lua
local M = { tests = {}, passed = 0, failed = 0 }

function M.test(name, fn)
  table.insert(M.tests, { name = name, fn = fn })
end

function M.eq(expected, actual, msg)
  if not vim.deep_equal(expected, actual) then
    error(
      (msg or "values not equal")
        .. "\n  expected: "
        .. vim.inspect(expected)
        .. "\n  actual:   "
        .. vim.inspect(actual),
      2
    )
  end
end

function M.truthy(v, msg)
  if not v then
    error((msg or "expected truthy") .. ", got: " .. vim.inspect(v), 2)
  end
end

function M.falsy(v, msg)
  if v then
    error((msg or "expected falsy") .. ", got: " .. vim.inspect(v), 2)
  end
end

function M.run()
  for _, t in ipairs(M.tests) do
    local ok, err = pcall(t.fn)
    if ok then
      M.passed = M.passed + 1
      print("ok   - " .. t.name)
    else
      M.failed = M.failed + 1
      print("FAIL - " .. t.name)
      print("       " .. tostring(err):gsub("\n", "\n       "))
    end
  end
  print(string.format("\n%d passed, %d failed", M.passed, M.failed))
  if M.failed > 0 then
    os.exit(1)
  end
end

return M
