-- wmti public entry: setup() + the command-facing open/select orchestration.
local M = {}

local config = require("wmti.config")

function M.setup(opts)
  config.setup(opts)
end

-- :WmtiOpen — 0 walkthroughs notify, 1 auto-open, many show the picker.
function M.open()
  local paths = require("wmti.walkthrough").discover(vim.fn.getcwd())
  if #paths == 0 then
    vim.notify("wmti: no walkthroughs found in " .. config.options.dir .. "/", vim.log.levels.WARN)
    return
  end
  if #paths == 1 then
    require("wmti.nav").open(paths[1])
    return
  end
  M.select()
end

-- :WmtiSelect — always prompt, even when only one exists.
function M.select()
  local paths = require("wmti.walkthrough").discover(vim.fn.getcwd())
  if #paths == 0 then
    vim.notify("wmti: no walkthroughs found in " .. config.options.dir .. "/", vim.log.levels.WARN)
    return
  end
  vim.ui.select(paths, {
    prompt = "Select walkthrough",
    format_item = function(p)
      return vim.fn.fnamemodify(p, ":t")
    end,
  }, function(choice)
    if choice then
      require("wmti.nav").open(choice)
    end
  end)
end

function M.close()
  require("wmti.nav").close()
end

function M.next()
  require("wmti.nav").next()
end

function M.prev()
  require("wmti.nav").prev()
end

function M.dive()
  require("wmti.nav").dive()
end

return M
