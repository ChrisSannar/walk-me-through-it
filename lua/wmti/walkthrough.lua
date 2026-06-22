-- wmti walkthrough: discovery + file load. The I/O boundary feeding core.load.
local M = {}

local core = require("wmti.core")
local config = require("wmti.config")

-- Absolute paths of every <dir>/*.json under root, sorted.
function M.discover(root)
  local pattern = root .. "/" .. config.options.dir .. "/*.json"
  local files = vim.fn.glob(pattern, false, true)
  table.sort(files)
  return files
end

-- Read a walkthrough JSON file and parse it through core.load.
function M.load_path(path)
  local f = io.open(path, "r")
  if not f then
    return nil, "cannot read " .. path
  end
  local content = f:read("*a")
  f:close()
  return core.load(content)
end

return M
