package corehttp

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/hamster-shared/hamster-provider/core/context"
)

func StartApi(ctx *context.CoreContext) error {
	r := NewMyServer(ctx)
	r.GET("/ws", TerminalHandle)
	r.GET("/wslog", TerminalLogHandle)
	// router
	v1 := r.Group("/api/v1")
	{

		// basic configuration
		config := v1.Group("/config")
		{
			config.GET("/settting", getConfig)
			config.POST("/settting", setConfig)
			config.POST("/boot", setBootState)
			config.GET("/boot", getBootState)
		}
		chain := v1.Group("/chain")
		{
			chain.GET("/resource", getChainResource)
			chain.GET("/expiration-time", getCalculateInstanceOverdue)
			chain.GET("/account-info", getAccountInfo)
			chain.GET("/staking-info", getStakingInfo)
			chain.POST("/pledge", stakingAmount)
			chain.POST("/withdraw-amount", withdrawAmount)
			chain.POST("/price", changeUnitPrice)
		}
		// container routing
		container := v1.Group("/container")
		{
			container.GET("/start", startContainer)
			container.GET("/delete", deleteContainer)
		}

		p2p := v1.Group("/p2p")
		// p2p
		{
			p2p.POST("/listen", listenP2p)
			p2p.POST("/forward", forwardP2p)
			p2p.GET("/ls", lsP2p)
			p2p.POST("/close", closeP2p)
			p2p.POST("/check", checkP2p)
		}
		vm := v1.Group("/vm")
		{
			vm.POST("/create", createVm)
		}
		resource := v1.Group("/resource")
		{
			resource.POST("/modify-price", modifyPrice)
			resource.POST("/add-duration", addDuration)
			resource.POST("/receive-income", receiveIncome)
			resource.POST("/rent-again", rentAgain)
			resource.POST("/delete-resource", deleteResource)
			resource.GET("/receive-income-judge", receiveIncomeJudge)
		}
	}
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	port := ctx.GetConfig().ApiPort
	return r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
