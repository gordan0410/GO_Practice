/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package api

import (
	"hw_thirteenth/api/apis"
	"log"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		con_addr, err := cmd.Flags().GetString("server_config")
		if err != nil {
			log.Println("config load error use default config")
			con_addr = "./configs/config.json"
		}
		apis.Load_config(con_addr)
		if len(args) > 0 {
			apis.Me(args[0])
		} else {
			apis.Me("1")
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
