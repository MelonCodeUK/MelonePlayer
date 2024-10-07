package main

import (
	. "MelonePlayer/MelonePlayer"
	"fmt"

	"github.com/getlantern/systray"
)

func main() {
	GetSettings()
	PrintInfo()
	go Server()
	go func() { // Player Window
		MainWindowSettings.WS_CAPTION = false
		MainWindowSettings.WS_MAXIMIZEBOX = false
		MainWindowSettings.WS_MINIMIZEBOX = false
		MainWindowSettings.WS_SYSMENU = false
		MainWindowSettings.WS_THICKFRAME = false
		MainWindowSettings.WS_EX_TOOLWINDOW = false
		MainWindowSettings.WS_EX_TOPMOST = false
		MainWindowSettings.SWP_NOZORDER = true
		MainWindowSettings.SWP_NOMOVE = true
		MajorUi(MainWindowSettings, fmt.Sprintf("%s:%d/player.window?cache-control=no-cache&pragma=no-cache", Server_Url, Port))
	}()
	go func() { // Settings Window
		SettingsWindows.WS_CAPTION = true
		SettingsWindows.WS_MAXIMIZEBOX = false
		SettingsWindows.WS_MINIMIZEBOX = false
		SettingsWindows.WS_SYSMENU = true
		SettingsWindows.WS_THICKFRAME = true
		SettingsWindows.WS_EX_TOOLWINDOW = true
		SettingsWindows.WS_EX_TOPMOST = true
		SettingsWindows.SWP_NOZORDER = true
		SettingsWindows.SWP_NOMOVE = true
		SettingsWindows.Height = 200
		SettingsWindows.Width = 300
		SettingsWindows.Title = "Settings"
		MajorUi(SettingsWindows, fmt.Sprintf("%s:%d/settings.window?cache-control=no-cache&pragma=no-cache", Server_Url, Port))
	}()
	go systray.Run(Init, OnExit)

	select {} // Блокируем основной поток для предотвращения завершения программы
}
