package MelonePlayer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/gonutz/w32/v2"
	"github.com/valyala/fastjson"
)

type JsonNode map[string]interface{}

// settings
var (
	Data             *fastjson.Value
	Scripts          *fastjson.Value
	StartScript      *fastjson.Value
	App_Settings     *fastjson.Value
	App_Info         *fastjson.Value
	Version          float32
	Version_Settings float32
	Port             int
	Server_Url       string
	Languages        []string
	Language_Default string
	IsUpdate         bool
	IsSave_Log       bool = true
	IsDebug          bool
	Resetting        bool
	// IsSettingsWindowVisible bool = false
)
var Message = make(chan string)

type Window struct {
	Title            string
	hwnd             w32.HWND
	Width            int
	Height           int
	X                int
	Y                int
	WS_SYSMENU       bool // Управляет наличием системного меню (и кнопки закрытия).
	WS_MINIMIZEBOX   bool // Управляет наличием кнопки минимизации.
	WS_MAXIMIZEBOX   bool // Управляет наличием кнопки максимизации.
	WS_CAPTION       bool // Добавляет или убирает заголовок окна.
	WS_THICKFRAME    bool // Управляет возможностью изменения размера окна.
	WS_EX_TOOLWINDOW bool // Устанавливает наличие значка в панали задач.
	WS_EX_TOPMOST    bool // Параметр который устанавливает окно по верх других
	SWP_NOZORDER     bool // Oкно должно менять свою позицию в Z-последовательности
	SWP_NOMOVE       bool // Окно нельзя перемещать по кординатам X и Y
}

const (
	ABS_AUTOHIDE = 0x0000001
)

// Абстракция структуры APPBARDATA
type APPBARDATA struct {
	cbSize           uint32
	hWnd             syscall.Handle
	uCallbackMessage uint32
	uEdge            uint32
	rc               w32.RECT
	lParam           int32
}

var MainWindowSettings Window
var SettingsWindows Window

type Paths struct {
	StaticDir        string
	HomePath         string
	SettingsFilePath string
	Lib              string
	Bin              string
	PlayLists        string
	DefaultPlayList  string
	Themes           string
}

var Path Paths
var now = time.Now()
var Log = NewLogger(fmt.Sprintf("log_%d.%s.%d.log", now.Year(), now.Month(), now.Day()))
var settingsNames = []string{"setting.json", ".setting.json", "settings.json", ".settings.json", "настройки.json", ".настройки.json", "настройка.json", ".настройка.json", "налаштування.json", ".налаштування.json"}

func Command(command string, args []string) (string, error) {
	// Выполнение команды
	cmd := exec.Command(command, args...)
	output, err := cmd.Output() // Получаем вывод команды

	if err != nil {
		Log.Error(err.Error())
		return "", err
	} else {
		Log.Info(string(output))
		return string(output), nil
	}
}

