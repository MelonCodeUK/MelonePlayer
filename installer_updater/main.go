// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/alexflint/go-arg"
// )

// var license_agreement_ua = "Ліцензійна угода\n" +
// 	"Увага! Будь ласка, уважно прочитайте цю угоду перед використанням програмного забезпечення.\n" +
// 	"1. Загальні положення\n\n" +
// 	"Ця ліцензійна угода регулює використання програмного забезпечення MelonePlayer, яке надається MelonCodeUK. Використовуючи це Програмне забезпечення, ви погоджуєтеся з умовами цієї Угоди. Якщо ви не згодні з умовами, ви не маєте права використовувати це Програмне забезпечення.\n" +
// 	"2. Ліцензія на використання\n\n" +
// 	"Ліцензіар надає вам обмежену, непередавану, невиключну ліцензію на використання Програмного забезпечення виключно відповідно до умов цієї Угоди.\n" +
// 	"3. Заборона на модифікацію та втручання\n\n" +
// 	"Ви не маєте права:\n\n" +
// 	"    Модифікувати, адаптувати, змінювати, декомпілірувати, дизасемблювати або будь-яким іншим чином намагатися витягти вихідний код Програмного забезпечення.\n" +
// 	"    Змінювати, видаляти або приховувати будь-які повідомлення про права інтелектуальної власності, що містяться в Програмному забезпеченні.\n" +
// 	"    Використовувати або встановлювати Програмне забезпечення в будь-якому місці, де це не дозволено відповідно до цієї Угоди.\n\n" +
// 	"Ви можете модифікувати вихідний код, після чого скомпілювати його та вказати оригінальний репозиторій. Також ви можете змінювати файл налаштувань settings.json. Дозволяється встановлювати розширення/пакети. Ви також можете встановлювати інші модифікації, але в таких випадках повинна бути інформація про оригінальний репозиторій.\n" +
// 	"4. Обмеження відповідальності\n\n" +
// 	"Ліцензіар не несе відповідальності за будь-які збитки, що виникають внаслідок використання або неможливості використання Програмного забезпечення, включаючи, але не обмежуючись, втратою прибутку, даних або інших нематеріальних збитків.\n" +
// 	"5. Зміни в Угоді\n\n" +
// 	"Ліцензіар залишає за собою право змінювати умови цієї Угоди в будь-який час. Ваше подальше використання Програмного забезпечення після внесення змін вважатиметься вашим погодженням з зміненими умовами.\n" +
// 	"6. Припинення дії\n\n" +
// 	"Ліцензія буде діяти доти, поки не буде припинена вами або Ліцензіаром. Ви можете припинити дію ліцензії, припинивши використання Програмного забезпечення. Ліцензіар може припинити дію ліцензії у разі порушення вами умов цієї Угоди.\n" +
// 	"7. Прийняття умов\n\n" +
// 	"Використовуючи це Програмне забезпечення, ви підтверджуєте, що прочитали та зрозуміли умови цієї Угоди і погоджуєтеся з ними."

// var updateUrl = "https://raw.githubusercontent.com/MelonCodeUK/MelonePlayer/refs/heads/main/installer_updater/version.json"

// func Українська(skip_license bool) {
// 	if skip_license == true {

// 	} else {
// 		isLicenseAccept := ""
// 		fmt.Println(license_agreement_ua)
// 		fmt.Println("Чи приймаете ви згоду?(yes - так\\no - ні)")
// 		fmt.Scanln(&isLicenseAccept)
// 		if isLicenseAccept == "yes" {

// 		} else {

// 			os.Exit(1)
// 		}

// 	}
// }

// func English() {

// }

// func Reinstall() {

// }

// func Delete() {

// }

// func Update() {

// }

// func UpdateInstaller() {

// }

// // Структура для команды установки (install)
// type InstallCmd struct {
// 	Path        string `arg:"positional" help:"installation path"`
// 	Language    string `arg:"--language,-l" help:"set a language"`
// 	SkipLicense bool   `arg:"--skipLicense,-sL" help:"skiping license"`
// 	Licence     string `arg:"--licence,-l,l,licence," help:"write --licence yes or no"`
// }

// type Args struct {
// 	Licence         string      `arg:"--licence,-l,l,licence," help:"write --licence and your land.\n Example: --licence ua"`
// 	Install         *InstallCmd `arg:"subcommand:install"`
// 	Reinstall       bool        `arg:"--reinstall,-RI,-r" help:"reinstalling"`
// 	Delete          bool        `arg:"--delete,-d,--remove,-r" help:"removing a App"`
// 	Update          bool        `arg:"--update,-u" help:"update a App"`
// 	UpdateInstaller bool        `arg:"--update-installer, -uI"  help:"update installer"`
// }

// func main() {
// 	var args Args
// 	arg.MustParse(&args)

// 	if args.Install != nil || len(args.Licence) != 0 || args.Reinstall || args.Delete || args.Update || args.UpdateInstaller {
// 		if len(args.Install.Path) != 0 {
// 			if args.Install.Licence == "yes" {
// 				if args.Install.SkipLicense {
// 					if args.Install.Language == "ua" {
// 						Українська(args.Install.SkipLicense)
// 					}
// 				}
// 			} else {
// 				os.Exit(1)
// 			}
// 		} else if args.Licence == "ua" {
// 			fmt.Println(license_agreement_ua)
// 		} else if args.Licence == "en" {
// 		} else if args.Reinstall {
// 			Reinstall()
// 		} else if args.Delete {
// 			Delete()
// 		} else if args.Update {
// 			Update()
// 		} else if args.UpdateInstaller {
// 			UpdateInstaller()
// 		} else {

// 		}

// 	} else {
// 		var choice = 0
// 		for choice == 0 {
// 			// Вывод меню выбора языка
// 			fmt.Println("1) Українська")
// 			fmt.Println("2) English")

// 			fmt.Scanln(&choice)

// 			switch choice {
// 			case 1:
// 				Українська(false)
// 				os.Exit(0)
// 			case 2:
// 				English()
// 				os.Exit(0)

// 			default:

// 			}
// 		}
// 	}

// }

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	downloadURL = "http://example.com/path/to/your/new/version" // URL для скачивания новой версии
	tempFile    = "new_version"                                 // Имя временного файла
)

func main() {
	// Скачиваем новую версию
	if err := downloadNewVersion(downloadURL, tempFile); err != nil {
		fmt.Println("Ошибка при скачивании новой версии:", err)
		return
	}

	// Завершаем текущий процесс
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Ошибка при получении пути к исполняемому файлу:", err)
		return
	}

	// Задержка для завершения работы программы
	time.Sleep(2 * time.Second)

	// Запускаем процесс обновления
	if err := os.Rename(tempFile, executablePath); err != nil {
		fmt.Println("Ошибка при замене файла:", err)
		return
	}

	// Перезапускаем приложение
	if err := exec.Command(executablePath).Start(); err != nil {
		fmt.Println("Ошибка при перезапуске приложения:", err)
		return
	}

	fmt.Println("Обновление завершено. Приложение перезапускается...")
}

func downloadNewVersion(url, filepath string) error {
	// Создаём файл для сохранения новой версии
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Получаем данные с URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Копируем содержимое ответа в файл
	_, err = io.Copy(out, resp.Body)
	return err
}
