package goo

import (
	"os"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"gopkg.in/yaml.v2"
)

func LoadConfig(yamlFile string, conf interface{}) (err error) {
	var buf []byte

	buf, err = os.ReadFile(yamlFile)
	if err != nil {
		goo_log.Error(err.Error())
		return
	}

	if err = yaml.Unmarshal(buf, conf); err != nil {
		goo_log.Error(err.Error())
	}
	return
}
