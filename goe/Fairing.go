package goe

import "github.com/gin-gonic/gin"

// Fairing 规范中间件代码和功能的接口
type Fairing interface {
	OnRequest(context *gin.Context) error
	OnResponse(ret interface{}) (interface{}, error)
}
