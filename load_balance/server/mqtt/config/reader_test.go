package config

import (
	"testing"
)

func TestConfigReader(t *testing.T) {
	reader := CReader{}
	reader.Read("test.cfg")
}
