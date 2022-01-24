/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package server

import (
	"hw_thirteenth/server/servers/b_server"
	"log"

	"github.com/spf13/cobra"
)

// bCmd represents the b command
var bCmd = &cobra.Command{
	Use:   "b",
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
		err = b_server.B_server(port)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(bCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	bCmd.Flags().StringP("port", "p", "3001", "server's port")

}
