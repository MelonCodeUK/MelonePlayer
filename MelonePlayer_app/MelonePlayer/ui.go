package MelonePlayer

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
	. "github.com/gonutz/w32/v2"
	webview "github.com/webview/webview_go"
)

func MajorUi(w Window, url string) {
	w_ := webview.New(IsDebug)
	w.hwnd = w32.HWND(w_.Window())
	if w.Title == MainWindowSettings.Title {
		MainWindowSettings.hwnd = w.hwnd
	} else if w.Title == SettingsWindows.Title {
		SettingsWindows.hwnd = w.hwnd
		w32.ShowWindow(SettingsWindows.hwnd, w32.SW_HIDE)
	}
	w_.SetSize(w.Width, w.Height, webview.HintNone)
	ShowWindowIcons(w)
	if SHAppBar_AUTOHIDE() == false {
		_, y := GetSizeTaskBar()
		w.Y = w.Y - y
	}
	SetPosition(w)
	w_.SetTitle(w.Title)
	originalWndProc = SetWindowLongPtr(w.hwnd, GWL_WNDPROC, syscall.NewCallback(wndProc))
	w_.Navigate(url)
	w_.Run()
}

func SetTitle(hwnd w32.HWND, title string) {
	SetWindowText(hwnd, title)
}

var (
	originalWndProc uintptr
)

func wndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case w32.WM_CLOSE:
		// Вместо закрытия окна скрываем его
		w32.ShowWindow(hwnd, SW_HIDE)
		return 0
	}
	return w32.CallWindowProc(originalWndProc, hwnd, msg, wParam, lParam)
}

// ShowWindowIcons управляет стилями окна
func ShowWindowIcons(w Window) {
	var rect RECT
	rect.Left = 0
	rect.Top = 0
	rect.Right = int32(w.Width)
	rect.Bottom = int32(w.Height)

	style := GetWindowLong(w.hwnd, GWL_STYLE)
	ex_style := GetWindowLongPtr(w.hwnd, GWL_EXSTYLE)
	var pos_style uint

	if w.WS_THICKFRAME {
		style |= WS_THICKFRAME
	} else {
		style &^= WS_THICKFRAME
	}

	if w.WS_CAPTION {
		style |= WS_CAPTION
	} else {
		style &^= WS_CAPTION
	}

	if w.WS_SYSMENU {
		style |= WS_SYSMENU
	} else {
		style &^= WS_SYSMENU
	}

	if w.WS_MINIMIZEBOX {
		style |= WS_MINIMIZEBOX
	} else {
		style &^= WS_MINIMIZEBOX
	}

	if w.WS_MAXIMIZEBOX {
		style |= WS_MAXIMIZEBOX
	} else {
		style &^= WS_MAXIMIZEBOX
	}
	if w.WS_EX_TOOLWINDOW {
		ex_style &^= WS_EX_TOOLWINDOW
	} else {
		ex_style |= WS_EX_TOOLWINDOW
	}
	if w.WS_EX_TOPMOST {
		pos_style |= WS_EX_TOPMOST
	} else {
		pos_style &^= WS_EX_TOPMOST
	}
	if w.SWP_NOMOVE {
		pos_style |= SWP_NOMOVE
	} else {
		pos_style &^= SWP_NOMOVE
	}
	if w.SWP_NOZORDER {
		pos_style |= SWP_NOZORDER
	} else {
		pos_style &^= SWP_NOZORDER

	}

	// Применение нового стиля
	SetWindowLong(w.hwnd, GWL_STYLE, style)
	SetWindowLongPtr(w.hwnd, GWL_EXSTYLE, ex_style)
	SetWindowPos(w.hwnd, HWND_TOPMOST, 0, 0, int(rect.Right-rect.Left), int(rect.Bottom-rect.Top), SWP_FRAMECHANGED|pos_style)

}

// SetPosition изменяет позицию окна
func SetPosition(w Window) {
	SetWindowPos(w.hwnd, HWND_TOPMOST, w.X, w.Y, 0, 0, SWP_NOSIZE)
}

// GetSystemMetrics возвращает ширину и высоту экрана
func GetDisplayResolution() (int, int) {
	x := GetSystemMetrics(SM_CXSCREEN)
	y := GetSystemMetrics(SM_CYSCREEN)
	return x, y
}

func Error(type_error string, error string) {

}

func SHAppBar_AUTOHIDE() bool {
	// Загружаем библиотеку shell32.dll
	shell32 := syscall.NewLazyDLL("shell32.dll")
	SHAppBarMessage := shell32.NewProc("SHAppBarMessage")

	// Создаем структуру APPBARDATA
	var abd APPBARDATA
	abd.cbSize = uint32(unsafe.Sizeof(abd))

	// Вызов функции SHAppBarMessage с параметром ABM_GETSTATE
	ret, _, _ := SHAppBarMessage.Call(uintptr(ABM_GETSTATE), uintptr(unsafe.Pointer(&abd)))

	// Проверяем результат
	if ret&ABS_AUTOHIDE == ABS_AUTOHIDE {
		return true
	} else {
		return false
	}
}

func GetSizeTaskBar() (int, int) {
	hwnd := FindWindow("Shell_TrayWnd", "")
	if hwnd == 0 {
		Log.Error("Failed to find the taskbar window.")
		return 0, 0
	}

	var rect = GetWindowRect(hwnd)
	if rect != nil {
		width := rect.Right - rect.Left
		height := rect.Bottom - rect.Top
		Log.Debug(fmt.Sprintf("Taskbar size: width=%d, height=%d\n", int(width), int(height)))
		return int(width), int(height)
	} else {
		Log.Error("Failed to get the taskbar size.")
		return 0, 0
	}
}
