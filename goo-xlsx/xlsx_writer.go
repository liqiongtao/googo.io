package gooxlsx

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	goolog "github.com/liqiongtao/googo.io/goo-log"
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

func (x *xlsxWrite) SetMergeCellValue(left, right string, value any, styles ...*excelize.Style) error {
	if x.RowNum() == 0 {
		x.IncrRowNum()
	}

	if err := x.Handler().MergeCell(x.sheetName, left, right); err != nil {
		goolog.Error(err)
		return err
	}
	if err := x.Handler().SetCellValue(x.sheetName, left, value); err != nil {
		goolog.Error(err)
		return err
	}

	if len(styles) == 0 {
		styles = append(styles, defaultTitleStyle)
	}

	_ = x.SetStyle(left, right, styles[0])

	return nil
}

func (x *xlsxWrite) SetTitles(titles []string, styles ...*excelize.Style) error {
	x.IncrRowNum()

	left := fmt.Sprintf("A%d", x.RowNum())

	if err := x.Handler().SetSheetRow(x.sheetName, left, &titles); err != nil {
		goolog.Error(err)
		return err
	}

	if l := len(titles); l > 0 {
		columns := generateColumns(l)
		right := fmt.Sprintf("%s%d", columns[l-1], x.RowNum())

		if len(styles) == 0 {
			styles = append(styles, defaultTitleStyle)
		}

		_ = x.SetStyle(left, right, styles[0])
	}

	return nil
}

func (x *xlsxWrite) SetData(data []any, styles ...*excelize.Style) error {
	x.IncrRowNum()

	left := fmt.Sprintf("A%d", x.RowNum())

	if err := x.Handler().SetSheetRow(x.sheetName, left, &data); err != nil {
		goolog.Error(err)
		return err
	}

	if l := len(data); l > 0 {
		columns := generateColumns(l)
		right := fmt.Sprintf("%s%d", columns[l-1], x.RowNum())

		if len(styles) == 0 {
			styles = append(styles, defaultCellStyle)
		}

		_ = x.SetStyle(left, right, styles[0])
	}

	return nil
}

func (x *xlsxWrite) SetRows(data [][]any, styles ...*excelize.Style) *xlsxWrite {
	for _, i := range data {
		x.IncrRowNum()

		left := fmt.Sprintf("A%d", x.RowNum())

		if err := x.Handler().SetSheetRow(x.sheetName, left, &i); err != nil {
			goolog.Error(err)
			continue
		}

		if l := len(i); l > 0 {
			columns := generateColumns(l)
			right := fmt.Sprintf("%s%d", columns[l-1], x.RowNum())

			if len(styles) == 0 {
				styles = append(styles, defaultCellStyle)
			}

			_ = x.SetStyle(left, right, styles[0])
		}
	}

	return x
}

func (x *xlsxWrite) SetSheetName(sheetName string) *xlsxWrite {
	_, _ = x.Handler().NewSheet(sheetName)
	x.sheetName = sheetName
	if _, ok := x.sheetRowNums[sheetName]; !ok {
		x.sheetRowNums[sheetName] = 0
	}
	return x
}

func (x *xlsxWrite) Save2File(filename string) (err error) {
	if x.sheetRowNums["Sheet1"] == 0 {
		_ = x.Handler().DeleteSheet("Sheet1")
	}

	defer func() { _ = x.Handler().Close() }()

	if err = x.Handler().SaveAs(filename); err != nil {
		goolog.Error(err)
		return
	}

	return nil
}

func (x *xlsxWrite) Output(ctx *gin.Context, filename string) (err error) {
	filename = filepath.Base(filename)
	if filename == "" || filename == "." || filename == string(filepath.Separator) {
		return fmt.Errorf("invalid filename")
	}

	tmpFile, err := os.CreateTemp("", "goo-xlsx-*.xlsx")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	if err = x.Save2File(tmpPath); err != nil {
		return err
	}

	fileInfo, err := os.Stat(tmpPath)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+url.PathEscape(filename))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Header("content-length", strconv.FormatInt(fileSize, 10))

	ctx.File(tmpPath)

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
