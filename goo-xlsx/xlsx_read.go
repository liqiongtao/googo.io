package gooxlsx

import (
	"io"
	"os"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	"github.com/xuri/excelize/v2"
)

func ReadBySheet(r io.Reader, sheet string, fn func(n int, row []string) error) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		goolog.Error(err)
		return err
	}
	defer xlsx.Close()

	if sheet == "" {
		sheet = xlsx.GetSheetName(0)
	}

	rows, err := xlsx.Rows(sheet)
	if err != nil {
		goolog.Error(err)
		return err
	}
	defer rows.Close()

	var n int
	for rows.Next() {
		n++
		row, err := rows.Columns()
		if err != nil {
			goolog.Error(err)
			return err
		}
		if err = fn(n, row); err != nil {
			return err
		}
	}
	if err := rows.Error(); err != nil {
		goolog.Error(err)
		return err
	}

	return nil
}

func Read(r io.Reader, fn func(n int, row []string) error) error {
	return ReadBySheet(r, "", fn)
}

func ReadFile(file string, fn func(n int, row []string) error) error {
	h, err := os.Open(file)
	if err != nil {
		goolog.Error(err)
		return err
	}
	defer h.Close()
	return Read(h, fn)
}
