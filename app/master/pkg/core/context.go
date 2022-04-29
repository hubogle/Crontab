package core

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type HandlerFunc func(c Context)

var _ Context = (*context)(nil)

type Context interface {
	init()
	JSON(code int, v interface{})
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
func (c *context) JSON(code int, obj interface{}) {
	c.ctx.JSON(code, obj)
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
