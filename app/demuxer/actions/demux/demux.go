package actiondemux

import (
	"fmt"

	"github.com/urfave/cli"
)

/*
ПРИМЕРЫ ПРИМЕНЕНИЯ
demuxer -tofile file.txt -update demux -i film.mp4
	-tofile file.txt - терминал будет писаться в указанный файл

*/

func Run(c *cli.Context) error {
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

	return nil
}
