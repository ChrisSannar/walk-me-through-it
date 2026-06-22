-- wmti command registration. Lazy entry point: a single :Wmti command that
-- dispatches to subcommands (open is the default). Native Vim command names
-- must be uppercase + alphanumeric, so this is the idiomatic stand-in for a
-- lowercase `wmti` / `wmti-<sub>` CLI feel (see the lowercase abbrev in setup).
if vim.g.loaded_wmti then
  return
end
vim.g.loaded_wmti = true

-- subcommand -> require("wmti") function name
local subcommands = {
  open = "open",
  next = "next",
  prev = "prev",
  close = "close",
  select = "select",
  dive = "dive",
}

vim.api.nvim_create_user_command("Wmti", function(opts)
  local sub = opts.fargs[1] or "open"
  local fn = subcommands[sub]
  if not fn then
    vim.notify("wmti: unknown subcommand '" .. sub .. "'", vim.log.levels.ERROR)
    return
  end
  require("wmti")[fn]()
end, {
  nargs = "?",
  desc = "wmti walkthrough viewer (open|next|prev|close|select|dive)",
  complete = function(arglead)
    local names = vim.tbl_keys(subcommands)
    table.sort(names)
    return vim.tbl_filter(function(n)
      return n:find(arglead, 1, true) == 1
    end, names)
  end,
})
