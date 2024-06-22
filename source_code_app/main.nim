
import jester
import melonGui/melongui
import functions
import strutils, crowngui
import crowngui/webview
import strformat, asyncdispatch
import ws, ws/jester_extra
import winim
import std/os
import jester/[request]
import json, nigui
import nigui/msgbox
import std/[logging]
import colored_logger, locks, asynchttpserver, threadpool, atomics
import nimclipboard/libclipboard
import melone_player/types
global_values = %*{}
var initialStr = ""
let initialStrPtr = cast[ptr string](addr initialStr)
atomic_variable.store(initialStrPtr)
app.init()

var channel: Channel[seq[string]]
channel.open()
var s = 0

var lock: Lock
initLock lock


const  appName: string = "MelonPlayer"

var
  version: string
  settings_version: string
  width: int
  height: int
  port: int
  urlServer: string
  urlAppGui: string
  urlAppOpen: string
  staticDir: string
  pathToDefaultPlayList: string
  urlSettings: string
  hwnds = %*[]
  msg: MSG
  isGUI = 1
  nid*: NOTIFYICONDATAW = NOTIFYICONDATAW()
  logger = newColoredLogger(fmtStr = coolerFmtStr)






proc echoError(error:string="error", typeError:string="warn") =
  var window = newWindow()
  var listButtons = @["OK", "Exit", "Copy Error"]
  var res: int

  if typeError == "warn":
    res = window.msgBox(message = "Error:\n\n" & $error,
        title = "Melon Player Error:", button1 = fmt"{listButtons[0]}")
    warn fmt"{error}"
  elif typeError == "error":
    res = window.msgBox("Error:\n\n" & $error, "Melon Player Error:",
        listButtons[0], listButtons[1], listButtons[2])
    error fmt"error"
  elif typeError == "fatal":
    listButtons[2] = "copy Error and Exit"
    res = window.msgBox("Error:\n\n" & $error, "Melon Player Error:",listButtons[0], listButtons[1], listButtons[2])
    fatal fmt"{error}"
  # hwnds.add(%*{"Melon_Player_Error":FindWindowA(nil,"Melon Player Error:")})
  # fatal $hwnds
  if res-1 == 0:
    discard
  elif res-1 == 1:
      quit()
  elif res-1 == 2:
    var cb = clipboard_new(nil)
    if cb.clipboard_set_text(error):
      if typeError=="fatal":
        quit()
      discard
    else:
      fatal fmt"error copy text: {error}"
  else:
    discard






proc consoleWindow() =
  return




proc getSettings() =
  try:
    var data_app_settings = getData(key = @["app_settings"])
    var data_app_info = getData(key = @["app_info"])

    version = data_app_info["version"].getStr()
    settings_version = data_app_info["settings_version"].getStr()
    port = data_app_settings["port"].getInt()
    urlServer = data_app_settings["url_server"].getStr()
    urlAppGui = fmt"{urlServer}:{port}/app.gmp"
    urlAppOpen = fmt"{urlServer}:{port}/appOpen.mp"

    staticDir = getPaths(getData(key = @["app_settings"]), "paths")[0]
    pathToDefaultPlayList = "./playList/"

    width = data_app_settings["width"].getInt()
    height = data_app_settings["height"].getInt()
  except CatchableError as error:
    fatal "error:\n" & "name: " & $error.name & "\n↳ exception: " & $error.msg
    echoError(error="error app initialization:\nNameException:" & $error.name & "\n↳ exception: " & $error.msg, typeError="fatal")







proc setSettings() =
  return






