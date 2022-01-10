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
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/hamster-shared/hamster-provider/core"
	context2 "github.com/hamster-shared/hamster-provider/core/context"
	chain2 "github.com/hamster-shared/hamster-provider/core/modules/chain"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/core/modules/p2p"
	"github.com/hamster-shared/hamster-provider/core/modules/pk"
	vm2 "github.com/hamster-shared/hamster-provider/core/modules/vm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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

		cm := config.NewConfigManager()
		pkManager := pk.NewManager(cm)
		config, err := cm.GetConfig()
		if err != nil {
			logrus.Error(err)
			return
		}
		p2pClient, err := p2p.NewP2pClient(34001, config.Identity.PrivKey, config.Identity.SwarmKey, config.Bootstraps)
		if err != nil {
			logrus.Error(err)
			return
		}
		var vmManager vm2.Manager
		if "docker" == config.Vm.Type {
			vmManager, err = vm2.NewDockerManager()
		} else {
			vmManager, err = vm2.NewVirtManager()
		}
		if err != nil {
			logrus.Error(err)
			return
		}

		substrateApi, err := gsrpc.NewSubstrateAPI(config.ChainApi)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		reportClient, err := chain2.NewChainClient(cm, substrateApi)
		if err != nil {
			logrus.Error(err)
			return
		}
		context := context2.CoreContext{
			P2pClient:    p2pClient,
			VmManager:    vmManager,
			Cm:           cm,
			PkManager:    pkManager,
			ReportClient: reportClient,
			SubstrateApi: substrateApi,
		}
		server := core.NewServer(context)
		server.Run()
	},
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
