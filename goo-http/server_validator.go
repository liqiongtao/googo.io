package goohttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strings"
)

func ValidationMessage(err error, messages map[string]string) string {
	if err == nil {
		return ""
	}

	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		return fmt.Sprintf("请求参数 %s 的类型是 %s, 不是 %s", ute.Field, ute.Type, ute.Value)
	}

	var v validator.ValidationErrors
	if errors.As(err, &v) {
		for _, i := range v {
			field := goo_utils.Camel2Case(i.Field())
			key := fmt.Sprintf("%s_%s", field, strings.ToLower(i.Tag()))
			if msg, ok := messages[key]; ok {
				return msg
			}
			return fmt.Sprintf("%s %s", field, i.Tag())
		}
	}

	return err.Error()
}
