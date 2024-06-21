package goo_xlsx

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
)

func ReadBySheet(r io.Reader, sheet string, fn func(n int, row []string) error) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		goo_log.Error(err)
		return err
	}
	defer xlsx.Close()

	if sheet == "" {
		sheet = xlsx.GetSheetName(0)
	}

	rows, err := xlsx.Rows(sheet)
	if err != nil {
		goo_log.Error(err)
		return err
	}
	defer rows.Close()

	var n int
	for rows.Next() {
		n++
		row, _ := rows.Columns()
		if err = fn(n, row); err != nil {
			return err
		}
	}

	return nil
}

func Read(r io.Reader, fn func(n int, row []string) error) error {
	return ReadBySheet(r, "", fn)
}

func ReadFile(file string, fn func(n int, row []string) error) error {
	h, err := os.Open(file)
	if err != nil {
		goo_log.Error(err)
		return err
	}
	return Read(h, fn)
}