// Функция для рекурсивного обхода файлов и папок и сбора путей в срез
func CollectDirContents(path string) ([]string, error) {
	// Если путь равен ".", получаем путь к текущему рабочему каталогу
	if path == "." {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	var paths []string

	// Открыть директорию
	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	// Пройтись по всем файлам и папкам в директории
	for _, file := range files {
		// Формируем полный путь
		fullPath := filepath.Join(path, file.Name())

		if file.IsDir() {
			// Если это директория, добавляем её путь в срез и рекурсивно обходим
			paths = append(paths, fullPath)
			subPaths, err := CollectDirContents(fullPath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		} else {
			// Если это файл, добавляем его путь в срез
			paths = append(paths, fullPath)
		}
	}

	return paths, nil
}

func RemoveDuplicate(s string, w string) string {
	// Split the string by spaces
	words := strings.Fields(s)
	// Join the words with a single space
	return strings.Join(words, w)
}

// функиця для поиска файла с настройками
func SearchSettingsFiles() string {
	Log.Info("search for a configuration file....")
	paths, _ := CollectDirContents(".")
	for _, path := range paths {
		for _, settingsPath := range settingsNames {
			if strings.HasSuffix(path, settingsPath) {
				Log.Info("settings file found!")
				return path
			}
		}
		Log.Info(path, fmt.Sprintf("%*s%s", 10, "", color.RedString("NO")))
	}
	return ""
}

func StartStartingsScripts() {

}

// Чтение и парсинг файла JSON

func GetData(keys []string) *fastjson.Value {
	Path.SettingsFilePath = SearchSettingsFiles()
	if Path.SettingsFilePath != "" {
		fileData, err := os.ReadFile(Path.SettingsFilePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}

		var p fastjson.Parser
		v, err := p.Parse(string(fileData))
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return nil
		}
		if len(keys) == 0 {
			return v
		}

		currentValue := v
		for _, key := range keys {
			if currentValue = currentValue.Get(key); currentValue == nil {
				fmt.Printf("Key %s not found\n", key)
				return nil
			}
		}
		return currentValue
	}
	return nil
}

func GetSettings() {
	Path.SettingsFilePath = SearchSettingsFiles()
	if Path.SettingsFilePath != "" {
		Data = GetData([]string{})
		App_Info = Data.Get("app_info")
		App_Settings = Data.Get("app_settings")
		go func() {
			dir, err := os.Getwd() // Получаем текущую директорию
			if err != nil {
				Log.Error(err.Error())
				return
			}
			Path.HomePath = dir
		}()

		Scripts := Data.GetArray("scripts")
		for _, item := range Scripts {
			if string(item.GetStringBytes("name")) == "SetHomePathVariable" {
				value, err := Command("where", []string{string(item.GetStringBytes("type"))})
				if err == nil {
					var a_temp = item.GetArray("script")
					// Преобразуем []*fastjson.Value в []string
					var result []string
					for _, v := range a_temp {
						if string(v.GetStringBytes()) == "\"{path}\"" {
							b_temp := strings.ReplaceAll(string(v.GetStringBytes()), "{path}", Path.HomePath)
							Log.Error(b_temp)
							result = append(result, b_temp) // Извлекаем строковое значение

						} else {
							result = append(result, string(v.GetStringBytes())) // Извлекаем строковое значение

						}
					}
					Command(value[:len(value)-2], result)
				}
			}

		}

		Path.StaticDir = string(App_Settings.GetStringBytes("static_dir"))
		Version = float32(App_Info.GetFloat64("version"))
		Version_Settings = float32(App_Info.GetFloat64("settings_version"))
		Port = App_Settings.GetInt("port")
		Server_Url = string(App_Settings.GetStringBytes("url_server"))
		MainWindowSettings.Width = App_Settings.GetInt("width")
		MainWindowSettings.Height = App_Settings.GetInt("height")
		MainWindowSettings.X = App_Settings.GetInt("x")
		MainWindowSettings.Y = App_Settings.GetInt("y")
		IsSave_Log = App_Settings.GetBool("save_log")
		IsUpdate = App_Settings.GetBool("auto_update")
		IsDebug = App_Settings.GetBool("debug")
		MainWindowSettings.Title = string(App_Settings.GetStringBytes("title_app"))

		lan := App_Info.GetArray("localization")
		for _, language := range lan {
			Languages = append(Languages, language.String())

		}
		Language_Default = string(App_Settings.GetStringBytes("default_language"))
		checkingСomponents := CheckingСomponents()
		if checkingСomponents != nil {
			os.Exit(1)
		}
		// resetting()
	}
}

func half(str string) int {
	if len(str) != 0 {
		return len(str) / 2
	}
	return 0
}

func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func CheckingСomponents() error {
	// paths
	paths_files := []string{}
	TraverseDirectories(App_Settings.Get("paths").Get("{static_dir}"), "{static_dir}", &paths_files)
	// Проходим по каждому элементу списка и заменяем слово
	for i, str := range paths_files {
		paths_files[i] = strings.ReplaceAll(str, "{static_dir}", Path.StaticDir)
	}
	for _, path := range paths_files {
		// Проверяем, существует ли папка
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Папка не существует, создаем её
			err := os.MkdirAll(path, 0755)
			if err != nil {
				Log.Fatal(fmt.Sprintf("Folder creation check: " + path + "\n\tERROR: " + err.Error()))
				return err
			}
			Log.Debug(fmt.Sprintf("Folder creation check: " + path + "\n\tCREATED!"))
		} else {
			Log.Info(fmt.Sprintf("Folder creation check: " + path + "\n\tEXISTS!"))
		}

		if strings.HasSuffix(path, "lib") {
			Path.Lib = path
		} else if strings.HasSuffix(path, "PlayLists") {
			Path.PlayLists = path
		} else if strings.HasSuffix(path, "Themes") {
			Path.Themes = path
		} else if strings.HasSuffix(path, "DefaultPlayList") {
			Path.DefaultPlayList = path
		} else if strings.HasSuffix(path, "bin") {
			Path.Bin = path
		}

	}
	return nil

}

