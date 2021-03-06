/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/environment"
	"os"

	"github.com/spf13/cobra"
)

// environmentsInfoCmd represents the info command
var environmentsInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about currently selected environment",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("enviroments info called")
		config, err := configuration.GetConfig()
		if err != nil {
			panic(fmt.Sprintf("get config failed: %v\n", err)) //TODO err
		}
		if config.CurrentEnvironment == uuid.Nil {
			fmt.Println("no environment used") //TODO warn?
			os.Exit(1)
		}
		environment, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			panic(fmt.Sprintf("environemtns details failed: %v\n", err)) //TODO err
		}
		fmt.Println(environment.String())
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
