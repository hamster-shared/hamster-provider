/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/hamster-shared/hamster-provider/core"
	context2 "github.com/hamster-shared/hamster-provider/core/context"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"github.com/hamster-shared/hamster-provider/core/modules/listener"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/docker"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	"github.com/hamster-shared/hamster-provider/log"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"path/filepath"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("daemon called")

		context := NewContext()
		server := core.NewServer(context)
		server.Run()
	},
}

func NewContext() context2.CoreContext {
	cm := config.NewConfigManager()
	cfg, err := cm.GetConfig()
	if err != nil {
		log.GetLogger().Error(err)
		return context2.CoreContext{}
	}
	var vmManager provider.Manager
	// set vm template
	template := provider.Template{
		Cpu:    cfg.Vm.Cpu,
		Memory: cfg.Vm.Mem,
		System: cfg.Vm.System,
		Image:  cfg.Vm.Image,
	}
	if "docker" == cfg.Vm.Type {
		vmManager, err = docker.NewDockerManager(template)
	} else {
		//vmManager, err = vm2.NewVirtManager(template)
		os.Exit(1)
	}
	if err != nil {
		log.GetLogger().Error(err)
		return context2.CoreContext{}
	}

	timeService := utils.NewTimerService()

	ec := event.EventContext{
		VmManager:    vmManager,
		Cm:           cm,
		TimerService: timeService,
	}

	eventService := event.NewEventService(ec)

	context := context2.CoreContext{
		VmManager:     vmManager,
		Cm:            cm,
		TimerService:  timeService,
		EventService:  eventService,
		ChainListener: listener.NewChainListener(eventService, cm),
	}
	return context
}
func saveGatewayNodes(ctx context2.CoreContext) {
	cfg, err := ctx.Cm.GetConfig()
	if err != nil {
		fmt.Println("save gateway failed", err)
	} else {
		data, err := ctx.ReportClient.GetGatewayNodes()
		var nodes []string
		if err != nil {
			cfg.Bootstraps = nodes
		}
		if len(data) <= 3 {
			cfg.Bootstraps = data
		} else {
			num := rand.Intn(len(data) - 1)
			nodes = append(nodes, data[num])
			num1 := rand.Intn(len(data) - 1)
			nodes = append(nodes, data[num1])
			num3 := rand.Intn(len(data) - 1)
			nodes = append(nodes, data[num3])
			cfg.Bootstraps = nodes
		}
		path := config.DefaultConfigPath()

		err = os.MkdirAll(filepath.Dir(path), os.ModeDir)
		if err != nil {
			log.GetLogger().Error(err)
		}

		err = os.Chmod(filepath.Dir(path), os.ModePerm)
		if err != nil {
			log.GetLogger().Error(err)
		}
		err = config.NewConfigManagerWithPath(path).Save(cfg)
		if err != nil {
			log.GetLogger().Error(err)
		}
	}
}
func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
