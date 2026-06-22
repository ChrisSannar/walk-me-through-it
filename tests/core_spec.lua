package.path = "./lua/?.lua;./lua/?/init.lua;" .. package.path
local h = dofile("tests/harness.lua")
local core = require("wmti.core")

h.test("load returns a normalized walkthrough for well-formed JSON", function()
  local json = [[{
    "title": "T",
    "steps": [
      { "file": "a.lua", "line_start": 1, "line_end": 2, "title": "s1", "why": "because" }
    ]
  }]]
  local wt, err = core.load(json)
  h.truthy(wt, "expected a walkthrough, got err: " .. tostring(err))
  h.eq("T", wt.title)
  h.eq(1, #wt.steps)
  h.eq("a.lua", wt.steps[1].file)
end)

h.test("load returns an error for unparseable JSON", function()
  local wt, err = core.load("{not valid json")
  h.falsy(wt, "expected nil walkthrough")
  h.truthy(err, "expected an error message")
end)

h.test("load rejects a step missing a required field, naming the step", function()
  local json = [[{"title":"T","steps":[{"file":"a.lua","line_start":1,"line_end":2,"title":"s1"}]}]]
  local wt, err = core.load(json)
  h.falsy(wt, "expected rejection")
  h.truthy(err and err:find("step 1"), "error should name the step, got: " .. tostring(err))
  h.truthy(err and err:find("why"), "error should name the missing field, got: " .. tostring(err))
end)

h.test("resolve returns ok at the recorded range when the anchor matches", function()
  local step = { line_start = 2, line_end = 3, anchor = "local x = 1" }
  local lines = { "first", "local x = 1", "second", "third" }
  local r = core.resolve(step, lines)
  h.eq({ line_start = 2, line_end = 3, status = "ok" }, r)
end)

h.test("resolve relocates to the exact anchor when lines have drifted", function()
  local step = { line_start = 2, line_end = 3, anchor = "local x = 1" }
  local lines = { "new", "new2", "first", "local x = 1", "second" }
  local r = core.resolve(step, lines)
  h.eq({ line_start = 4, line_end = 5, status = "ok" }, r)
end)

h.test("resolve relocates via whitespace-trimmed match when no exact match exists", function()
  local step = { line_start = 1, line_end = 1, anchor = "local x = 1" }
  local lines = { "other", "    local x = 1" } -- moved down and reindented
  local r = core.resolve(step, lines)
  h.eq({ line_start = 2, line_end = 2, status = "ok" }, r)
end)

h.test("resolve picks the anchor match nearest the original line_start", function()
  local step = { line_start = 5, line_end = 5, anchor = "dup" }
  local lines = { "a", "dup", "b", "c", "d", "dup", "e" } -- matches at 2 and 6
  local r = core.resolve(step, lines)
  h.eq({ line_start = 6, line_end = 6, status = "ok" }, r)
end)

h.test("resolve flags moved and falls back to the recorded line when the anchor is gone", function()
  local step = { line_start = 3, line_end = 4, anchor = "gone forever" }
  local lines = { "a", "b", "c", "d", "e" }
  local r = core.resolve(step, lines)
  h.eq({ line_start = 3, line_end = 4, status = "moved" }, r)
end)

h.test("resolve flags missing when the file's lines are absent", function()
  local step = { line_start = 1, line_end = 2, anchor = "x" }
  local r = core.resolve(step, nil)
  h.eq({ line_start = 1, line_end = 2, status = "missing" }, r)
end)

h.test("resolve uses the raw range with ok when the step has no anchor", function()
  local step = { line_start = 2, line_end = 4 }
  local lines = { "a", "b", "c", "d", "e" }
  local r = core.resolve(step, lines)
  h.eq({ line_start = 2, line_end = 4, status = "ok" }, r)
end)

h.test("staleness is true only when base and head are present and differ", function()
  h.falsy(core.staleness("abc", "abc"), "same sha is not stale")
  h.truthy(core.staleness("abc", "def"), "different sha is stale")
  h.falsy(core.staleness(nil, "def"), "absent base is not stale")
  h.falsy(core.staleness("abc", nil), "absent head is not stale")
end)

h.run()
