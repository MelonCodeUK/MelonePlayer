package MelonePlayer

import "github.com/gen2brain/beeep"

func Notify(title string, msg string) {
	err := beeep.Notify(title, msg, "")
	if err != nil {
		panic(err)
	}
}
