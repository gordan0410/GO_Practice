/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package server

import (
	"hw_thirteenth/server/servers/a_server"
	"hw_thirteenth/server/servers/b_server"
	"hw_thirteenth/server/servers/c_server"
	"log"

	"github.com/spf13/cobra"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		go b_server.B_server("3001")
		go c_server.C_server("3002")
		err := a_server.A_server("3000")
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(allCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
