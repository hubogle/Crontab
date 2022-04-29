package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RouterGroup 包装gin的RouterGroup，调用封装的 Router 方法
type RouterGroup interface {
	Group(string, ...HandlerFunc) RouterGroup
	IRoutes
}
type Mux interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}
type mux struct {
	engine *gin.Engine
}

// ServeHTTP 实现http.Handler接口
func (m *mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.engine.ServeHTTP(w, req)
}

// Group 封装gin的Group方法，单独对 handlers 进行了处理
func (m *mux) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &router{
		group: m.engine.Group(relativePath, wrapHandlers(handlers...)...),
	}
}

// 单独对 gin.HandlerFunc 进行处理，主要就是对 gin.Context 进行包装，通过 Pool 实现资源池
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcObj := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
		funcObj[i] = func(c *gin.Context) {
			ctx := newContext(c) // 从 Pool 获取 ctx
			defer releaseContext(ctx)
			handler(ctx)
		}
	}
	return funcObj
}

func New() (Mux, error) {
	mux := &mux{engine: gin.Default()}
	return mux, nil
}