proc startApp() =
  try:
    info "Starting app"
    var m_app = newWebView(path = "http://localhost:3487/app.gmp", title = appName, width = width, height = height, debug = false)
    var trayApp = newSystemTray($appName, $appName,
        "./assets/icons/iconMelonPlayer.ico", onClickLeftButton = proc() = isGUI = 1)
    trayApp.menu.items.add(ContextItem(title: "Открыть",
        onClickRightButton: proc() = isGUI = 1))
    trayApp.menu.items.add(ContextItem(title: "Скрыть",
        onClickRightButton: proc() = isGUI = 0))
    trayApp.menu.items.add(ContextItem(title: "Выйти",
        onClickRightButton: proc() = quit()))

    SetWindowLongPtr(m_app.priv.windowHandle, GWL_EXSTYLE, GetWindowLongPtr(
        m_app.priv.windowHandle, GWL_EXSTYLE) or WS_EX_TOOLWINDOW)


    let (screenWidth, screenHeight) = mGetSystemMetrics()
    mSetWindowPosition(m_app, windowOrder = HWND_TOPMOST, x = screenWidth-width-1,
        y = screenHeight-height-50)
    mShowWindowIcons(m_app, height = LONG(height), width = LONG(width))
    mSetTitle(m_app, $appName)
    ShowWindow(m_app.priv.windowHandle, SW_SHOW)

    info "started Application cycle..."
    while GetMessage(msg.addr, 0, 0, 0) != -1:
      if global_values.len() != 0:
        if "echoError" in global_values:
          debug "echoError:" & global_values["echoError"][0].getStr() & "type:" & global_values["echoError"][1].getStr()
          echoError(error=global_values["echoError"][0].getStr(), typeError=global_values["echoError"][1].getStr())
          global_values = %*{}
        else:
          if "s_Windows" in global_values and global_values["s_Windows"].getStr() == "open":
            debug "s_Windows: open" & "\n↳ " & "link: http://localhost:3487/appSettings.gmp"
            global_values = %*{}
            var settingWindow = newWebView(fmt"http://localhost:3487/appSettings.gmp", title = "Settings MelonPlayer", width = width, height = height, debug = false)
          else:
            discard
      nid = trayApp.tInit()
      if isGUI == 0:
        ShowWindow(m_app.priv.windowHandle, SW_HIDE)
      elif isGUI == 1 and IsWindowVisible(m_app.priv.windowHandle) == 0:
        ShowWindow(m_app.priv.windowHandle, SW_SHOW)


      if msg.hwnd != 0:
        # debug "msg.hwnd != 0"
        TranslateMessage(msg.addr)
        DispatchMessage(msg.addr)
        continue
      case msg.message:
        of WM_MOVE:
          debug "↳ WM_MOVE"
          continue
        of WM_APP:
          debug "↳ WM_APP"
          let fn = cast[proc(env: pointer): void {.stdcall.}](msg.lParam)
          fn(cast[pointer](msg.wParam))
        of WM_CLOSE:
          debug "↳ WM_QUIT"
        of WM_QUIT:
          debug "↳ WM_QUIT"
        of WM_COMMAND,
          WM_KEYDOWN,
          WM_KEYUP:
          if (msg.wParam == VK_F5):
            debug "↳ WM_COMMAND, WM_KEYDOWN, WM_KEYUP"
            return
        else:
          discard

    Shell_NotifyIcon(NIM_DELETE, &nid)
  except CatchableError as error:
    fatal "error app:\n" & "name: " & $error.name & "\n↳ exception: " & $error.msg
    echoError(error="error app initialization:\nNameException:" & $error.name & "\n↳ texception: " & $error.msg, typeError="fatal")


