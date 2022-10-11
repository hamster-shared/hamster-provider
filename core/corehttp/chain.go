package corehttp

import (
	"github.com/hamster-shared/hamster-provider/core/modules/factory"
	"net/http"

	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

func pullImageChain(c *MyContext) {
	chain, err := getChain(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err != chain.InitParam(c.Context) {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if err := chain.PullImage(); err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, Success("ok"))
	}
}

func startChain(c *MyContext) {
	chain, err := getChain(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := chain.Start(); err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, Success("ok"))
	}
}

func stopChain(c *MyContext) {
	chain, err := getChain(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := chain.Stop(); err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, Success("ok"))
	}
}

func getChainStatus(c *MyContext) {
	chain, err := getChain(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if status, err := chain.GetStatus(); err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, Success(status))
	}
}

func getChain(c *MyContext) (provider.Chain, error) {
	deployType := c.CoreContext.GetConfig().ChainRegInfo.DeployType

	return factory.GetChain(deployType)
}
