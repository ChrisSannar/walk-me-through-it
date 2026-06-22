-- wmti core: pure load / resolve / staleness. No Neovim UI, git, or filesystem.
local M = {}

local STEP_REQUIRED = { "file", "line_start", "line_end", "title", "why" }

function M.load(json_string)
  local ok, decoded = pcall(vim.json.decode, json_string)
  if not ok then
    return nil, "invalid JSON: " .. tostring(decoded)
  end
  for i, step in ipairs(decoded.steps or {}) do
    for _, field in ipairs(STEP_REQUIRED) do
      if step[field] == nil then
        return nil, string.format("step %d: missing required field '%s'", i, field)
      end
    end
  end
  return decoded
end

local function trim(s)
  return (s:gsub("^%s+", ""):gsub("%s+$", ""))
end

-- Index of the line satisfying `pred` whose position is closest to `origin`.
-- Ties keep the earlier (lower) index.
local function find_nearest(lines, pred, origin)
  local best, best_dist
  for i, line in ipairs(lines) do
    if pred(line) then
      local dist = math.abs(i - origin)
      if not best_dist or dist < best_dist then
        best, best_dist = i, dist
      end
    end
  end
  return best
end

local function located(line_start, line_end, status)
  return { line_start = line_start, line_end = line_end, status = status }
end

function M.resolve(step, lines)
  if lines == nil then
    return located(step.line_start, step.line_end, "missing")
  end
  local len = step.line_end - step.line_start + 1
  local origin = step.line_start
  local found
  if step.anchor then
    local target = trim(step.anchor)
    found = find_nearest(lines, function(l)
      return l == step.anchor
    end, origin) or find_nearest(lines, function(l)
      return trim(l) == target
    end, origin)
  end
  if found then
    return located(found, found + len - 1, "ok")
  end
  return located(origin, step.line_end, step.anchor and "moved" or "ok")
end

function M.staleness(base_sha, head_sha)
  if not base_sha or not head_sha then
    return false
  end
  return base_sha ~= head_sha
end

return M
