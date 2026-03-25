package goo_xlsx

import (
	"fmt"
	"testing"
)

func TestXlsxWrite_Handler(t *testing.T) {
	w := Writer()

	w.SetSheetName("数据导出")

	w.SetMergeCellValue("A1", "D2", "大标题1")

	w.SetRowNum(3)
	w.SetMergeCellValue("A3", "B3", "子标题1")
	w.SetMergeCellValue("C3", "D3", "子标题2")

	left, right, err := w.SetTitles([]string{"姓名", "手机号"})
	fmt.Println(left, right, err)
	
	w.SetData([]interface{}{"李涛", "18510381580"})

	w.Save2File("user.xlsx")
}
