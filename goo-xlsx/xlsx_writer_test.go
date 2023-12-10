package goo_xlsx

import (
	"testing"
)

func TestXlsxWrite_Handler(t *testing.T) {
	w := Writer()

	w.SetSheetName("数据导出")
	w.SetTitles([]string{"姓名", "手机号"})
	w.SetData([]interface{}{"李涛", "18510381580"})

	w.Save2File("user.xlsx")
}
