package corehttp

import (
	"fmt"
	"net/http"

	"github.com/hamster-shared/hamster-provider/core/modules/provider/optimism"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/polygon"

	"github.com/hamster-shared/hamster-provider/core/modules/provider"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/aptos"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/avalanche"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/ethereum"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/near"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/starkware"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/sui"
)

func pullImageChain(c *MyContext) {
	chain, err := getChain(c)
	if err != nil {
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

	switch deployType {
	case 1:
		return aptos.New(), nil
	case 2:
		return sui.New(), nil
	case 3:
		return ethereum.New(), nil
	case 4:
		return nil, fmt.Errorf("not support deployType %d", deployType)
	case 5:
		return polygon.New(), nil
	case 6:
		return avalanche.New(), nil
	case 7:
		return optimism.New(), nil
	case 8:
		return nil, fmt.Errorf("not support deployType %d", deployType)
	case 9:
		return starkware.New(), nil
	case 10:
		return near.New(), nil
	case 11:
		return nil, fmt.Errorf("not support deployType %d", deployType)
	default:
		return nil, fmt.Errorf("not support deployType %d", deployType)
	}
}
