package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Galdoba/ffstuff/pkg/config"
	"gopkg.in/yaml.v3"
)

type grabberConfig struct {
	Description       string
	External_Log_path string
	Local_Log_path    string
	//Actions
	Actions []Action
}

/*
CONFIG:
пути:
SOURCE_ROOT - рутовая папка в которой искать ready файлы
DEFAULT_DESTINATION - путь куда летят файлы, если неуказано иного (обязательно)
TASK_DIRECTORY - путь куда будут лететь отложенные задачи для перемещений файлов (не обязательно)
LOG_FILEPATH - общий файл куда пишутся логи (обязательно)
LOG_LEVEL - минимальный уровень важности сообщений
LOG_BY_SESSION - true/false - кроме общего файла лог пишется в отдельный файл для каждого запуска программы
поведение:
EXIT_WHEN_DONE - true/false - отключиться когда нечего качать
TRIGGER_SCHEDULE - true/false - запускать граббер по расписанию?
SCHEDULE - расписание для автозапуска по кроновской схеме
TRIGGER_TIMEOUT - true/false - запускать граббер по КД (каждые Х секунд)?
TIMEOUT - собственно Х секунд
порядок закачек:
SORT_ORDER - перечень тегов влияющий на порядок скачиваний
GRAB_BY_SIZE - true/false - качать сначала маленькие файлы (конфликтует с SORT_ORDER)
COPY_PREFIX - маска для копий "copy" для file.mp4 ==> copy_1_file.mp4 //индекс увеличивается по количеству копий
COPY_SUFFIX - маска для копий "copy" для file.mp4 ==> file_copy_1.mp4
оповещения (на перспективу):
tgChannel - в какой чат плевать текст
system    - выскакивающее уведомление на локальной машине

стандартные уровни лога - можно добавить свои:
	FATAL  - ошибка приводящая к завершению программы (нет отлика от удаленного компа/аномальные данные)
	ERROR  - ошибка приводящая к завершению операции (не могу скачать файл - оборвалась связь)
	WARN   - незначительная ошибка которая ни к чему не ведет или которую решает сама программа (файл уже скачан - скипуем/переписываем/переименовываем)
	INFO   - начало/окончание закачки файлового сета/операции
	DEBUG  - начало/окончание закачки конкретного файла/отдельного действия (переименование/удаление/решение по алгоритму)
	TRACE  - текущее состояние системы со всеми аргументами (очень много словно чтобы разобраться что происходит)

флаги:
-dest, -d       - папка куда качать файлы (приоритет над конфигом)
-overwrite, -ow - если файл уже есть, перезаписываем
-rename, -r     - если файл уже есть, переименовываем
-move, -mv      - удаляем исходник после копирования
-loglevel       - уровень сообщений для консоли (по умолчанию равен LOG_LEVEL для логфайла в конфиге)



команды:
setup  - настроить конфиг
health - проверить програмные файлы
search - искать ready файлы и принтовать результаты
grab   - переместить файлы по ключам/аргументам и закончить сессию
queue  - создать отложенную задачу на перемещение файлов
run    - запустить грабер (будет работать по тригерам и отложенным задачам)
help   - чтиво на ночь
*/

type Action struct {
	ActionName string
	Triggers   []string
}

func CreateDefaultConfig() error {
	gc := grabberConfig{}
	gc.Description = "config file for 'grabber.exe'"

	gc.External_Log_path = "TODO"
	gc.Local_Log_path = "TODO"
	gc.Actions = []Action{
		{
			ActionName: "MOVE_CURSOR_UP",
			Triggers:   []string{"UP"},
		},
		{
			ActionName: "MOVE_CURSOR_PU",
			Triggers:   []string{"PgUp"},
		},
		{
			ActionName: "MOVE_CURSOR_TOP",
			Triggers:   []string{"HOME"},
		},
		{
			ActionName: "MOVE_CURSOR_DOWN",
			Triggers:   []string{"DOWN"},
		},
		{
			ActionName: "MOVE_CURSOR_PD",
			Triggers:   []string{"PgDn"},
		},
		{
			ActionName: "MOVE_CURSOR_BOTTOM",
			Triggers:   []string{"END"},
		},
		{
			ActionName: "CURSOR_DOWN_AND_TOGGLE_SELECTION",
			Triggers:   []string{"Insert"},
		},
		{
			ActionName: "TOGGLE_SELECTION_STATE",
			Triggers:   []string{"SPACE"},
		},
		{
			ActionName: "SELECT_ALL_WITH_SAME_EXTENTION",
			Triggers:   []string{"Ctrl+SPACE"},
		},
		{
			ActionName: "DROP_SELECTIONS",
			Triggers:   []string{"~", "BACKSPACE"},
		},
		{
			ActionName: "MOVE_SELECTED_TOP",
			Triggers:   []string{"ENTER", "Ctrl+T"},
		},
		{
			ActionName: "MOVE_SELECTED_BOTTOM",
			Triggers:   []string{"Ctrl+B"},
		},
		{
			ActionName: "MOVE_SELECTED_UP",
			Triggers:   []string{"W"},
		},
		{
			ActionName: "MOVE_SELECTED_DOWN",
			Triggers:   []string{"S"},
		},
		{
			ActionName: "DECIDION_CONFIRM",
			Triggers:   []string{"ENTER"},
		},
		{
			ActionName: "DELETE_SELECTED",
			Triggers:   []string{"Delete"},
		},
		{
			ActionName: "DECIDION_DENY",
			Triggers:   []string{"~", "BACKSPACE"},
		},
		{
			ActionName: "DOWNLOAD_PAUSE",
			Triggers:   []string{"P"},
		},
		{
			ActionName: "UNDO_MOVEMENT",
			Triggers:   []string{"Ctrl+Z"},
		},
		{
			ActionName: "ADD_NEW_SOURCE_FROM_CLIPBOARD",
			Triggers:   []string{"Ctrl+V"},
		},
		{
			ActionName: "ACTION_QUIT_PROGRAM",
			Triggers:   []string{"Ctrl+Q"},
		},
	}
	fileBts, err := yaml.Marshal(gc)
	if err != nil {
		return err
	}

	cDir, cFile := config.StdConfigPath(programName)
	confPath := fmt.Sprintf("%v%v", cDir, cFile)
	fmt.Println("will go here:", confPath)
	//confPath, err = filepath.Abs(confPath)

	if err := os.MkdirAll(cDir, 0777); err != nil {
		return err
	}
	fmt.Printf("'%v' created\n", cDir)
	//panic(confPath)
	// read the whole file at once
	_, err = os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("'%v' opened\n", confPath)
	// write the whole body at once
	err = ioutil.WriteFile(confPath, fileBts, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadConfig(path string) (*grabberConfig, error) {
	fl, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	gc := &grabberConfig{}
	err = yaml.Unmarshal(fl, gc)
	if err != nil {
		return nil, err
	}
	if gc.Validate() != nil {
		return nil, err
	}

	return gc, nil
}

func (gc *grabberConfig) Validate() error {
	if len(gc.Actions) < 5 {
		return fmt.Errorf("")
	}
	return nil
}
