package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alexflint/go-arg"
)

var updateUrl = "https://raw.githubusercontent.com/MelonCodeUK/MelonePlayer/refs/heads/main/installer_updater/version.json"

func Українська(skip_license bool) {
	if skip_license == true {

	} else {
		isLicenseAccept := ""
		fmt.Println(GetData())
		fmt.Println("Чи приймаете ви згоду?(yes - так\\no - ні)")
		fmt.Scanln(&isLicenseAccept)
		if isLicenseAccept == "yes" {

		} else {

			os.Exit(1)
		}

	}
}

func GetData(url string) (string, error) {
	// Выполняем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Возвращаем тело ответа как строку
	return string(body), nil
}

func English() {

}

func Reinstall() {

}

func Delete() {

}

func Update() {

}

func UpdateInstaller() {

}

// Структура для команды установки (install)
type InstallCmd struct {
	Path        string `arg:"positional" help:"installation path"`
	Language    string `arg:"--language,-l" help:"set a language"`
	SkipLicense bool   `arg:"--skipLicense,-sL" help:"skiping license"`
	Licence     string `arg:"--licence,-l,l,licence," help:"write --licence yes or no"`
}

type Args struct {
	Licence         string      `arg:"--licence,-l,l,licence," help:"write --licence and your land.\n Example: --licence ua"`
	Install         *InstallCmd `arg:"subcommand:install"`
	Reinstall       bool        `arg:"--reinstall,-RI,-r" help:"reinstalling"`
	Delete          bool        `arg:"--delete,-d,--remove,-r" help:"removing a App"`
	Update          bool        `arg:"--update,-u" help:"update a App"`
	UpdateInstaller bool        `arg:"--update-installer, -uI"  help:"update installer"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	if args.Install != nil || len(args.Licence) != 0 || args.Reinstall || args.Delete || args.Update || args.UpdateInstaller {
		if len(args.Install.Path) != 0 {
			if args.Install.Licence == "yes" {
				if args.Install.SkipLicense {
					if args.Install.Language == "ua" {
						Українська(args.Install.SkipLicense)
					}
				}
			} else {
				os.Exit(1)
			}
		} else if args.Licence == "ua" {
			fmt.Println(license_agreement_ua)
		} else if args.Licence == "en" {
		} else if args.Reinstall {
			Reinstall()
		} else if args.Delete {
			Delete()
		} else if args.Update {
			Update()
		} else if args.UpdateInstaller {
			UpdateInstaller()
		} else {

		}

	} else {
		var choice = 0
		for choice == 0 {
			// Вывод меню выбора языка
			fmt.Println("1) Українська")
			fmt.Println("2) English")

			fmt.Scanln(&choice)

			switch choice {
			case 1:
				Українська(false)
				os.Exit(0)
			case 2:
				English()
				os.Exit(0)

			default:

			}
		}
	}

}
