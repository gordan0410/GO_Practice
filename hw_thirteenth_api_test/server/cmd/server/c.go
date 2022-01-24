/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package server

import (
	"hw_thirteenth/server/servers/c_server"
	"log"

	"github.com/spf13/cobra"
)

// cCmd represents the c command
var cCmd = &cobra.Command{
	Use:   "c",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Println(err)
		}
		err = c_server.C_server(port)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cCmd.Flags().StringP("port", "p", "3002", "server's port")

}
