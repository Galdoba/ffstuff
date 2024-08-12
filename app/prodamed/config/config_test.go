package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := Load(defaultConfigPath())
	fmt.Println(err)
	fmt.Println(cfg)

}
