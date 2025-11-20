package goo_pprof

import (
	"testing"

	goo_context "github.com/liqiongtao/googo.io/goo-context"
)

func TestRun(t *testing.T) {
	Run()

	<-goo_context.WithCancel().Done()
}
