package goowechat

import (
	"fmt"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"net/url"
	"strings"
	"time"
)

func JsApi(appid, secret, urlStr string) map[string]interface{} {
	ticket := CGITicket(appid, secret).Get()

	ts := time.Now().Unix()
	nonceStr := goo_utils.NonceStr()

	urlStr, _ = url.QueryUnescape(urlStr)
	urlStr = strings.Split(urlStr, "#")[0]

	rawstr := fmt.Sprintf(jsapi_ticket_qs, ticket, nonceStr, ts, urlStr)
	goolog.DebugF("Wechat-JsApi: %s", rawstr)
	rawstr = goo_utils.SHA1([]byte(rawstr))

	params := map[string]interface{}{
		"debug":     false,
		"appId":     appid,
		"timestamp": ts,
		"nonceStr":  nonceStr,
		"signature": rawstr,
		"jsApiList": []string{"checkJsApi", "updateAppMessageShareData", "updateTimelineShareData", "chooseWXPay",
			"openLocation", "getLocation", "chooseImage", "previewImage", "uploadImage", "downloadImage",
			"startRecord", "stopRecord", "onVoiceRecordEnd", "playVoice", "pauseVoice", "stopVoice",
			"onVoicePlayEnd", "uploadVoice", "downloadVoice", "translateVoice", "getNetworkType", "scanQRCode",
			"addCard", "chooseCard", "openCard"},
	}

	return params
}
