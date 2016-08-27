package helper

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// JSON2Yaml _
func JSON2Yaml(b []byte) ([]byte, error) {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(&m)
}

// JSON2YamlFile _
func JSON2YamlFile(src, dest string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	b, err = JSON2Yaml(b)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, b, 0644)
}
