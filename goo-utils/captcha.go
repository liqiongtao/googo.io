package goo_utils

import (
	"github.com/mojocn/base64Captcha"
	"sync"
)

var (
	__store     base64Captcha.Store
	__storeOnce sync.Once
)

// 获取图片验证码
func CaptchaGet(width, height int) map[string]string {
	// 初始化
	__storeOnce.Do(func() {
		__store = base64Captcha.DefaultMemStore
	})

	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	driver := base64Captcha.DriverDigit{
		Height:   height,
		Width:    width,
		Length:   4,
		MaxSkew:  0.7,
		DotCount: 80,
	}
	ca := base64Captcha.NewCaptcha(&driver, __store)
	id, b64s, _ := ca.Generate()

	return map[string]string{
		"id":          id,
		"base64image": b64s,
	}
}

// 验证图片验证码
func CaptchaVerify(id, code string) bool {
	if __store.Verify == nil {
		return false
	}
	return __store.Verify(id, code, true)
}
