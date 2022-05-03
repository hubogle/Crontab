package core

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"sync"
)

type HandlerFunc func(c Context)

var _ Context = (*context)(nil)

type Context interface {
	init()
	JSON(code int, v interface{})
	JSONError(err BusinessError) // 返回错误信息
	// ShouldBindJSON 反序列化，tag: `json:"xxx"`
	ShouldBindJSON(obj interface{}) error
	ShouldBindURI(obj interface{}) error // 需要传递指针
	ShouldBind(obj interface{}) error
	ShouldBindQuery(obj interface{}) error
}

// context 进行封装
type context struct {
	ctx *gin.Context
}

func (c *context) init() {
}

// ResponseWriter 获取 ResponseWriter
func (c *context) ResponseWriter() gin.ResponseWriter {
	return c.ctx.Writer
}

// ShouldBindQuery 反序列化querystring tag: `form:"xxx"`
func (c *context) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Query)
}

// ShouldBindURI 反序列化path参数(如路由路径为 /user/:name)	tag: `uri:"xxx"`
func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

func (c *context) ShouldBind(obj interface{}) error {
	return c.ctx.ShouldBind(obj)
}

func (c *context) JSON(code int, obj interface{}) {
	c.ctx.JSON(code, obj)
}
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}
func (c *context) JSONError(err BusinessError) {
	if err != nil {
		httpCode := err.HTTPCode()
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}
		c.ctx.JSON(httpCode, map[string]interface{}{
			"code": err.BusinessCode(),
			"msg":  err.Message(),
		})
	}
}

var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

// GetContext 通过从 pool 中获取 context
func newContext(ctx *gin.Context) Context {
	context := contextPool.Get().(*context)
	context.ctx = ctx
	return context
}

// ReleaseContext 释放 context
func releaseContext(ctx Context) {
	c := ctx.(*context)
	c.ctx = nil
	contextPool.Put(c)
}
