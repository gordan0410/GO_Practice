/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package server

import (
	"hw_thirteenth/server/servers/a_server"
	"log"

	"github.com/spf13/cobra"
)

// aCmd represents the a command
var aCmd = &cobra.Command{
	Use:   "a",
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
		err = a_server.A_server(port)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(aCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// aCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	aCmd.Flags().StringP("port", "p", "3000", "server's port")
}
