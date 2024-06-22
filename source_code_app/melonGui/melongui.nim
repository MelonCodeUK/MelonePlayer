import winim
import crowngui
import winim/mean
import strutils
import winim/lean
import winim/com
import std/enumerate


proc mSetTitle*(w: Webview; title: string ): void =
  SetWindowTextW(w.priv.windowHandle, &T(title))
  # текст окна


proc mShowWindowIcons*(w: Webview; 
                       mWS_SYSMENU: bool = false, #Управляет наличием системного меню (и кнопки закрытия).
                       mWS_MINIMIZEBOX: bool = false, #Управляет наличием кнопки минимизации.
                       mWS_MAXIMIZEBOX: bool = false, #Управляет наличием кнопки максимизации.
                       mWS_CAPTION: bool = false, #Добавляет или убирает заголовок окна.
                       mWS_THICKFRAME: bool = false, #Управляет возможностью изменения размера окна.
                       width: LONG, height: LONG
                       ): void =
  ## Получаем дескриптор окна из объекта Webview.
  let hwnd = w.priv.windowHandle
  var rect: RECT
  rect.left = 0
  rect.top = 0
  rect.right = width
  rect.bottom = height


  # Получаем текущий стиль окна.
  var style = GetWindowLong(hWnd, GWL_STYLE).int32

  # WS_THICKFRAME: Управляет возможностью изменения размера окна.
  if mWS_THICKFRAME :
    style = style or WS_THICKFRAME 
  else:
    style = style and not WS_THICKFRAME 

  # WS_CAPTION: Добавляет или убирает заголовок окна.
  if mWS_CAPTION:
    style = style or WS_CAPTION
  else:
    style = style and not WS_CAPTION

  # WS_SYSMENU: Управляет наличием системного меню (и кнопки закрытия).
  if mWS_SYSMENU:
    echo "mWS_SYSMENU"
    style = style or WS_SYSMENU
  else:
    style = style and not WS_SYSMENU

  # WS_MINIMIZEBOX: Управляет наличием кнопки минимизации.
  if mWS_MINIMIZEBOX:
    style = style or WS_MINIMIZEBOX
  else:
    style = style and not WS_MINIMIZEBOX 

  # WS_MAXIMIZEBOX: Управляет наличием кнопки максимизации.
  if mWS_MAXIMIZEBOX:
    style = style or WS_MAXIMIZEBOX
  else:
    style = style and not WS_MAXIMIZEBOX
    
  # Выводим в консоль текущие стили для отладки.
  echo "styles:" & $style
  # Применяем новый стиль окна.
  SetWindowLong(hWnd, GWL_STYLE, style or WS_POPUP)
  SetWindowPos(hwnd, nil, 0, 0, rect.right - rect.left, rect.bottom - rect.top, SWP_FRAMECHANGED or SWP_NOMOVE or SWP_NOZORDER)





proc mSetWindowPosition*(w: Webview; windowOrder: HWND, x: int, y: int): void =
  # позицыя окна
  # x и y: Новые координаты левого верхнего угла окна.
  # Параметр windowOrder - размешает окно:
    # HWND_TOP: Размещает окно на верхней позиции порядка Z. Окно становится самым верхним окном, которое не является топмост (всегда поверх других окон).
    # HWND_BOTTOM: Размещает окно на самой нижней позиции порядка Z. Все другие окна (кроме всегда видимых окон) располагаются над этим окном.
    # HWND_TOPMOST: Размещает окно над всеми не-топмост окнами, даже если окно не активно. Окно остается над другими окнами, даже когда оно неактивно. Это полезно для создания окон, которые должны быть всегда видимыми.
    # HWND_NOTOPMOST: Размещает окно над всеми обычными (не-топмост) окнами, но под любыми топмост окнами. Это значение используется, когда окно должно быть над обычными окнами, но не всегда на переднем плане.
  let hwnd = w.priv.windowHandle
  SetWindowPos(hWnd=hwnd, hWndInsertAfter=windowOrder, X=int32(x),Y=int32(y), cx=0, cy=0, uflags = SWP_NOSIZE)


proc mGetSystemMetrics*():(int, int)=
  # возврощает ширину и высоту дисплея
  let x = int(GetSystemMetrics(SM_CXSCREEN))
  let y = int(GetSystemMetrics(SM_CYSCREEN))
  return (x, y)


proc nilProc*()=
  discard



type ContextItem* = object
  title*: string
  # onClickLeftButton: proc()
  onClickRightButton*: proc()

type ContextMenu* = object
  items*: seq[ContextItem]

type SystemTray* = object
  tooltip*: string
  className*: string
  icon*: string
  onClickLeftButton*: proc()
  onClickRightButton*: proc()
  menu*: ContextMenu



var globalSystemTray* = SystemTray()


proc newSystemTray*(title: string, className: string, icon: string, onClickLeftButton:proc() = nilProc, onClickRightButton: proc() = nilProc): SystemTray =
  var tray = SystemTray()
  tray.tooltip = title
  tray.className = className
  tray.icon = icon
  tray.onClickLeftButton = onClickLeftButton
  tray.onClickRightButton = onClickRightButton
  tray.menu = ContextMenu()
  return tray




