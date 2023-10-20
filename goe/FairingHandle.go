package goe

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type FairingHandler struct {
	fairings []Fairing
}

func NewFairingHandler() *FairingHandler {
	return &FairingHandler{}
}

var fairingHandler *FairingHandler
var fairingHandlerOnce sync.Once

func getFairingHandler() *FairingHandler {
	fairingHandlerOnce.Do(func() {
		fairingHandler = &FairingHandler{}
	})
	return fairingHandler
}

func (this *FairingHandler) AddFairing(f ...Fairing) {
	if f != nil && len(f) > 0 {
		this.fairings = append(this.fairings, f...)
	}
}

// 中间件onRequest处理
func (this *FairingHandler) before(ctx *gin.Context) {
	for _, f := range this.fairings {
		err := f.OnRequest(ctx)
		if err != nil {
			Throw(err.Error(), 400, ctx)
		}
	}
}

// 中间件onResponse处理
func (this *FairingHandler) after(ctx *gin.Context, ret interface{}) interface{} {
	var result = ret
	for _, f := range this.fairings {
		res, err := f.OnResponse(result)
		if err != nil {
			Throw(err.Error(), 400, ctx)
		}
		result = res
	}
	return result
}

func (this *FairingHandler) handlerFairing(responder Responder, ctx *gin.Context) interface{} {
	this.before(ctx)
	var ret interface{}
	innerNode := getInnerRouter().getRoute(ctx.Request.Method, ctx.Request.URL.Path)
	var innerFairingHandler *FairingHandler
	if innerNode.fullPath != "" && innerNode.handlers != nil {
		if fs, ok := innerNode.handlers.([]Fairing); ok {
			innerFairingHandler = NewFairingHandler()
			innerFairingHandler.AddFairing(fs...)
		}
	}
	if innerFairingHandler != nil {
		innerFairingHandler.before(ctx)
	}
	if s1, ok := responder.(StringResponder); ok {
		ret = s1(ctx)
	}
	if s2, ok := responder.(JsonResponder); ok {
		ret = s2(ctx)
	}
	if s3, ok := responder.(SqlResponder); ok {
		ret = s3(ctx)
	}
	if s4, ok := responder.(SqlQueryResponder); ok {
		ret = s4(ctx)
	}
	if s5, ok := responder.(VoidResponder); ok {
		s5(ctx)
		ret = struct{}{}
	}
	// exec route-level middleware
	if innerFairingHandler != nil {
		ret = innerFairingHandler.after(ctx, ret)
	}

	return getFairingHandler().after(ctx, ret)
}
