package goo_xlsx

import (
	"fmt"
	"github.com/liqiongtao/googo.io/goo"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/xuri/excelize/v2"
)

func Writer() *xlsxWrite {
	return &xlsxWrite{
		fh:        excelize.NewFile(),
		sheetName: "Sheet1",
	}
}

type xlsxWrite struct {
	fh        *excelize.File
	sheetName string
	titles    *[]string
	rows      []*[]interface{}
}

func (x *xlsxWrite) SetTitles(titles []string) *xlsxWrite {
	x.titles = &titles
	return x
}

func (x *xlsxWrite) SetData(data []interface{}) *xlsxWrite {
	x.rows = append(x.rows, &data)
	return x
}

func (x *xlsxWrite) SetRows(data [][]interface{}) *xlsxWrite {
	for _, i := range data {
		x.rows = append(x.rows, &i)
	}
	return x
}

func (x *xlsxWrite) SetSheetName(sheetName string) *xlsxWrite {
	x.sheetName = sheetName
	return x
}

func (x *xlsxWrite) Save2File(filename string) (err error) {
	if err = x.fh.SetSheetRow(x.sheetName, "A1", x.titles); err != nil {
		goo_log.Error(err)
		return
	}

	for i := 0; i < len(x.rows); i++ {
		if err = x.fh.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), x.rows[i]); err != nil {
			goo_log.Error(err)
			return
		}
	}

	if err = x.fh.SaveAs(filename); err != nil {
		goo_log.Error(err)
		return
	}

	return nil
}

func (x *xlsxWrite) Output(ctx *goo.Context, filename string) (err error) {
	if err = x.fh.SetSheetRow(x.sheetName, "A1", x.titles); err != nil {
		goo_log.Error(err)
		return
	}

	for i := 0; i < len(x.rows); i++ {
		if err = x.fh.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), x.rows[i]); err != nil {
			goo_log.Error(err)
			return
		}
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err = x.fh.Write(ctx.Writer); err != nil {
		goo_log.Error(err)
		return
	}

	return
}
