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

var (
	defaultCellStyle = &excelize.Style{
		Font: &excelize.Font{
			Size:   11,
			Family: "微软雅黑",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	}
	defaultTitleStyle = &excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
			Family: "微软雅黑",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#B8CCE4"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	}
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

func (x *xlsxWrite) SetStyle(left, right string, style *excelize.Style) error {
	styleId, _ := x.Handler().NewStyle(style)
	return x.Handler().SetCellStyle(x.sheetName, left, right, styleId)
}

func (x *xlsxWrite) SetMergeCellValue(left, right string, value any) error {
	if x.RowNum() == 0 {
		x.IncrRowNum()
	}

	if err := x.Handler().MergeCell(x.sheetName, left, right); err != nil {
		goo_log.Error(err)
		return err
	}
	if err := x.Handler().SetCellValue(x.sheetName, left, value); err != nil {
		goo_log.Error(err)
		return err
	}

	_ = x.SetStyle(left, right, defaultTitleStyle)

	return nil
}

func (x *xlsxWrite) SetTitles(titles []string) (left string, right string, err error) {
	x.IncrRowNum()

	left = fmt.Sprintf("A%d", x.RowNum())

	if err = x.Handler().SetSheetRow(x.sheetName, left, &titles); err != nil {
		goo_log.Error(err)
		return
	}

	if l := len(titles); l > 0 {
		columns := generateColumns(l)
		right = fmt.Sprintf("%s%d", columns[l-1], x.RowNum())
		_ = x.SetStyle(left, right, defaultTitleStyle)
	}

	return
}

func (x *xlsxWrite) SetData(data []interface{}) (left string, right string, err error) {
	x.IncrRowNum()

	left = fmt.Sprintf("A%d", x.RowNum())

	if err = x.Handler().SetSheetRow(x.sheetName, left, &data); err != nil {
		goo_log.Error(err)
		return
	}

	if l := len(data); l > 0 {
		columns := generateColumns(l)
		right = fmt.Sprintf("%s%d", columns[l-1], x.RowNum())
		_ = x.SetStyle(left, right, defaultCellStyle)
	}

	return
}

func (x *xlsxWrite) SetRows(data [][]interface{}) *xlsxWrite {
	for _, i := range data {
		x.IncrRowNum()

		left := fmt.Sprintf("A%d", x.RowNum())

		if err := x.Handler().SetSheetRow(x.sheetName, left, &i); err != nil {
			goo_log.Error(err)
			continue
		}

		if l := len(i); l > 0 {
			columns := generateColumns(l)
			right := fmt.Sprintf("%s%d", columns[l-1], x.RowNum())
			_ = x.SetStyle(left, right, defaultCellStyle)
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

// 生成列标识列表 (A 到 ZZ 示例)
func generateColumns(maxCols int) []string {
	columns := make([]string, maxCols)
	for i := 0; i < maxCols; i++ {
		columns[i] = excelColumnName(i + 1) // 列索引从 1 开始
	}
	return columns
}

// 将列索引(1-based)转换为 Excel 列名 (A, B, ..., Z, AA, AB, ...)
func excelColumnName(colIndex int) string {
	if colIndex < 1 {
		return ""
	}
	name := ""
	for colIndex > 0 {
		colIndex-- // 转换为 0-based 索引
		remainder := colIndex % 26
		name = string(rune('A'+remainder)) + name
		colIndex /= 26
	}
	return name
}
