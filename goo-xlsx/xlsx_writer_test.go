package goo_xlsx

import (
	"testing"
)

func TestXlsxWrite_Handler(t *testing.T) {
	w := Writer()

	w.SetSheetName("数据导出")

	w.SetMergeCellValue("A1", "D2", "大标题1", defaultCellStyle)

	w.SetRowNum(3)
	w.SetMergeCellValue("A3", "B3", "子标题1")
	w.SetMergeCellValue("C3", "D3", "子标题2")

	w.SetTitles([]string{"姓名", "手机号"})

	w.SetData([]any{"李涛", "18510381580"}, defaultTitleStyle)

	w.Save2File("user.xlsx")
}
