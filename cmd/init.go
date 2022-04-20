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
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// initialize config
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init  config",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("init hamster provider")

		path := config.DefaultConfigPath()

		cfg := getDefaultConfig()

		err := os.MkdirAll(filepath.Dir(path), os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Chmod(filepath.Dir(path), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		// init config
		log.Info("init context")
		err = config.NewConfigManagerWithPath(path).Save(&cfg)
		if err != nil {
			log.Error(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getDefaultConfig() config.Config {
	identity, err := config.CreateIdentity()

	if err != nil {
		log.Error("create identity error")
		os.Exit(0)
	}

	return config.Config{
		ApiPort:      10771,
		Identity:     identity,
		Keys:         []config.PublicKey{},
		LinkApi:      CONFIG_DEFAULT_LINK_API,
		ChainApi:     CONFIG_DEFAULT_CHAIN_API,
		Vm:           getDockerDefaultConfig(),
		SeedOrPhrase: "betray extend distance category chimney globe employ scrap armor success kiss forum",
		ChainRegInfo: config.ChainRegInfo{},
		ConfigFlag:   config.NONE,
	}
}

func getKvmDefaultConfig() config.VmOption {
	return config.VmOption{
		Cpu:        1,
		Mem:        1,
		Disk:       50,
		System:     "Centos 7",
		Image:      "https://s3.ttchain.tntlinking.com/compute/CentOS7.qcow2.tar.gz",
		AccessPort: 22,
		Type:       "kvm",
	}
}

func getDockerDefaultConfig() config.VmOption {
	return config.VmOption{
		Cpu:        1,
		Mem:        1,
		Disk:       50,
		System:     "Ubuntu 18",
		Image:      "rastasheep/ubuntu-sshd:18.04",
		AccessPort: 22,
		Type:       "docker",
	}
}
