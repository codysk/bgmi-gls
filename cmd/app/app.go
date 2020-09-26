package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"moe.two.bgmi-gls/pkg/common"
	"moe.two.bgmi-gls/pkg/http"
)

var RootCmd = &cobra.Command{
	Use: "bgmi-gls",
	Run: func(cmd *cobra.Command, args []string) {
		config := &common.GLSConfig{
			Port: viper.GetString("PORT"),
		}
		common.Init(config)

		server := http.NewServer(common.Config.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	viper.SetDefault("PORT", ":8080")

	viper.AutomaticEnv()
}
