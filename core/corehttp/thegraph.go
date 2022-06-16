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

func deployTheGraph(c *MyContext) {

	var data thegraph.DeployParams
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = thegraph.Deploy(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, Success(""))
}

func execHandler(c *MyContext) {

	containerName := c.Param("serviceName")

	// websocket握手
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	defer conn.Close()

	thegraph.GetWebSocket(conn, containerName)
}

func logHandler(c *MyContext) {

	containerName := c.Param("serviceName")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	defer conn.Close()

	thegraph.GetDockerLog(conn, containerName)
}