proc newNOTIFYICONDATA*(hWnd: HWND, 
                      ID_TRAY_ICON: UINT, 
                      WM_TRAY: UINT, 
                      icon: string = "0", 
                      tooltip: string = "App", 
                      uFlags: UINT = NIF_MESSAGE or NIF_ICON or NIF_TIP): NOTIFYICONDATA =
  var tooltip = tooltip
  var nid: NOTIFYICONDATA
  nid.cbSize = windef.DWORD(sizeof(NOTIFYICONDATA))
  nid.hWnd = hWnd
  nid.uID = ID_TRAY_ICON
  nid.uFlags = uFlags
  nid.uCallbackMessage = WM_TRAY
  if icon.endsWith".ico":
    nid.hIcon = cast[HICON](LoadImage(nil, icon, IMAGE_ICON, 0, 0, LR_LOADFROMFILE))
  else:
    nid.hIcon = LoadIcon(0, IDI_APPLICATION)
  if tooltip.len > 63: tooltip = tooltip[0 ..< 63]
  let wTooltip = +$tooltip
  for i in 0 ..< wTooltip.len: 
      nid.szTip[i] = wTooltip[i]
  return nid


proc newWindowForSystemTray*(WndProc: WNDPROC, tray: SystemTray): WNDCLASSEX =
  globalSystemTray = tray
  var wc: WNDCLASSEX
  wc.cbSize = sizeof(WNDCLASSEX).UINT
  wc.lpfnWndProc = WndProc
  wc.hInstance = GetModuleHandle(nil)
  wc.lpszClassName = tray.className
  wc.style = CS_HREDRAW or CS_VREDRAW
  wc.hIcon = LoadIcon(0, IDI_APPLICATION)
  wc.hCursor = LoadCursor(0, IDC_ARROW)
  wc.hbrBackground = (HBRUSH) (COLOR_WINDOW + 1)
  return wc

const
  WM_TRAY = WM_USER + 1
  ID_TRAY_ICON = 1000

proc WndProc*(hWnd: HWND, message: UINT, wParam: WPARAM, lParam: LPARAM): LRESULT {.stdcall.} =
  # echo int(LOWORD(wParam))
  case message
  of WM_TRAY:
    if lParam == WM_RBUTTONDOWN:
      globalSystemTray.onClickRightButton()
      var
        pt: POINT
        hMenu = CreatePopupMenu()
      for index, i in enumerate(globalSystemTray.menu.items):
        InsertMenu(hMenu, UINT(index), MF_BYPOSITION or MF_STRING, int32(1000+index), newWideCString(i.title))  # Добавляем новый пункт 
      SetForegroundWindow(hWnd)
      GetCursorPos(&pt)
      TrackPopupMenu(hMenu, TPM_BOTTOMALIGN or TPM_LEFTALIGN, pt.x, pt.y, 0, hWnd, nil)
      PostMessage(hWnd, WM_NULL, 0, 0) # to allow SetForegroundWindow to work
    if lParam == WM_LBUTTONDOWN:
      echo "diggar"
      PostMessage(hWnd, WM_NULL, 0, 0)
      globalSystemTray.onClickLeftButton()
  of WM_COMMAND:
    echo int(LOWORD(wParam))
    for index, i in enumerate(globalSystemTray.menu.items):
      # echo LOWORD(wParam)
      if int(LOWORD(wParam)) == int(1000+index):
        i.onClickRightButton()
    # if LOWORD(wParam) == ID_TRAY_EXIT:
    #   PostQuitMessage(0)
    # elif LOWORD(wParam) == ID_TRAY_HELLO:
    #   echo "HELLO"
    
  else:
    result = DefWindowProc(hWnd, message, wParam, lParam)

proc tInit*(tray: SystemTray): NOTIFYICONDATAW =
  var wc = newWindowForSystemTray(WndProc, tray)

  if RegisterClassEx(wc) != 0:
    let hWnd = CreateWindowEx(0, tray.className, tray.tooltip, WS_OVERLAPPEDWINDOW, 
                              CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, 
                              nil, nil,  GetModuleHandle(nil), nil)
    
    var nid = newNOTIFYICONDATA(hWnd=hWnd,ID_TRAY_ICON=ID_TRAY_ICON, WM_TRAY=WM_TRAY, icon=tray.icon, tooltip=tray.tooltip)

    Shell_NotifyIcon(NIM_ADD, &nid)

    var msg: MSG
    # while GetMessage(&msg, nil, 0, 0) != 0:
    TranslateMessage(&msg)
    DispatchMessage(&msg)
    return nid
  return

    # Shell_NotifyIcon(NIM_DELETE, &nid)

# var trat= newSystemTray(title="System", className="systems2", icon="./../assets/icons/iconMelonPlayer.ico", onClickLeftButton=proc()=echo "active")
# trat.menu.items.add(ContextItem(title:"hi", onClickRightButton:proc()=echo 123))
# trat.menu.items.add(ContextItem(title:"hi0", onClickRightButton:proc()=echo 1243))


# trat.tInit()
