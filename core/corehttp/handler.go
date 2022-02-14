package corehttp

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type Keys struct {
	PublicKey string `json:"public_key"`
}

// @Summary startContainer
// @Description startContainer
// @Tags containerOperation
// @Accept json
// @Produce json
// @Success 200 {string} startContainer
// @Router /container/start [get]
func startContainer(c *MyContext) {

	name := c.Query("name")
	if err := c.CoreContext.VmManager.Start(name); err != nil {
		c.String(http.StatusBadRequest, "err:%v", err)
		return
	}
	c.String(http.StatusOK, "start container success")
}

// @Summary delete container
// @Description delete container
// @Tags containerOperation
// @Accept json
// @Produce json
// @Success 200 {string} deleteContainer
// @Router /container/delete [get]
func deleteContainer(c *MyContext) {

	name := c.Query("name")
	if err := c.CoreContext.VmManager.Destroy(name); err != nil {
		c.String(http.StatusBadRequest, "err:%v", err)
		return
	}
	c.String(http.StatusOK, "delete container success")
}

// @Summary grant key
// @Description grant key
// @Tags key operation
// @Accept json
// @Produce json
// @Param param body Keys true "the key need upload"
// @Success 200 {string} grantKey
// @Router /pk/grantKey [POST]
func grantKey(gin *MyContext) {

	name := gin.Query("name")
	json := Keys{}
	if err := gin.BindJSON(&json); err != nil {
		gin.String(http.StatusBadRequest, "err:%v", err)
		return
	}

	if err := gin.CoreContext.PkManager.AddPublicKey(json.PublicKey); err != nil {
		gin.String(http.StatusBadRequest, "err:%v", err)
		return
	}

	if err := gin.CoreContext.VmManager.InjectionPublicKey(name, json.PublicKey); err != nil {
		gin.String(http.StatusBadRequest, "err:%v", err)
		return
	}

	gin.String(http.StatusOK, "you public key is :%s\n", json.PublicKey)
	gin.String(http.StatusOK, "public key upload success \n")
}

// @Summary delete key
// @Description delete key
// @Tags key operation
// @Accept json
// @Produce json
// @Param param body Keys true "the key want delete"
// @Success 200 {string} grantKey
// @Router /pk/deleteKey [POST]
func deleteKey(gin *MyContext) {

	json := Keys{}
	if err := gin.BindJSON(&json); err != nil {
		gin.String(http.StatusBadRequest, "err:%v", err)
		return
	}

	if err := gin.CoreContext.PkManager.DeletePublicKey(json.PublicKey); err != nil {
		gin.String(http.StatusBadRequest, "err:%v", err)
		return
	}
	gin.String(http.StatusOK, "the key you want delete is :%s\n", json.PublicKey)
	gin.String(http.StatusOK, "the key delete success\n")
}

// @Summary query key
// @Description uery key
// @Tags key operation
// @Accept json
// @Produce json
// @Param param body Keys true "the key you want query"
// @Success 200 {string} queryKey
// @Router /pk/queryKey [POST]
func queryKey(gin *MyContext) {
	json := Keys{}
	if err := gin.BindJSON(&json); err != nil {
		gin.String(http.StatusBadRequest, "err:%v", err)
		return
	}

	res, err := gin.CoreContext.PkManager.QueryPublicKey(json.PublicKey)
	if err != nil {
		gin.String(http.StatusBadRequest, "err:%v", err)
		return
	}

	gin.String(http.StatusOK, "the key you query is :%s\n", json.PublicKey)

	if res == true {
		gin.String(http.StatusOK, "the key you query is  exists")
	} else {
		gin.String(http.StatusOK, "the key you query is not exists")
	}
}

// @Summary p2p port listen
// @Description p2p port listen
// @Tags p2p
// @Accept json
// @Produce json
// @Param param from port true ""
// @Success 200 {string} listenP2p
// @Router /p2p/listen [POST]
func listenP2p(gin *MyContext) {

	portStr := gin.Query("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logrus.Error("port format is invalid")
		gin.String(400, "port is not integer")
		return
	}

	if port > 65535 || port <= 0 {
		logrus.Error("port range invalid ，getter 65535 or smaller 0")
		gin.String(400, "port range invalid ，getter 65535 or smaller 0")
		return
	}

	target := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)
	err = gin.CoreContext.P2pClient.Listen(target)
	if err != nil {
		logrus.Error("p2p port create fail")
		gin.String(400, "p2p port create fail")
		return
	}

	gin.String(http.StatusOK, "p2p port create success")
}

// @Summary p2p port forward to local
// @Description p2p port forward to local
// @Tags p2p
// @Accept json
// @Produce json
// @Param param from port true ""
// @Success 200 {string} forwardP2p
// @Router /p2p/forward [POST]
func forwardP2p(gin *MyContext) {

	portStr := gin.Query("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logrus.Error("port format is invalid")
		gin.String(400, "port is not integer")
		return
	}

	if port > 65535 || port <= 0 {
		logrus.Error("port range invalid ，getter 65535 or smaller 0")
		gin.String(400, "port range invalid ，getter 65535 or smaller 0")
		return
	}

	targetPeerId := gin.Query("peerId")

	err = gin.CoreContext.P2pClient.Forward(port, targetPeerId)
	if err != nil {
		logrus.Error("p2p port create fail")
		gin.String(400, "p2p port create fail")
		return
	}

	gin.String(http.StatusOK, "p2p port create success")
}

