package goo_xlsx

import (
	"testing"
)

func TestXlsxWrite_Save2File(t *testing.T) {
	Writer().
		SetTitles([]string{"姓名", "年龄"}).
		SetData([]interface{}{"张三", 20}).
		SetData([]interface{}{"李四", 23}).
		Save2File("1.xlsx")
}
