package cmd

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/urfave/cli/v2"
)

const (
	flag_Extended = "extended"
)

func Config() *cli.Command {
	cm := &cli.Command{
		Name:      "config",
		Usage:     "print config file",
		UsageText: "mfline config",
		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(*cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			cfg, _ := config.Load(c.App.Name)
			fmt.Println(cfg.String())
			if c.Bool(flag_Extended) {
				fmt.Printf(strings.Join([]string{
					"%v config fields explanations:",
					"Default Storage Directory   - тут хранятся результаты сканов",
					"Log File                    - пишем лог сюда",
					"Old Scan Age (hours)        - если файл с данными старше этого количества часов",
					"                              считаем его старым",
					"Auto Delete Old Scans       - если true - удаляем старые сканы по завершению работы",
					"Scan In Tracked Directories - запускаем сканирование файлов в выделенных директориях",
					"Repeat Scans if Error       - запускаем повторные сканы если произошла ошибка",
					"Track Directories           - перечень директорий за которыми следим",
				}, "\n"), c.App.Name)
			}
			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return err
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        flag_Extended,
				Category:    "",
				DefaultText: "",
				FilePath:    "",
				Usage:       "print extended explanation on each field",
				Required:    false,
				Hidden:      false,
				HasBeenSet:  false,
				Aliases:     []string{},
				EnvVars:     []string{},
			},
		},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
	return cm
}