func PrintInfo() {

	name := string(App_Info.GetStringBytes("name"))
	name1 := name[:half(name)]
	name2 := name[half(name):]

	Log.Write(
		color.GreenString(name1),
		color.RedString(name2),
		color.BlackString(" v."),
		Version,
		color.RedString(string(App_Info.GetStringBytes("version_type"))), " (",
		color.HiBlackString(MainWindowSettings.Title), ")\n",
	)
	Log.Write(color.WhiteString("Settings version: ") + color.HiCyanString("%.1f\n", Version_Settings))
	Log.Write(
		color.RedString("Width") +
			" and " +
			color.HiBlueString("Height") +
			": " +
			color.RedString("%d", MainWindowSettings.Width) +
			" " +
			color.HiBlueString("%d\n", MainWindowSettings.Height),
	)
	Log.Write(color.BlueString("Server url: %s:%d\n", Server_Url, Port))
	settingsfile := Path.SettingsFilePath
	if settingsfile == "" {
		settingsfile = color.RedString("SETTINGS FILE NOT FOUND!")
		Log.Fatal("SETTINGS FILE NOT FOUND!")
	}
	Log.Write(color.GreenString("Settings file path: ") + settingsfile + "\n")
	Static_Dirr, err := os.Getwd()
	if err != nil {
	}
	Log.Write(color.CyanString("Working directory: ") + Static_Dirr + "\\" + Path.StaticDir + "\n")
	Log.Write(color.HiYellowString("Save a log file: "))
	if IsSave_Log == true {
		Log.Write(color.HiGreenString("YES\n"))

	} else if IsSave_Log == false {
		Log.Write(color.HiRedString("NO\n"))

	}
	Log.Write(color.HiYellowString("Auto update: "))

	if IsUpdate == true {
		Log.Write(color.HiGreenString("YES\n"))

	} else if IsUpdate == false {
		Log.Write(color.HiRedString("NO\n"))

	}
	Log.Write(color.HiYellowString("Debug mode: "))
	if IsDebug == true {
		Log.Write(color.HiRedString("YES\n"))

	} else if IsDebug == false {
		Log.Write(color.HiGreenString("NO\n"))

	}
	temp_string := strings.Join(Languages, ", ")
	Log.Write(color.New(color.BgHiCyan).Sprint(
		color.CyanString("Languages supported: %s", color.New(color.BgGreen).Sprint(temp_string)) + "\n",
	))
	Log.Write(color.BlueString("Default language: %s", color.HiGreenString(Language_Default)) + "\n")

}

// Функция для обхода ключей и значений в fastjson.Value и возвращения массива путей
func TraverseDirectories(node *fastjson.Value, path string, paths *[]string) {
	/*
	   EXAMPLE:
	   ##Input data:
	   {
	   	"paths": {
	   		"./assets": {
	   			"/lib": ["/bin"],
	   			"/files": [
	   				{ "/PlayList": ["/Music", "/Previews", "/Video"] },
	   				"/Themes"
	   			]
	   		}
	   	}
	   }
	   ##Output data:
	   [./assets ./assets/lib
	   ./assets/lib/bin
	   ./assets/files
	   ./assets/files/PlayList
	   ./assets/files/PlayList/Music
	   ./assets/files/PlayList/Previews
	   ./assets/files/PlayList/Video
	   ./assets/files/Themes]
	*/
	switch node.Type() {
	case fastjson.TypeObject:
		obj, err := node.Object()
		if err != nil {
			fmt.Println("Error parsing object:", err)
			return
		}
		obj.Visit(func(k []byte, v *fastjson.Value) {
			newPath := path
			if len(path) > 0 {
				newPath += string(k)
			} else {
				newPath = string(k)
			}
			*paths = append(*paths, newPath)
			TraverseDirectories(v, newPath, paths)
		})
	case fastjson.TypeArray:
		for _, item := range node.GetArray() {
			TraverseDirectories(item, path, paths)
		}
	default:
		// Преобразуем []byte в string перед добавлением к path
		if len(path) > 0 {
			path += string(node.GetStringBytes())
		} else {
			path = string(node.GetStringBytes())
		}
		*paths = append(*paths, path)
	}
}

// Функция для извлечения подстрок из круглых скобок
func ExtractFromBrackets(input string) ([]string, error) {
	// Регулярное выражение для поиска текста в круглых скобках
	re := regexp.MustCompile(`\(([^)]+)\)`)

	// Поиск всех совпадений
	matches := re.FindAllStringSubmatch(input, -1)

	// Создание среза для хранения результатов
	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}

	return results, nil
}

func resetting() {
	screenWidth, screenHeight := GetDisplayResolution()
	MainWindowSettings.X = screenWidth - MainWindowSettings.Width - 1
	MainWindowSettings.Y = screenHeight - MainWindowSettings.Height - 1
	SetPosition(MainWindowSettings)
	App_Settings.Set("resetting", fastjson.MustParse(`false`))
	App_Settings.Set("x", fastjson.MustParse(fmt.Sprintf(`%s`, fmt.Sprintf("%d", MainWindowSettings.X))))
	App_Settings.Set("y", fastjson.MustParse(fmt.Sprintf(`%s`, fmt.Sprintf("%d", MainWindowSettings.Y))))

	// Получение обновленных данных JSON
	updatedData := Data.MarshalTo(nil)
	// Запись обновленных данных в файл
	if err := os.WriteFile(Path.SettingsFilePath, updatedData, 0644); err != nil {
		Log.Fatal(fmt.Sprintf("Ошибка при записи файла: %s", err.Error()))
	}

}
