/**
 * 收集爬虫代理服务器地址
 */

package main 

import (
 	"mjz_spider/handlers/proxy"
)

func main() {
	for _, spider := range proxy.Handlers {
		sp := spider()
		sp.Run()
	}
}