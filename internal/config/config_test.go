package config

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	LoadConfig("test.toml")
	fmt.Println(Config)
}