// @Summary show p2p listen list
// @Description show p2p listen list
// @Tags p2p
// @Accept json
// @Produce json
// @Param param from port true ""
// @Success 200 {string} forwardP2p
// @Router /p2p/forward [POST]
func lsP2p(gin *MyContext) {
	gin.JSON(http.StatusOK, gin.CoreContext.P2pClient.List())
}

// @Summary close matcher p2p connection
// @Description close matcher p2p connection
// @Tags p2p
// @Accept json
// @Produce json
// @Param param from port true ""
// @Success 200 {string} closeP2p
// @Router /p2p/close [POST]
func closeP2p(gin *MyContext) {
	target := gin.Query("target")

	done, err := gin.CoreContext.P2pClient.Close(target)

	if err != nil {
		gin.String(http.StatusBadRequest, "", err.Error())
		return
	}

	gin.String(http.StatusOK, "%d Stream closed", done)
}

// @Summary check p2p connection can be connected
// @Description check p2p connection can be connected
// @Tags p2p
// @Accept json
// @Produce json
// @Param param from port true ""
// @Success 200 {string} closeP2p
// @Router /p2p/check [POST]
func checkP2p(gin *MyContext) {
	targetPeerId := gin.Query("peerId")

	err := gin.CoreContext.P2pClient.CheckForwardHealth(targetPeerId)
	if err != nil {
		gin.String(http.StatusBadRequest, "p2p connection is not ready")
	} else {
		gin.String(http.StatusOK, "p2p connection is ready")
	}
}

// createVm create virtual machine
func createVm(gin *MyContext) {
	name := gin.Query("name")
	manage := gin.CoreContext.VmManager
	logrus.Info("create virtual machine  start")
	err := manage.Create(name)
	logrus.Info("create virtual machine  end")
	if err != nil {
		logrus.Error("create virtual machine  fail")
		fmt.Println(err)
	}
}

func modifyPrice(gin *MyContext) {
	reportClient := gin.CoreContext.ReportClient
	price, err := strconv.Atoi(gin.Query("price"))
	if err != nil {
		gin.String(400, "Incorrect parameter format : %s", gin.Query("price"))
		return
	}
	ri := gin.CoreContext.GetConfig().ChainRegInfo.ResourceIndex
	err = reportClient.ModifyResourcePrice(ri, int64(price))
	if err != nil {
		gin.String(400, "modify price fail: %d", price)
	} else {
		gin.String(200, "modify price success")
	}
}

func addDuration(gin *MyContext) {
	durationParam := gin.Query("duration")

	duration, err := time.ParseDuration(durationParam)
	if err != nil {
		gin.String(400, "duration format is not invalid ")
		return
	}
	ri := gin.CoreContext.GetConfig().ChainRegInfo.ResourceIndex
	err = gin.CoreContext.ReportClient.AddResourceDuration(ri, int(duration.Hours()))
	if err != nil {
		gin.String(400, "add duration fail: %d", duration.Hours())
	} else {
		gin.String(200, "add duration success")
	}
}

func getConfig(gin *MyContext) {
	cfg := gin.CoreContext.GetConfig()
	gin.JSON(http.StatusOK, Success(cfg))
}

func setConfig(gin *MyContext) {
	cfg := gin.CoreContext.GetConfig()
	reqBody := config.Config{}
	if err := gin.BindJSON(&reqBody); err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest())
		return
	}

	cfg.Vm = reqBody.Vm
	cfg.ChainApi = reqBody.ChainApi
	// 校验seed 是否合法
	_, err := signature.KeyringPairFromSecret(reqBody.SeedOrPhrase, 42)
	if err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest("seed not invalid"))
		return
	}

	cfg.SeedOrPhrase = reqBody.SeedOrPhrase
	cfg.ConfigFlag = config.DONE

	err = gin.CoreContext.Cm.Save(cfg)
	if err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest("save config fail"))
		return
	}

	gin.JSON(http.StatusOK, Success(""))
}

func setBootState(gin *MyContext) {

	var op BootParam

	if err := gin.BindJSON(&op); err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest())
		return
	}

	gin.CoreContext.ChainListener.SetState(op.Option)

	gin.JSON(http.StatusOK, Success(""))
}

func getBootState(gin *MyContext) {
	gin.JSON(http.StatusOK, Success(gin.CoreContext.ChainListener.GetState()))
}

func getChainResource(gin *MyContext) {
	resourceIndex := gin.CoreContext.GetConfig().ChainRegInfo.ResourceIndex
	info, err := gin.CoreContext.ReportClient.GetResource(resourceIndex)
	if err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest("query resource fail"))
	} else {
		gin.JSON(http.StatusOK, Success(info))
	}
}

func changeUnitPrice(gin *MyContext) {
	var price UnitPriceParam

	if err := gin.BindJSON(&price); err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest())
		return
	}
	ri := gin.CoreContext.GetConfig().ChainRegInfo.ResourceIndex
	err := gin.CoreContext.ReportClient.ModifyResourcePrice(ri, int64(price.UnitPrice))
	if err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest("modify price fail"))
	} else {
		gin.JSON(http.StatusOK, Success(""))
	}
}
