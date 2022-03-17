package goo

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadConfig(yamlFile string, conf interface{}) (err error) {
	var buf []byte

	buf, err = ioutil.ReadFile(yamlFile)
	if err != nil {
		goo_log.WithTrace().Error(err.Error())
		return
	}

	if err = yaml.Unmarshal(buf, conf); err != nil {
		goo_log.WithTrace().Error(err.Error())
	}
	return
}
