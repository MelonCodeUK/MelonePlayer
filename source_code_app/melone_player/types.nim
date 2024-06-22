import atomics, json

var atomic_variable*: Atomic[ptr string]
var global_values*: JsonNode