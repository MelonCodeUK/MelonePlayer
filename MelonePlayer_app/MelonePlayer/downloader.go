package MelonePlayer

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/valyala/fastjson"
)

func download(path string, args []string) {
	// Команда и её аргументы
	cmd := exec.Command(path, args...)

	// Создание пайпов для чтения Stdout и Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Log.Error("Ошибка получения StdoutPipe: %v\n", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		Log.Error("Ошибка получения StderrPipe: %v\n", err)
		return
	}

	// Запуск команды
	if err := cmd.Start(); err != nil {
		Log.Error("Ошибка запуска команды: %v\n", err)
		return
	}

	// Горутина для чтения вывода из Stdout в реальном времени
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			t_ := strings.Split(RemoveDuplicate(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(scanner.Text(), "of", ""), "at", ""), "ETA", ""), "~", ""), "frag", ""), "(", ""), ")", ""), "in", ""), " "), " ")
			if strings.Contains(scanner.Text(), "[download]") {
				Message <- strings.Join(t_, " ") + "\n"
			} else {
				Message <- scanner.Text() + "\n"
			}
			Log.Info(path + ": " + scanner.Text())
		}
		Message <- ""
	}()

	// Горутина для чтения вывода из Stderr в реальном времени
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			Log.Error(path + ":" + scanner.Text())
		}
	}()

	// Ожидание завершения процесса
	if err := cmd.Wait(); err != nil {
		Log.Error("Команда завершилась с ошибкой: %v\n", err)
		return
	} else {
		Message <- ""

	}

}

func Download(url string, type_ string) {
	downloads_services := Data.Get("lib_settings").Get("binaries")

	obj, err := downloads_services.Object()
	if err == nil {
		obj.Visit(func(key []byte, value *fastjson.Value) {
			for _, item := range value.GetArray("support_links") {
				if strings.Contains(url, string(item.GetStringBytes())) {
					c_temp, err := value.Get("download_commands").Object()
					if err == nil {
						c_temp.Visit(func(key1 []byte, value1 *fastjson.Value) {
							if string(key1) == type_ {
								var command []string
								err := json.Unmarshal(value1.MarshalTo(nil), &command)
								if err == nil {
									command[0] = Path.Bin + "/" + command[0]
									for index, item := range command {
										if strings.Contains(item, "{path}") {
											command[index] = Path.DefaultPlayList + "/" + strings.ToUpper(string(type_[0])) + strings.ToLower(type_[1:]) + command[index][6:]
										} else if item == "{url}" {
											command[index] = url
										}
									}
									if len(command) != 0 {
										Log.Info(command)
										Log.Info(command[0])
										Log.Info(command[1:])
										download(command[0], command[1:])
									}
								} else {
									Log.Error(err.Error())
								}
							}
						})
					} else {
						Log.Error(err.Error())
					}

				}
			}

		})
	}

}
