package cmd

import (
	"github.com/urfave/cli/v2"
)

func Listen() *cli.Command {

	cm := &cli.Command{
		Name: "listen",
		Action: func(c *cli.Context) error {

			return nil
		},
		Flags: []cli.Flag{},
	}
	return cm
}
