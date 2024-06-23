import os
import json, ws
import strutils, strformat, json
from osproc import execCmdEx
import osproc, streams, strutils, sequtils, asyncdispatch, atomics
import /melone_player/types

# proc add*[T](arry: ref seq[T], element: T) =


# proc add*[T](atomicArray: var Atomic[ref seq[T]], element: T) =
#   var arrayReddy = atomicArray.load[]
#   arrayReddy.add(element)
#   var temp_array = cast[ref seq[T]](addr arrayReddy)
#   atomicArray.store(temp_array)


# proc remove*[T](atomicArray: var Atomic[ref seq[T]], element: T) =
#   var arrayReddy = atomicArray.load[]
#   for index, data in arrayReddy:
#     if data == element:
#       arrayReddy.delete(index)
#   var temp_array = cast[ref seq[T]](addr arrayReddy)
#   atomicArray.store(temp_array)



proc runCommandInRealTime*(command: string = "", args: seq[string] = @[]): Future[void] {.async.} =

  try:
    var process = startProcess(command, args=args, options = {poUsePath, poStdErrToStdOut})
    var outp = process.outputStream
    var line = newStringOfCap(120)

    # Проверяем, что процесс успешно запущен
    if not process.running:
      echo "Failed to start process"

    # Чтение вывода процесса в реальном времени
    while process.running:
      if outp.readLine(line):
        echo "ytdl.exe: " & line
        if line.contains("[download]") and not line.contains("Destination:"):
          var rs = $line.replace("[download]", "").replace("of", "").replace("at", "").replace("ETA", "").replace("~", "").replace("frag", "").replace("(","").replace(")", "").split(" ").filterIt(it.len > 0)
          var rs_ptr = cast[ptr string](addr rs)
          atomic_variable.store(rs_ptr)
        elif line.contains("Destination:"):
          var rs = $line.replace("[download] Destination:", "Download:").split(" ")
          var rs_ptr = cast[ptr string](addr rs)
          atomic_variable.store(rs_ptr)
    echo "fertisch"
    var rs = ""
    var rs_ptr = cast[ptr string](addr rs)
    atomic_variable.store(rs_ptr)
    # Ожидание завершения процесса
    process.close()
  except CatchableError as e:
    echo e.msg 




# Функция для получения содержимого папки
proc getDirectoryContents*(path: string): JsonNode =
  # Проверяем, существует ли путь и является ли он директорией
  if not dirExists(path):
    raise newException(ValueError, "Указанный путь не существует или не является директорией")

  # Создаем пустой список для хранения результатов
  var contents = %*{}

  # Итерируемся по содержимому директории
  for items in walkDir(path):
    if dirExists(items.path):
      echo splitPath(items.path)
      contents[items.path] = %*[]
      for item in walkDir(items.path):
        contents[items.path].add(%*item.path) # Добавляем путь к каждому элементу в список
    else:
      contents[items.path].add(%*items.path)

  return contents

proc isFileData*(path: string = "./"): bool =
  for kind, path in walkDir(path):
    if kind == pcFile:
      if path.endsWith("settings.json"):
        return true
  return false

proc getData*(path: string = "./", key: seq[string] = @["app_info"]): JsonNode =
  if isFileData(path):
    try:
      var currentNode = parseFile(fmt"{path}settings.json")
      for field in key:
        if currentNode.kind == JObject:
          currentNode = currentNode[field]
        else:
          raise newException(ValueError, fmt"{key} not found!")
      return currentNode

    except IOError:
      return %*{"error": "config_file_not_found", "error_number": 404}
    except Exception as Except:
      return %*{"error": fmt"{Except.msg}", "error_number": 520}
  else:
    return %*{"error": "config_file_not_found", "error_number": 404}

type
  BracketedText* = tuple
    index: int
    text: string

proc findTextInBrackets*(inputStr: string): seq[BracketedText] =
  var result: seq[BracketedText] = @[]
  var startIndex = 0

  while true:
    let openIndex = inputStr.find('(', startIndex)
    if openIndex == -1:
      break # Если скобка '(' не найдена, завершаем цикл
    startIndex = openIndex + 1 # Начинаем поиск после найденной скобки '('

    let closeIndex = inputStr.find(')', startIndex)
    if closeIndex == -1:
      break # Если закрывающая скобка ')' не найдена, завершаем цикл

    let startIdx = openIndex
    let textInsideBrackets = inputStr[openIndex+1..<closeIndex]
    result.add((startIdx, textInsideBrackets))

    startIndex = closeIndex + 1 # Начинаем поиск после найденной скобки ')'

  return result


proc jsonNodeToSeq(node: JsonNode): seq[string] =
  if node.kind != JArray:
    raise newException(ValueError, "Expected a JSON array")
  result = @[]
  for item in node:
    if item.kind == JString:
      result.add(item.getStr)
    else:
      raise newException(ValueError, "Expected all elements to be strings")


# Функция для обхода ключей и значений в JsonNode и возвращения массива путей
proc traverseDirectories(node: JsonNode, path: string = "", paths: var seq[string]) =
  case node.kind
  of JObject:
    for k, v in node:
      let newPath = if path.len > 0: path & k else: k
      paths.add(newPath)
      traverseDirectories(v, newPath, paths)
  of JArray:
    for item in node:
      traverseDirectories(item, path, paths)
  else:
    let newPath = if path.len > 0: path & node.str else: node.str
    paths.add(newPath)




# Начать обход с узла "paths"
proc getPaths*(node: JsonNode, key:string = ""): seq[string] =
  var paths: seq[string]
  if node.kind == JObject and node.hasKey($key):
    traverseDirectories(node[$key], "", paths)
  else:
    echo "JSON does not have 'paths' key or is not an object."
  return paths


proc ifAllPath*(): bool=
  var data = getPaths(getData(key = @["app_settings"]), "paths")
  for dirPath in data:
    if not dirExists(dirPath):
      try:
        createDir(dirPath)
      except IOError:
        echo "Не вдалось створити папку: ", dirPath
        return false
    else:
      discard
  return true






proc download_youtube_music*(url:seq[string]): Future[void] {.async.} =
  var command = jsonNodeToSeq(getData(key = @["lib_settings", "exe","ytdl.exe","download_command", "audio"]))
  if ifAllPath() == true:
    var data = getPaths(getData(key = @["app_settings"]), "paths")
    echo data[7]
    command[0] = $data[2] & "/ytdl.exe"
    command[2] = getAppDir() & $data[5].replace(".","").replace("/","\\") & "\\%(title)s.%(ext)s"
    for i1 in url:
      command[6] = i1
      echo command[1..command.len-1]
      await runCommandInRealTime(command=command[0], args=command[1..command.len-1])
