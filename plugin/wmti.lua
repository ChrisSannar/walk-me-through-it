-- wmti command registration. Lazy entry point: defines the :Wmti* commands.
if vim.g.loaded_wmti then
  return
end
vim.g.loaded_wmti = true

local function cmd(name, fn)
  vim.api.nvim_create_user_command(name, function()
    require("wmti")[fn]()
  end, {})
end

cmd("WmtiOpen", "open")
cmd("WmtiSelect", "select")
cmd("WmtiClose", "close")
cmd("WmtiNext", "next")
cmd("WmtiPrev", "prev")
