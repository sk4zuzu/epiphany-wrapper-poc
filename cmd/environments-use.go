/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/promptui"
	"github.com/spf13/cobra"
)

// environmentsUseCmd represents the use command
var environmentsUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Allows to select environment to be used",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use called")
		config, err := configuration.GetConfig()
		if err != nil {
			panic(fmt.Sprintf("get config failed: %v\n", err)) //TODO err
		}
		uuid, err := promptui.PromptForEnvironmentSelect("Environments") //TODO check if any environment exists
		if err != nil {
			panic(fmt.Sprintf("prompt for environment select failed: %v\n", err)) //TODO err
		}
		fmt.Printf("Choosen environment UUID is: %s\v", uuid)
		err = config.SetUsedEnvironment(uuid)
		if err != nil {
			panic(fmt.Sprintf("setting used environment failed: %v\n", err)) //TODO err
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsUseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsUseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsUseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
