package corehttp

import (
	"github.com/gorilla/websocket"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"github.com/hamster-shared/hamster-provider/log"
	"net/http"
	"os"
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
		log.GetLogger().Info("binding json error")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.GetLogger().Info("received deploy command:", data)

	//data := thegraph.DeployParams{}
	err = thegraph.Deploy(data)
	if err != nil {
		log.GetLogger().Info("deploy error , error is : ", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, Success(""))
}

func pullImage(c *MyContext) {
	err := thegraph.PullImage()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, Success(""))
}

func execHandler(c *MyContext) {

	_ = os.Setenv("DOCKER_API_VERSION", "1.41")
	containerName := c.Query("serviceName")

	// websocket握手
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.GetLogger().Info("error", err.Error())
		return
	}
	defer conn.Close()

	thegraph.GetWebSocket(conn, containerName)
}

func logHandler(c *MyContext) {

	containerName := c.Query("serviceName")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.GetLogger().Info("error", err.Error())
		return
	}

	defer conn.Close()

	thegraph.GetDockerLog(conn, containerName)
}

func deployStatus(c *MyContext) {
	containerNames := c.QueryArray("serviceName")

	if len(containerNames) > 0 {
		status, err := thegraph.GetDockerComposeStatus(containerNames...)
		if err != nil {
			c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Success(status))
	} else {
		c.JSON(http.StatusBadRequest, BadRequest("not found serviceName"))
	}
}
