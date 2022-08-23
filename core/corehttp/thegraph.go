package corehttp

import (
	"github.com/gorilla/websocket"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"github.com/hamster-shared/hamster-provider/log"
	"net/http"
	"os"
	"strings"
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

func SS58AuthMiddleware(c *MyContext) {
	ss58AuthData := c.Request.Header.Get("SS58AuthData")
	if ss58AuthData == "" {
		c.JSON(http.StatusBadRequest, BadRequest("not found SS58AuthData"))
		c.Abort()
		return
	}
	ss58Address, data, signHex, ok := parseSS58AuthData(ss58AuthData)
	if !ok {
		c.JSON(http.StatusBadRequest, BadRequest("parse SS58AuthData error"))
		c.Abort()
		return
	}
	if ss58Address == "" || data == "" || signHex == "" {
		c.JSON(http.StatusBadRequest, BadRequest("SS58AuthData have empty value"))
		c.Abort()
		return
	}
	if verifyResult := thegraph.VerifyWithSS58AndHexSign(ss58Address, data, signHex); !verifyResult {
		c.JSON(http.StatusBadRequest, BadRequest("SS58AuthData verify failed"))
		c.Abort()
		return
	}
	c.Next()
}

func parseSS58AuthData(str string) (ss58Address, data, signHex string, ok bool) {
	if str == "" {
		return
	}
	split := strings.Split(str, ":")
	if len(split) != 3 {
		return
	}
	return split[0], split[1], split[2], true
}

func graphStart(c *MyContext) {
	deploymentID := c.QueryArray("deploymentID")
	if len(deploymentID) == 0 {
		c.JSON(http.StatusBadRequest, BadRequest("not found deploymentID"))
		return
	}
	err := thegraph.DefaultComposeGraphConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, BadRequest(err.Error()))
		return
	}
	err = thegraph.DefaultComposeGraphStart(deploymentID...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, BadRequest(err.Error()))
	} else {
		c.JSON(http.StatusOK, Success("ok"))
	}
}

func graphStop(c *MyContext) {
	deploymentID := c.QueryArray("deploymentID")
	if len(deploymentID) == 0 {
		c.JSON(http.StatusBadRequest, BadRequest("not found deploymentID"))
		return
	}
	err := thegraph.DefaultComposeGraphConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, BadRequest(err.Error()))
		return
	}
	err = thegraph.DefaultComposeGraphStop(deploymentID...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, BadRequest(err.Error()))
	} else {
		c.JSON(http.StatusOK, Success("ok"))
	}
}

func graphRules(c *MyContext) {
	result, err := thegraph.DefaultComposeGraphRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, BadRequest(err.Error()))
	} else {
		c.JSON(http.StatusOK, Success(result))
	}
}
