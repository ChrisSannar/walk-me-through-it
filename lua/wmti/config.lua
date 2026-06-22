-- wmti config: defaults + setup merge. No I/O, no Neovim UI.
local M = {}

M.defaults = {
  side = "right", -- "right" | "left"
  width = 44,
  dir = ".wmti", -- relative to cwd; walkthroughs are <dir>/*.json
}

M.options = vim.deepcopy(M.defaults)

function M.setup(opts)
  M.options = vim.tbl_deep_extend("force", vim.deepcopy(M.defaults), opts or {})
end

return M
