package goo

import (
	"os"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	"gopkg.in/yaml.v2"
)

func LoadConfig(yamlFile string, conf any) (err error) {
	var buf []byte

	buf, err = os.ReadFile(yamlFile)
	if err != nil {
		goolog.Error(err.Error())
		return
	}

	if err = yaml.Unmarshal(buf, conf); err != nil {
		goolog.Error(err.Error())
	}
	return
}
