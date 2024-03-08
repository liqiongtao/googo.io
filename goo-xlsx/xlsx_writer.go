package goo_xlsx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/xuri/excelize/v2"
	"net/url"
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

func (x *xlsxWrite) SetTitles(titles []string) error {
	x.sheetRowNums[x.sheetName]++
	if err := x.fh.SetSheetRow(x.sheetName, fmt.Sprintf("A%d", x.sheetRowNums[x.sheetName]), &titles); err != nil {
		goo_log.Error(err)
		return err
	}
	return nil
}

func (x *xlsxWrite) SetData(data []interface{}) error {
	x.sheetRowNums[x.sheetName]++
	if err := x.fh.SetSheetRow(x.sheetName, fmt.Sprintf("A%d", x.sheetRowNums[x.sheetName]), &data); err != nil {
		goo_log.Error(err)
		return err
	}
	return nil
}

func (x *xlsxWrite) SetRows(data [][]interface{}) *xlsxWrite {
	for _, i := range data {
		x.sheetRowNums[x.sheetName]++
		if err := x.fh.SetSheetRow(x.sheetName, fmt.Sprintf("A%d", x.sheetRowNums[x.sheetName]), &i); err != nil {
			goo_log.Error(err)
			continue
		}
	}
	return x
}

func (x *xlsxWrite) SetSheetName(sheetName string) *xlsxWrite {
	x.fh.NewSheet(sheetName)
	x.sheetName = sheetName
	x.sheetRowNums[sheetName] = 0
	return x
}

func (x *xlsxWrite) Save2File(filename string) (err error) {
	if err = x.fh.SaveAs(filename); err != nil {
		goo_log.Error(err)
		return
	}

	return nil
}

func (x *xlsxWrite) Output(ctx *gin.Context, filename string) (err error) {
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+url.PathEscape(filename))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")

	if err = x.fh.Write(ctx.Writer); err != nil {
		goo_log.Error(err)
		return
	}

	return
}
