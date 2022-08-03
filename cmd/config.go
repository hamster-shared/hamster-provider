// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config ",
}

// bootstrap config
var (
	bootstrapCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "add,remove,change you bootstrap nodes",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// bootstrap sub command
	addBootCmd = &cobra.Command{
		Use:   "add",
		Short: "add Bootstrap node ",
		Run: func(cmd *cobra.Command, args []string) {
			cm := config.NewConfigManager()
			if args[0] != "" {
				c, err := cm.GetConfig()
				if err != nil {
					fmt.Println(err)
					return
				}

				c.Bootstraps = append(c.Bootstraps, args[0])
				err = cm.Save(c)
				if err != nil {
					log.GetLogger().Error(err)
					return
				}

			}
		},
	}
	rmBootstrapCmd = &cobra.Command{

		Use:   "rm",
		Short: "remove Bootstrap node ",
		Run: func(cmd *cobra.Command, args []string) {
			boot := args[0]
			cm := config.NewConfigManager()
			if boot != "" {
				c, err := cm.GetConfig()
				if err != nil {
					fmt.Println(err)
					return
				}

				var res []string
				for _, value := range c.Bootstraps {
					if boot != value {
						res = append(res, value)
					}
				}

				if len(res) == len(c.Bootstraps) {
					fmt.Println("the bootstrap node you want delete is not exists")
				} else {
					c.Bootstraps = res
					err = cm.Save(c)
					if err != nil {
						logrus.Error(err)
						return
					}
				}

			}

		},
	}
	clearBootstrapCmd = &cobra.Command{

		Use:   "clear",
		Short: "clear Bootstrap node ",
		Run: func(cmd *cobra.Command, args []string) {
			cm := config.NewConfigManager()
			c, err := cm.GetConfig()
			if err != nil {
				fmt.Println(err)
				return
			}
			c.Bootstraps = []string{}
			err = cm.Save(c)
			if err != nil {
				logrus.Error(err)
				return
			}

		},
	}
)

// linkApi config
var (
	linkApiCmd = &cobra.Command{
		Use:   "linkapi",
		Short: "set your linkApi",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// link api subcommand
	setLinkApiCmd = &cobra.Command{

		Use:   "set",
		Short: "add your linkApi",
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "" {
				cm := config.NewConfigManager()
				c, err := cm.GetConfig()
				if err != nil {
					fmt.Println(err)
					return
				}
				c.LinkApi = args[0]
				err = cm.Save(c)
				if err != nil {
					logrus.Error(err)
					return
				}

			}

		},
	}
)

// SeedOrPhrase  config
var (
	chainSeedCmd = &cobra.Command{
		Use:   "chainseed",
		Short: "set your chain seed or phrase",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// link api subcommand
	setChainSeedCmd = &cobra.Command{

		Use:   "set",
		Short: "set your chain seed or phrase",
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "" {
				cm := config.NewConfigManager()
				c, err := cm.GetConfig()
				if err != nil {
					fmt.Println(err)
					return
				}
				c.SeedOrPhrase = args[0]
				err = cm.Save(c)
				if err != nil {
					logrus.Error(err)
					return
				}

			}

		},
	}
)

var showCmd = &cobra.Command{

	Use:   "show",
	Short: "show bootstrap or linkApi",
	Run: func(cmd *cobra.Command, args []string) {

		cm := config.NewConfigManager()
		c, err := cm.GetConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("=====================")
		i := 0
		for _, value := range c.Bootstraps {
			if value != "" {
				fmt.Printf("the %d is :%s\n", i+1, value)
				i++
			}
		}
		if i == 0 {
			fmt.Println("Bootstrap is empty")
		}

		fmt.Println("=====================")
		if c.LinkApi != "" {
			fmt.Printf("LinkApi is :%s\n", c.LinkApi)
		} else {
			fmt.Println("LinkApi is empty")
		}

	},
}

func init() {
	// configCmd add root
	rootCmd.AddCommand(configCmd)

	// bootstrap add config
	configCmd.AddCommand(bootstrapCmd, linkApiCmd, chainSeedCmd, showCmd)

	// bootstrap
	bootstrapCmd.AddCommand(addBootCmd, rmBootstrapCmd, clearBootstrapCmd)

	// linkApi
	linkApiCmd.AddCommand(setLinkApiCmd)
	//  chainSeedCmd
	chainSeedCmd.AddCommand(setChainSeedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
