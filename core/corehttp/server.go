package corehttp

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/hamster-provider/core/context"
)

type MyServer struct {
	*gin.Engine
	CoreContext *context.CoreContext
}

//MyContext extend gin.Context
type MyContext struct {
	*gin.Context
	CoreContext *context.CoreContext
}

//MyRouterGroup extend gin.RouterGroup
type MyRouterGroup struct {
	*gin.RouterGroup
	CoreContext *context.CoreContext
}

type HandlerFunc func(c *MyContext)

//NewMyServer new gin server
func NewMyServer(ctx *context.CoreContext) *MyServer {
	myServer := &MyServer{
		Engine:      gin.Default(),
		CoreContext: ctx,
	}
	return myServer
}

//handleFunc  implement gin.Context to MyContext。
func handleFunc(handler func(c *MyContext), ctx *context.CoreContext) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		handler(&MyContext{Context: c, CoreContext: ctx})
	}
}

//Group rewrite route register
func (server *MyServer) Group(relativePath string, handlers ...HandlerFunc) *MyRouterGroup {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle, server.CoreContext))
	}
	return &MyRouterGroup{server.Engine.Group(relativePath, RHandles...), server.CoreContext}
}

//GET extent GET
func (server *MyServer) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle, server.CoreContext))
	}
	return server.Engine.GET(relativePath, RHandles...)
}

//POST extend POST
func (server *MyServer) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle, server.CoreContext))
	}
	return server.Engine.POST(relativePath, RHandles...)
}

func (group *MyRouterGroup) Group(relativePath string, handlers ...HandlerFunc) *MyRouterGroup {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle, group.CoreContext))
	}
	return &MyRouterGroup{group.RouterGroup.Group(relativePath, RHandles...), group.CoreContext}
}

//GET extend group GET
func (r *MyRouterGroup) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle, r.CoreContext))
	}
	return r.RouterGroup.GET(relativePath, rHandles...)
}

//POST extend group POST
func (r *MyRouterGroup) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle, r.CoreContext))
	}
	return r.RouterGroup.POST(relativePath, rHandles...)
}

//Use extend register with middle
func (r *MyRouterGroup) Use(middlewares ...func(c *MyContext)) gin.IRoutes {
	rMiddlewares := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		rMiddlewares = append(rMiddlewares, handleFunc(middleware, r.CoreContext))
	}
	return r.RouterGroup.Use(rMiddlewares...)
}
