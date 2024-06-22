# # # example.nim
# # import htmlgen
# # import jester

# # routes:
# #   get "/":
# #     resp h1("Hello world")
# #   get "/hello/@name?@forname?":
# #     if @"name" == "":
# #       resp "No name received :("
# #     else:
# #       resp "Hello " & @"name" & @"forname"


# # import nimclipboard/libclipboard

# # var cb = clipboard_new(nil)
# # cb.clipboard_set_text("nimclipboard rocks")
# import strutils

# type
#   BracketedText* = tuple
#     index: int
#     text: string

# proc findTextInBrackets*(inputStr: string): seq[BracketedText] =
#   var result: seq[BracketedText] = @[]
#   var startIndex = 0

#   while true:
#     let openIndex = inputStr.find('(', startIndex)
#     if openIndex == -1:
#       break # Если скобка '(' не найдена, завершаем цикл
#     startIndex = openIndex + 1 # Начинаем поиск после найденной скобки '('

#     let closeIndex = inputStr.find(')', startIndex)
#     if closeIndex == -1:
#       break # Если закрывающая скобка ')' не найдена, завершаем цикл

#     let startIdx = openIndex
#     let textInsideBrackets = inputStr[openIndex+1..<closeIndex]
#     result.add((startIdx, textInsideBrackets))

#     startIndex = closeIndex + 1 # Начинаем поиск после найденной скобки ')'

#   return result

# # Пример использования функции
# let inputString = "Hello (world), (this) is a (test) string"
# let textsInBrackets = findTextInBrackets(inputString)

# # Выводим найденные тексты в скобках с номерами символов, с которых начинались скобки
# for (index, text) in textsInBrackets:
#   echo "Text:", text
#   echo "Start Index:", index

# This example shows the basic use of the NiGui toolkit.


# import osproc, streams, strutils, sequtils

# proc runCommandInRealTime(command: string, args: openArray[string]) =
#   var process = startProcess(command, args=args, options = {poUsePath, poStdErrToStdOut})
#   var outp = process.outputStream
#   var line = newStringOfCap(120)

#   # Проверяем, что процесс успешно запущен
#   if not process.running:
#     echo "Failed to start process"
#     return

#   # Чтение вывода процесса в реальном времени
#   while process.running:
#     if outp.readLine(line):
#       if line.contains("[download]") and not line.contains("Destination:"):
#         echo $line.replace("[download]", "").replace("of", "").replace("at", "").replace("ETA", "").split(" ").filterIt(it.len > 0)
#       elif line.contains("Destination:"):
#         echo $line.replace("[download] Destination:", "Download:")

#   # Ожидание завершения процесса
#   process.close()

# # Пример использования:
# let command = "ytdl.exe"
# runCommandInRealTime(command, ["-o", "%(title)s.%(ext)s", "--newline", "--format", "best", 
#                                "https://www.youtube.com/watch?v=OVh0bMNSFss"])


import os
import json
import strutils, strformat, json
from osproc import execCmdEx
import osproc, streams, strutils, sequtils, asyncdispatch, atomics

proc runCommandInRealTime*() {.nimcall, gcsafe.} =
  try:
    var process = startProcess("ytdl.exe", args=["-o",
              "{path}%(title)s.%(ext)s",
              "--newline",
              "--audio-format",
              "mp3", "https://youtu.be/hc4rbyXmXcQ?si=_19ttY8raXj93d1c"], options = {poUsePath, poStdErrToStdOut})
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
        elif line.contains("Destination:"):
          var rs = $line.replace("[download] Destination:", "Download:").split(" ")
    echo "fertisch"
    # Ожидание завершения процесса
    process.close()
  except CatchableError as e:
    echo e.msg 
runCommandInRealTime()




# import mummy, mummy/routers
# import nigui, strformat
# import nigui/msgbox
# app.init()
# proc echoError(error: string = "error", typeError: string = "warn") =
#   var window = newWindow()
#   var listButtons = @["OK", "Exit", "Copy Error"]
#   var res = window.msgBox(message = "Error:\n\n" & $error, title = "Melon Player Error:", button1=fmt"{listButtons[0]}")
#   window.show()


# proc indexHandler(request: Request) =
#   {.gcsafe.}:
#     echoError()
#   var headers: HttpHeaders
#   headers["Content-Type"] = "text/plain"
#   request.respond(200, headers, "Hello, World!")


# var router: Router
# router.get("/", indexHandler)

# let server = newServer(router)
# echo "Serving on http://localhost:8080"
# server.serve(Port(8080))