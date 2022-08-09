package actiondemux

import (
	"fmt"

	"github.com/Galdoba/devtools/cli/command"

	"github.com/Galdoba/ffstuff/app/demuxer/handle"

	"github.com/urfave/cli"
)

/*
ПРИМЕРЫ ПРИМЕНЕНИЯ
demuxer -tofile file.txt -update demux -i film.mp4
	-tofile file.txt - терминал будет писаться в указанный файл

Для демукса требуются:
1. Исходник(и)
2. Информация по заданию (данные из таблицы)
3. Задание (ввод в ручную что это фильм/трейлер/сериал)

ПЛАН:
1. Собираем данные:
	1.1 Подтверждаем исходник(и)
	1.2 Запрашиваем задание
	1.3 Читаем таблицу.
	1.4 ДЕБАГ: Выводим имена и пути предпологаемых результатов.


*/

func Run(c *cli.Context) error {
	fmt.Println("RUN Precheck")
	if err := Precheck(c); err != nil {
		return err
	}
	fmt.Println("Precheck complete")
	return nil
}

func Precheck(c *cli.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}
	for _, arg := range args {
		fmt.Println(arg)
	}
	selected := handle.SelectionSingle("Перечень исходных файлов корректен?", []string{"ДА", "НЕТ"}...)
	if selected != "ДА" {
		return fmt.Errorf("User abort")
	}
	for _, arg := range args {
		com, err := command.New(
			command.CommandLineArguments("ffprobe", "-i "+arg),
			command.Set(command.TERMINAL_ON),
		)
		if err != nil {
			return err
		}
		fmt.Println(" ")
		com.Run()
	}
	fmt.Println(" ")
	return nil
}
