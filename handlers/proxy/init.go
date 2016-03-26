package proxy

type spiderHandler func() SpiderHandler

type SpiderHandler interface {
	Run()
} 

var Handlers []spiderHandler

/**
 * 注册spider handler
 */
func Register(hanlder spiderHandler) {
	Handlers = append(Handlers, hanlder)
}