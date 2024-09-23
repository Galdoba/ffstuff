package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

var dest string

var cfg config.Configuration

func Grab() *cli.Command {
	return &cli.Command{
		Name:        "grab",
		Aliases:     []string{},
		Usage:       "TODO: Direct command for transfering operation(s)",
		UsageText:   "",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",
		Before: func(*cli.Context) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if validation.FileValidation(cfg.LOG) != nil {
				return fmt.Errorf("log filepath in config: %v")
			}

			switch cfg.LOG_LEVEL {
			case "":
				return fmt.Errorf("поле 'Minimum Log Level' не содержит данных")
			default:
				return fmt.Errorf("поле 'Minimum Log Level' содержит некоректные данные: '%v'\n"+
					"ожидаю: TRACE, DEBUG, INFO, WARN, ERROR или FATAL", cfg.LOG_LEVEL)
			case strings.ToUpper(logman.TRACE):
				logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceTRACE))
			case strings.ToUpper(logman.DEBUG):
				logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceDEBUG))
			case strings.ToUpper(logman.INFO):
				logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceINFO))
			case strings.ToUpper(logman.WARN):
				logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceWARN))
			case strings.ToUpper(logman.ERROR):
				logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceERROR))
			case strings.ToUpper(logman.FATAL):
				logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceFATAL))
			}
			logman.ClearOutput(logman.ALL)
			logman.SetOutput(cfg.LOG, logman.ALL)
			logman.SetOutput(logman.Stderr, logman.ALL)
			if err := validation.DirectoryValidation(cfg.DEFAULT_DESTINATION); err != nil {
				return fmt.Errorf("config defaut destination invalid: %v", err)
			}
			dest = cfg.DEFAULT_DESTINATION
			fmt.Println(cfg)
			return nil
		},
		Action: func(c *cli.Context) error {

			logman.Debug(logman.NewMessage("start 'grub'"))
			time.Sleep(1 * time.Second)
			switch c.String("dest") {
			case "":

			}
			if c.String("dest") != "" {
				dest = c.String("dest")

			}

			logman.Debug(logman.NewMessage("end   'grub'"))
			return nil
		},

		Subcommands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dest",
				Usage:       "destination where files will be downloaded to\n",
				DefaultText: "from config",
				Required:    false,
				Hidden:      false,
				HasBeenSet:  false,
				Value:       "",
				Destination: new(string),
				Aliases:     []string{"d"},
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
}
