import atomics, json, ws

var atomic_variable*: Atomic[ptr string]
var global_values*: JsonNode
var temp_WebSocket*: seq[WebSocket] = @[WebSocket()]