proc startServer() {.thread, nimcall.} =
  info "server thread started!"
  router webServer:
    get "/app.gmp":
      resp readfile("./assets/html/index.html")
    get "/appSettings.gmp":
      resp readfile("./assets/html/settings.html")
    get "/addMusik@url?":
      resp "Adds musik....!"
    get "/appOpen.mp":
      isGUI = 1
      resp "app start  "
    get "/appClose.mp":
      isGUI = 0
      resp "Closing app..."
    get "/send@command?":
      if @"command" != "" and @"command".len() != 0:
        var arguments = findTextInBrackets(@"command")
        var packet: seq[string]
        if arguments.len() != 0:
          for (index, text) in arguments:
            let commandPart = @"command"[0..<index] # Получаем подстроку "command" до заданного индекса
            var parts = commandPart.split(".") # Разбиваем подстроку по символу '.'
            # Добавляем полученные части в список packet
            packet.add(parts)
            packet.add(text)

        else:
          packet = @"command".split(".") # Разбиваем подстроку по символу '.'

        debug "a new design: " & @"command"
        if @"command" == "Слава Украине!" or @"command" == "Слава Україні!":
          resp "Героям Слава!"

        elif packet.len() <= 1:
          resp "{:red:}" & @"command" & " not found!"
          error @"command" & " not found!"
        elif packet[1] == "":
          error @"command" & " not found!"
          resp "{:red:}" & @"command" & " not found!"


        elif packet[0] == "get":
          info "get"
          case packet[1]
          of "Default":
            info "↳ Default"
            if packet[2] == "PlayList":
              info "  ↳ Default"
              let PlayList = getDirectoryContents("./assets/playList")
              notice "      ↳ Ok: " & $PlayList
              resp $PlayList
          of "Settings":
            info "↳ Settings"
            if packet[2] == "Window":
              info "  ↳ Window"
              withLock lock:
                {.gcsafe.}:
                  global_values["s_Windows"] = %*"open"
              notice "      ↳ s_Windows: ok"
              resp "s_Windows"
          of "Data":
            info "↳ Data"
            var pac = getData(key = packet[2..<packet.len])
            notice "  ↳ Data: " & $pac
            resp $pac
          else:
            error @"command" & " not found!"
            resp "{:red:}error:" & @"command" & " not found!"


        elif packet[0] == "send":
          info "send"
          case packet[1]
          of "error":
            info "↳ error"
            withLock lock:
              {.gcsafe.}:
                global_values["echoError"] = %*packet[packet.len()-1].split(",")
                notice "echoError: ok"
            resp("Start the Error windows")
          else:
            error @"command" & " not found!"
            resp "{:red:}" & @"command" & " not found!"
        elif packet[0] == "add":
          case packet[1]:
            of "music": 
              var url: seq[string] = @[packet[2]]
              var info: seq[string] = @["123","123"]
              # withLock lock:
              {.gcsafe.}:
                asyncCheck download_youtube_music(url=url)
              resp "warte"
            of "music_and_preview":discard
            of "preview":discard

                
            else:
              error @"command" & " not found!"
              resp "{:red:}" & @"command" & " not found!"
        else:
          error @"command" & " not found!"
          resp "{:red:}" & @"command" & " not found!"
      else:
        error @"command" & " not found!"
        resp @"command" & " not found!"


        
  var jester = initJester(webServer, settings = newSettings(port = Port(port),
      staticDir ="./assets"))
  info "web server startet!"
  jester.serve()
  runForever()


proc webSocketServer(req: asynchttpserver.Request) {.async.} =
  if req.url.path == "/ws":
    var ws: WebSocket
    try:
      ws = await newWebSocket(req)
      echo "Connected..."
      await ws.send("Welcome to simple chat server")


      proc reader() {.async.} =
        # Loops while socket is open, looking for messages to read
        while ws.readyState == Open:
          echo "Waiting for message..."
          
          try:
            var packet = await ws.receiveStrPacket()
            echo packet
          except CatchableError as e:
            echo "Error receiving packet: ", e.msg
            break


      proc writer() {.async.} =
        try:
          ## Loops while socket is open, looking for messages to write
          echo "writer"
          while ws.readyState == Open:
            # {.gcsafe.}:
              var el = atomic_variable.load[]
              if el != "":
                await ws.send(el)
                await sleepAsync(500)
        except CatchableError as e:
          echo "Error in WebSocket connection: ", e.msg 
      asyncCheck reader()
      await writer()
    except CatchableError as e:
      echo "Error in WebSocket connection: ", e.msg
    finally:
      if ws.readyState != Closed:
        ws.close()
      echo "Connection closed."








proc sc() {.thread, nimcall.} =
  try:
    var server = newAsyncHttpServer()
    waitFor server.serve(Port(9001), webSocketServer)
  except CatchableError as error:
    fatal "error app:\n" & "name: " & $error.name & "\n↳ exception: " & $error.msg
    # {.gcsafe.}:
    #   echoError(error="error wss initialization:\nNameException:" & $error.name & "\n↳ exception: " & $error.msg, typeError="fatal")

    
proc main() =
  try:
    info "create a new theard..."
    var startServerThread: Thread[void]
    createThread(startServerThread, startServer)
    var startSocketServerThread: Thread[void]
    createThread(startSocketServerThread, sc)

    notice "new theard created!"
    notice "open url: " & urlAppOpen
    info "starting gui app..."
    while true:
      startApp()
  except CatchableError as error:
    fatal "error app:\n" & "name: " & $error.name & "\n↳ exception: " & $error.msg
    echoError(error="error app initialization:\nNameException:" & $error.name & "\n↳ exception: " & $error.msg, typeError="fatal")




when isMainModule:
  getSettings()
  addHandler(logger)
  main()
