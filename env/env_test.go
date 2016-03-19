package env

import (
	"testing"
)

func TestInitEnv(t *testing.T) {
	defer func() {
		if v := recover(); v != nil {
			t.Errorf("Init env fail, %v\n", v)
		}
	}()

	InitEnv("../testing_env.json")
	Config() // don't care the value
}
