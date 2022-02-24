package main

import (
	"github.com/liqiongtao/googo.io/goo"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"log"
)

func main() {
	//goo_log.SetAdapter(goo_log.NewFileAdapter())

	s := goo.NewServer(goo.BaseDirOption(goo_utils.DIR()))

	s.Use(func(ctx *goo.Context) {
		//ctx.AbortWithStatusJSON(401, goo.Error(5002, "bbb"))
	}, func(ctx *goo.Context) {
		//ctx.AbortWithStatusJSON(401, goo.Error(5001, "aaa"))
	})

	s.POST("/a", goo.Handler(A{}))

	s.Run(":9003")
}

type A struct {
	Request struct {
		ContractNo    string `json:"contract_no"`                        // 合同编号
		Amount        int64  `json:"amount" binding:"required"`          // 预收款金额，单位分
		InvoiceCode   string `json:"invoice_code" binding:"required"`    // 发票代码
		InvoiceNo     string `json:"invoice_no" binding:"required"`      // 发票号码
		PaymentTime   string `json:"payment_time" binding:"required"`    // 付款时间
		PaymentPartyA string `json:"payment_party_a" binding:"required"` // 付款甲方
	}
}

func (this A) DoHandle(ctx *goo.Context) *goo.Response {
	defer func() {
		if err := recover(); err != nil {
			log.Println("err:", err)
		}
	}()

	if err := ctx.ShouldBind(&this.Request); err != nil {
		return goo.Error(7001, goo.ValidationMessage(err, map[string]string{
			"amount_required":          "预收款金额 为空",
			"invoice_code_required":    "发票代码 为空",
			"invoice_no_required":      "发票号码 为空",
			"payment_time_required":    "付款时间 为空",
			"payment_party_a_required": "付款甲方 为空",
		}))
	}

	return goo.Success(nil)
}

func b() {
	//c()
	panic("aa")
}

func c() {
	panic("aa")
}
