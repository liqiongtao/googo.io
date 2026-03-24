package goo_xlsx

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/xuri/excelize/v2"
)

func Writer() *xlsxWrite {
	return &xlsxWrite{
		fh:        excelize.NewFile(),
		sheetName: "Sheet1",
		sheetRowNums: map[string]int{
			"Sheet1": 0,
		},
	}
}

type xlsxWrite struct {
	fh           *excelize.File
	sheetName    string
	sheetRowNums map[string]int
}

func (x *xlsxWrite) Handler() *excelize.File {
	return x.fh
}

func (x *xlsxWrite) SheetName() string {
	return x.sheetName
}

func (x *xlsxWrite) RowNum() int {
	return x.sheetRowNums[x.sheetName]
}

func (x *xlsxWrite) SetRowNum(num int) {
	x.sheetRowNums[x.sheetName] = num
}

func (x *xlsxWrite) IncrRowNum() {
	x.sheetRowNums[x.sheetName]++
}

func (x *xlsxWrite) SetStyle(start, end string, style *excelize.Style) error {
	left := fmt.Sprintf("%s%d", start, x.RowNum())
	right := fmt.Sprintf("%s%d", end, x.RowNum())

	styleId, _ := x.Handler().NewStyle(style)

	return x.Handler().SetCellStyle(x.sheetName, left, right, styleId)
}

func (x *xlsxWrite) SetStyleCenter(start, end string) error {
	style := &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center", // 水平居中
			Vertical:   "center", // 垂直居中
			WrapText:   true,     // 自动换行（可选）
		},
		//Font: &excelize.Font{
		//	Size:  16,       // 字体大小
		//	Color: "000000", // 字体颜色（十六进制，FF0000=红色）
		//	Bold:  true,     // 是否加粗（可选）
		//},
	}
	return x.SetStyle(start, end, style)
}

func (x *xlsxWrite) SetMergeCellValue(start, end string, value any) error {
	if x.RowNum() == 0 {
		x.IncrRowNum()
	}

	left := fmt.Sprintf("%s%d", start, x.RowNum())
	right := fmt.Sprintf("%s%d", end, x.RowNum())

	if err := x.Handler().MergeCell(x.sheetName, left, right); err != nil {
		goo_log.Error(err)
		return err
	}
	if err := x.Handler().SetCellValue(x.sheetName, left, value); err != nil {
		goo_log.Error(err)
		return err
	}

	x.SetStyleCenter(start, end)

	return nil
}

func (x *xlsxWrite) SetTitles(titles []string) error {
	x.IncrRowNum()
	if err := x.Handler().SetSheetRow(x.sheetName, fmt.Sprintf("A%d", x.RowNum()), &titles); err != nil {
		goo_log.Error(err)
		return err
	}
	return nil
}

func (x *xlsxWrite) SetData(data []interface{}) error {
	x.IncrRowNum()
	if err := x.Handler().SetSheetRow(x.sheetName, fmt.Sprintf("A%d", x.RowNum()), &data); err != nil {
		goo_log.Error(err)
		return err
	}
	return nil
}

func (x *xlsxWrite) SetRows(data [][]interface{}) *xlsxWrite {
	for _, i := range data {
		x.IncrRowNum()
		if err := x.Handler().SetSheetRow(x.sheetName, fmt.Sprintf("A%d", x.RowNum()), &i); err != nil {
			goo_log.Error(err)
			continue
		}
	}
	return x
}

func (x *xlsxWrite) SetSheetName(sheetName string) *xlsxWrite {
	_, _ = x.Handler().NewSheet(sheetName)
	x.sheetName = sheetName
	x.sheetRowNums[sheetName] = 0
	return x
}

func (x *xlsxWrite) Save2File(filename string) (err error) {
	if x.sheetRowNums["Sheet1"] == 0 {
		_ = x.Handler().DeleteSheet("Sheet1")
	}

	if err = x.Handler().SaveAs(filename); err != nil {
		goo_log.Error(err)
		return
	}
	defer func() { _ = x.Handler().Close() }()

	return nil
}

func (x *xlsxWrite) Output(ctx *gin.Context, filename string) (err error) {
	tmpFile := fmt.Sprintf("/tmp/%s", filename)
	defer func() { _ = os.Remove(tmpFile) }()

	if err = x.Save2File(tmpFile); err != nil {
		return err
	}

	fileInfo, _ := os.Stat(tmpFile)
	fileSize := fileInfo.Size()

	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+url.PathEscape(filename))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Header("content-length", strconv.FormatInt(fileSize, 10))

	ctx.File(tmpFile)

	return
}
