// package main

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// // Счетчик времени
// var runTime int64

// // Обновление времени
// func updateTime() {
// 	startTime := time.Now()
// 	for {
// 		time.Sleep(time.Second)
// 		atomic.StoreInt64(&runTime, int64(time.Since(startTime).Seconds()))
// 	}
// }

// // HTTP сервер
// func httpServer() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Привет, мир!")
// 	})
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// // WebSocket сервер
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// // Хранилище всех подключений
// var clients = make(map[*websocket.Conn]bool)
// var mutex = sync.Mutex{}

// func wsServer() {
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			log.Println("Ошибка при обновлении:", err)
// 			return
// 		}
// 		defer conn.Close()

// 		// Добавляем новое соединение в список
// 		mutex.Lock()
// 		clients[conn] = true
// 		mutex.Unlock()

// 		// Удаляем соединение из списка при завершении
// 		defer func() {
// 			mutex.Lock()
// 			delete(clients, conn)
// 			mutex.Unlock()
// 		}()

// 		// Горутина для чтения сообщений от клиента
// 		go func() {
// 			for {
// 				_, message, err := conn.ReadMessage()
// 				if err != nil {
// 					log.Println("Ошибка при чтении сообщения:", err)
// 					return
// 				}

// 				// Вывод сообщения в терминал
// 				log.Printf("Сообщение от клиента: %s", message)

// 				// Отправка сообщения всем подключенным клиентам
// 				mutex.Lock()
// 				for client := range clients {
// 					if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
// 						log.Println("Ошибка при отправке сообщения клиенту:", err)
// 						client.Close()
// 						delete(clients, client)
// 					}
// 				}
// 				mutex.Unlock()
// 			}
// 		}()

// 		// Цикл для отправки случайных чисел клиентам
// 		for {
// 			randomNumber := rand.Intn(100)
// 			message := fmt.Sprintf("Случайное число: %d", randomNumber)
// 			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
// 				log.Println("Ошибка при отправке случайного числа:", err)
// 				return
// 			}
// 			time.Sleep(5 * time.Second)
// 		}
// 	})

// 	log.Fatal(http.ListenAndServe(":8081", nil))
// }

// func main() {
// 	// Запуск потоков
// 	go updateTime() // Подсчет времени работы приложения
// 	go httpServer() // HTTP сервер
// 	go wsServer()   // WebSocket сервер

// 	// Основной поток
// 	select {} // Блокируем основной поток для предотвращения завершения программы
// }

// // package main

// // import (
// // 	"fmt"
// // 	"syscall"
// // 	"unsafe"
// // )

// // var (
// // 	user32             = syscall.NewLazyDLL("user32.dll")
// // 	procEnumWindows    = user32.NewProc("EnumWindows")
// // 	procGetWindowTextW = user32.NewProc("GetWindowTextW")
// // 	procGetClassNameW  = user32.NewProc("GetClassNameW")
// // )

// // func EnumWindows(enumFunc uintptr, lParam uintptr) bool {
// // 	ret, _, _ := procEnumWindows.Call(
// // 		enumFunc,
// // 		lParam,
// // 	)
// // 	return ret != 0
// // }

// // func GetWindowText(hwnd uintptr, buf *uint16, nMaxCount int) int {
// // 	ret, _, _ := procGetWindowTextW.Call(
// // 		hwnd,
// // 		uintptr(unsafe.Pointer(buf)),
// // 		uintptr(nMaxCount),
// // 	)
// // 	return int(ret)
// // }

// // func GetClassName(hwnd uintptr, buf *uint16, nMaxCount int) int {
// // 	ret, _, _ := procGetClassNameW.Call(
// // 		hwnd,
// // 		uintptr(unsafe.Pointer(buf)),
// // 		uintptr(nMaxCount),
// // 	)
// // 	return int(ret)
// // }

// // func main() {
// // 	callback := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
// // 		var buf [256]uint16

// // 		// Получаем заголовок окна
// // 		length := GetWindowText(hwnd, &buf[0], int(len(buf)))
// // 		windowText := syscall.UTF16ToString(buf[:length])

// // 		// Получаем имя класса окна
// // 		length = GetClassName(hwnd, &buf[0], int(len(buf)))
// // 		className := syscall.UTF16ToString(buf[:length])

// // 		fmt.Printf("HWND: 0x%x, Class: %s, Title: %s\n", hwnd, className, windowText)

// // 		// Возвращаем 1, чтобы продолжить перечисление окон
// // 		return 1
// // 	})

// // 	EnumWindows(callback, 0)
// // }

package main

import (
	"fmt"
	"os"
)

func isValidPath(path string) bool {
	// Используем os.Stat для получения информации о файле или директории
	_, err := os.Stat(path)
	return !os.IsNotExist(err) // Если ошибка не "файл не существует", путь валиден
}

func main() {
	path := "C:/Windows/System32/setx.exe" // Укажи свой путь здесь

	if isValidPath(path) {
		fmt.Println("Путь валиден.")
	} else {
		fmt.Println("Путь не валиден.")
	}
}
