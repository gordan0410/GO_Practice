/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package api

import (
	"hw_thirteenth/api/apis"
	"log"

	"github.com/spf13/cobra"
)

// v2Cmd represents the v2 command
var v2Cmd = &cobra.Command{
	Use:   "v2",
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
			apis.Load_config("./configs/config.json")
		}
		apis.Load_config(con_addr)
		if len(args) > 1 {
			if args[0] == "get" {
				apis.V2_get(args[1])
			} else if args[0] == "post" {
				apis.V2_post(args[1])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(v2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// v2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// v2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
