package corehttp

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"net/http"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func TerminalHandle(c *MyContext) {

	// websocket握手
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	defer conn.Close()

	thegraph.GetWebSocket(conn)
}

func TerminalLogHandle(c *MyContext) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	defer conn.Close()

	thegraph.GetDockerLog(conn)
}
