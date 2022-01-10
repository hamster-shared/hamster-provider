package event

import (
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/context"
	log "github.com/sirupsen/logrus"
	"time"
)

func successDealOrder(ctx context.CoreContext, orderNo uint64, name string) error {
	err := forwardSSHToP2p(ctx, name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// report heartbeat
	agreementIndex := ctx.GetConfig().ChainRegInfo.AgreementIndex
	_ = ctx.ReportClient.Heartbeat(agreementIndex)

	// send timed heartbeats
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		ctx.TimerService.SubTicker(agreementIndex, ticker)
		for {
			<-ticker.C
			// report heartbeat
			agreementIndex := ctx.GetConfig().ChainRegInfo.AgreementIndex
			_ = ctx.ReportClient.Heartbeat(agreementIndex)
		}
	}()

	dealOverdueOrder(ctx, agreementIndex, name)
	return nil
}

func forwardSSHToP2p(ctx context.CoreContext, name string) error {
	// P2P listen port exposure
	ip, err := ctx.VmManager.GetIp(name)
	if err != nil {
		log.Error(err)
		return err
	}

	targetOpt := fmt.Sprintf("/ip4/%s/tcp/%d", ip, ctx.VmManager.GetAccessPort(name))
	err = ctx.P2pClient.Listen(targetOpt)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func dealOverdueOrder(ctx context.CoreContext, agreementIndex uint64, name string) bool {
	// calculate instance expiration time
	overdue := ctx.ReportClient.CalculateInstanceOverdue(ctx.GetConfig().ChainRegInfo.OrderIndex)
	instanceTimer := time.NewTimer(overdue)
	ctx.TimerService.SubTimer(agreementIndex, instanceTimer)

	go func(t *time.Timer) {
		<-t.C
		cfg := ctx.GetConfig()

		// expires triggers close
		_ = ctx.VmManager.Stop(name)
		_ = ctx.VmManager.Destroy(name)
		// modify the resource status on the chain to unused
		_ = ctx.ReportClient.ChangeResourceStatus(ctx.GetConfig().ChainRegInfo.ResourceIndex)
		// delete agreement number
		cfg.ChainRegInfo.OrderIndex = 0
		cfg.ChainRegInfo.AgreementIndex = 0
		cfg.ChainRegInfo.RenewOrderIndex = 0
		_ = ctx.Cm.Save(cfg)
	}(instanceTimer)

	return overdue < 0
}
