package MelonePlayer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gonutz/w32/v2"
	"github.com/gorilla/websocket"
)

func CommandHandler(command string) interface{} {
	commands := strings.Split(command, ".")
	switch commands[0] {
	case "get":
		switch commands[1] {
		case "Default":
			switch commands[2] {
			case "PlayList":
				return_ := map[string]interface{}{
					"Music":   []string{},
					"Preview": []string{},
				}

				var Music []string
				var Preview []string

				var ret1 []string
				ret, err := CollectDirContents(Path.DefaultPlayList)
				if err == nil {
					for i, str := range ret {
						if strings.HasSuffix(str, "Music") {
							ret1 = Remove(ret, i)
						} else if strings.HasSuffix(str, "Preview") {
							ret1 = Remove(ret, i)
						}
					}

					for i, str := range ret1 {
						ret1[i] = strings.ReplaceAll(str, Path.StaticDir, "")
					}

					for _, item := range ret1 {
						if strings.Contains(item, "Music") {
							Music = append(Music, item)
						} else if strings.Contains(item, "Preview") {
							Preview = append(Preview, item)
						}
					}
					return_["Music"] = Music
					return_["Preview"] = Preview
				}
				return return_
			}
		}
	}

	Log.Warn(fmt.Sprintf(command))
	return command
}

// WebSocket сервер
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Хранилище всех подключений
var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

func fileServerHandler(w http.ResponseWriter, r *http.Request) {
	// Определяем директорию в зависимости от пути запроса
	var dir string
	if strings.HasPrefix(r.URL.Path, "/static1/") {
		dir = "./static1"
		http.StripPrefix("/static1/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/static2/") {
		dir = "./static2"
		http.StripPrefix("/static2/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func substitute(htmlFilePath string) (string, error) {
	// Чтение HTML-файла
	htmlContent, err := os.Open(htmlFilePath)
	if err != nil {
		return "", err
	}
	defer htmlContent.Close() // Закрываем файл после завершения работы
	// Создаем срез для хранения строк
	var lines []string

	// Создаем сканер для чтения файла
	scanner := bufio.NewScanner(htmlContent)

	// Чтение файла построчно
	for scanner.Scan() {
		line := scanner.Text()      // Получаем строку
		lines = append(lines, line) // Добавляем строку в срез
	}

	// Проверка на ошибки
	if err := scanner.Err(); err != nil {
		Log.Error("Error reading file: %v\n", err)
		return "", err
	}

	for index, item := range lines {
		if strings.Contains(item, "{{") && strings.Contains(item, "}}") {
			Log.Info(item)
			start := strings.Index(item, "{{")
			end := strings.Index(item, "}}")
			if Translation.Get(Language_Default+"."+item[start+2:end]) != nil {
				item = strings.ReplaceAll(item, item[start:end+2], Translation.Get(Language_Default+"."+item[start+2:end]).(string))
				lines[index] = item
			} else {
				Log.Error("There is no translation for: " + item[start+2:end])
			}
			Log.Info(item)

		}
	}
	var result string

	// Соединяем строки из среза
	for _, line := range lines {
		result += line + "\n" // Добавляем строку и перевод строки
	}
	return result, nil
}

func Server() {
	// Обработка статических файлов
	fs_0 := http.FileServer(http.Dir(Path.HomePath + "\\" + Path.Themes + "\\" + string(App_Settings.GetStringBytes("theme"))))
	fs_1 := http.FileServer(http.Dir(Path.HomePath + "\\" + Path.PlayLists))
	//
	//
	http.Handle("/", fs_0)
	http.Handle("/PlayLists/", http.StripPrefix("/PlayLists/", fs_1))
	//
	//
	go func() {
		// Обработка "/player.window" и возврат конкретного HTML-файла
		http.HandleFunc("/player.window", func(w http.ResponseWriter, r *http.Request) {
			// Путь к HTML-файлу
			theme := string(App_Settings.GetStringBytes("theme"))
			htmlFilePath := Path.Themes + "/" + theme + "/html/index.html"
			result, err := substitute(htmlFilePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				Log.Error(err)
			} else {
				// Установка заголовков ответа
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte(result))
			}

		})
	}()
	//
	//
	go func() {
		// Обработка "/settings.window" и возврат конкретного HTML-файла
		http.HandleFunc("/settings.window", func(w http.ResponseWriter, r *http.Request) {
			// Путь к HTML-файлу
			theme := string(App_Settings.GetStringBytes("theme"))
			htmlFilePath := Path.Themes + "/" + theme + "/html/settings.html"
			result, err := substitute(htmlFilePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				Log.Error(err)
			} else {
				// Установка заголовков ответа
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte(result))
			}
		})
	}()

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		// Получение параметра "command" из URL
		command := r.URL.Query().Get("command")
		if command == "" {
			http.Error(w, "Параметр 'command' не найден", http.StatusBadRequest)
			Log.Warn("Параметр 'command' не найден: %v", http.StatusBadRequest)
			return
		}

		value := CommandHandler(command)
		switch value.(type) {
		case map[string]interface{}:
			jsonBytes, err := json.Marshal(value)
			if err != nil {
				http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
				Log.Error(err.Error())
				return
			} else if err == nil {
				// Устанавливаем заголовок Content-Type
				w.Header().Set("Content-Type", "application/json")

				// Записываем JSON в ответ
				w.Write(jsonBytes)

			}
		}

	})
	//
	//
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Ошибка при обновлении:", err)
			return
		}
		defer conn.Close()

		// Добавляем новое соединение в список
		mutex.Lock()
		clients[conn] = true
		mutex.Unlock()

		// Удаляем соединение из списка при завершении
		defer func() {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
		}()

		// Горутина для чтения сообщений от клиента
		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					message = []byte(err.Error())
					Log.Error("Ошибка при чтении сообщения:", err)
					return
				}

				commands := strings.Split(string(message), ".")
				switch commands[0] {
				case "Player":
					switch commands[1] {
					case "Hide":
						w32.ShowWindow(MainWindowSettings.hwnd, w32.SW_HIDE)
						message = []byte("ok")
					case "Show":
						w32.ShowWindow(MainWindowSettings.hwnd, w32.SW_SHOW)
						message = []byte("ok")
					}
				case "SettingsW":
					switch commands[1] {
					case "Hide":
						w32.ShowWindow(SettingsWindows.hwnd, w32.SW_HIDE)
						message = []byte("ok")
					case "Show":
						w32.ShowWindow(SettingsWindows.hwnd, w32.SW_SHOW)
						message = []byte("ok")
					}
				case "Download":
					url, err := ExtractFromBrackets(string(message))
					if err == nil {
						Log.Info(url)
						Download(url[0], commands[1])
					}
				default:
					message = []byte("404")
				}

				go func() {
					for {
						msg := <-Message
						if len(msg) != 0 {
							// Отправка сообщения всем подключенным клиентам
							mutex.Lock()
							for client := range clients {
								if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
									Log.Error("Ошибка при отправке сообщения клиенту:", err)
									client.Close()
									delete(clients, client)
								}
							}
							mutex.Unlock()
							time.Sleep(100 * time.Millisecond)
						}
					}
				}()

				// Отправка сообщения всем подключенным клиентам
				mutex.Lock()
				for client := range clients {
					if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
						Log.Error("Ошибка при отправке сообщения клиенту:", err)
						client.Close()
						delete(clients, client)
					}
				}
				mutex.Unlock()
			}
		}()

		for {
			time.Sleep(10 * time.Millisecond) // Это сделает ваш цикл менее агрессивным
			continue
		}
	})

	Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil).Error())
}
