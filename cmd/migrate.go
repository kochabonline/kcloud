/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kit/log"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "自动迁移数据库, 生成初始化数据",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		m, cleanup, err := initializeMigrate(config.Cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()
		m.Start()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
