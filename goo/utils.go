package goo

import goo_utils "github.com/liqiongtao/googo.io/goo-utils"

// LocalIP 获取本地IP地址
func LocalIP() (string, error) {
	return goo_utils.LocalIP()
}
