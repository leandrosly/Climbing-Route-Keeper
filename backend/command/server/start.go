package main

import (
	"crk/core/config"
	"crk/core/log"
	"fmt"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	mode       string
	configPath string
	startCmd   = &cobra.Command{
		Use:   "start",
		Short: "start server",
		Long:  `start server, default port is 5000`,
		Run:   startServer,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(startCmd)
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.crk/default.yaml)")
	startCmd.PersistentFlags().Int("port", 5000, "Port to run Application server on")
	config.Viper().BindPFlag("port", startCmd.PersistentFlags().Lookup("port"))
}

func initConfig() {
	config.Viper().SafeWriteConfig()
	config.Viper().WriteConfigAs("$HOME/.crk/.config")
	if len(configPath) != 0 {
		config.Viper().SetConfigFile(configPath)
	} else {
		config.Viper().AddConfigPath("$HOME/.crk/")
		config.Viper().AddConfigPath("./config")
		config.Viper().SetConfigName("default")
	}
	config.Viper().SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.Viper().AutomaticEnv()
	if err := config.Viper().ReadInConfig(); err != nil {
		log.Fatalf("Using config file [%s]: %v", config.Viper().ConfigFileUsed(), err)
	}
	log.Info("Config paths:", config.Viper().ConfigFileUsed())
	log.Info("DBConnection:", len(config.Viper().GetString("database.url")))
}

func startServer(cmd *cobra.Command, agrs []string) {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info("Start http-server")
	router := gin.Default()
	pprof.Register(router, "monitor/pprof")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(fmt.Sprintf(":%d", config.AllConf().Port))
}
