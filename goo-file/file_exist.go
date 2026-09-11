package goo_file

import (
	"os"
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
