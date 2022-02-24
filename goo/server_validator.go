package goo

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strings"
)

func ValidationMessage(err error, msgs map[string]string) string {
	if v, ok := err.(*json.UnmarshalTypeError); ok {
		return fmt.Sprintf("请求参数 %s 的类型是 %s, 不是 %s", v.Field, v.Type, v.Value)
	}

	if v, ok := err.(validator.ValidationErrors); ok {
		for _, i := range v {
			field := goo_utils.Camel2Case(i.Field())
			key := fmt.Sprintf("%s_%s", field, strings.ToLower(i.Tag()))
			if msg, ok := msgs[key]; ok {
				return msg
			}
			msg := fmt.Sprintf("%s %s", field, i.Tag())
			return msg
		}
	}

	return err.Error()
}
