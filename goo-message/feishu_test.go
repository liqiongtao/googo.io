package goo_message

import (
	"fmt"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	"testing"
)

func TestFeiShu(t *testing.T) {
	hookUrl := "https://open.feishu.cn/open-apis/bot/v2/hook/ce66f466-aac3-44a2-97dd-949da4e24853"

	for i := 0; i < 10; i++ {
		//if i%10 != 0 {
		go FeiShu(hookUrl, fmt.Sprintf("测试"))
		//	continue
		//}
		//go FeiShu(hookUrl, fmt.Sprintf("测试%d", i))
	}

	<-goo_context.WithCancel().Done()
}
