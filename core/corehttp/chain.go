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
	chainName := c.Query("chainName")
	if chainName == "" {
		c.JSON(http.StatusBadRequest, BadRequest("not found chainName"))
		return
	}
	chain, err := getChain(chainName)
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
	chainName := c.Query("chainName")
	if chainName == "" {
		c.JSON(http.StatusBadRequest, BadRequest("not found chainName"))
		return
	}
	chain, err := getChain(chainName)
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
	chainName := c.Query("chainName")
	if chainName == "" {
		c.JSON(http.StatusBadRequest, BadRequest("not found chainName"))
		return
	}
	chain, err := getChain(chainName)
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
	chainName := c.Query("chainName")
	if chainName == "" {
		c.JSON(http.StatusBadRequest, BadRequest("not found chainName"))
		return
	}
	chain, err := getChain(chainName)
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

func getChain(name string) (provider.Chain, error) {
	switch name {
	case "aptos":
		return aptos.New(), nil
	case "avalanche":
		return avalanche.New(), nil
	case "ethereum":
		return ethereum.New(), nil
	case "near":
		return near.New(), nil
	case "starkware":
		return starkware.New(), nil
	case "sui":
		return sui.New(), nil
	case "polygon":
		return polygon.New(), nil
	case "optimism":
		return optimism.New(), nil
	default:
		return nil, fmt.Errorf("not found chain %s", name)
	}
}
