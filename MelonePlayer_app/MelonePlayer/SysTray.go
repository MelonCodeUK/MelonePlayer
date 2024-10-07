package MelonePlayer

import (
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/gonutz/w32/v2"
)

func Init() {
	iconData, err := os.ReadFile(Path.StaticDir + "/icons/iconMelonPlayer.ico")
	if err != nil {
		log.Fatal(err)
	}
	systray.SetIcon(iconData)
	systray.SetTitle(MainWindowSettings.Title)
	systray.SetTooltip(MainWindowSettings.Title)
	hidePlayer := systray.AddMenuItem("Hide player", "")
	showPlayer := systray.AddMenuItem("Show player", "")
	Quit := systray.AddMenuItem("Quit", "Close FULL application")

	go func() {
		for {
			select {
			case <-hidePlayer.ClickedCh:
				w32.ShowWindow(MainWindowSettings.hwnd, w32.SW_HIDE)
			case <-showPlayer.ClickedCh:
				w32.ShowWindow(MainWindowSettings.hwnd, w32.SW_SHOW)
			case <-Quit.ClickedCh:
				os.Exit(0)

			}
		}
	}()
}

func OnExit() {
	// clean up here
}
