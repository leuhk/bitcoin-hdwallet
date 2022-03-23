package app

import (
	"btcwallet.com/src/pkg/handlers"
	"btcwallet.com/src/pkg/managers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start Rest server",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := gin.Default()

			r.Use(gin.Recovery())

			var (
				manager managers.Manager = managers.NewManager()
				handler handlers.Handler = handlers.NewHandler(manager)
			)
			util := r.Group("/util")
			{
				util.GET("/mnemonic", func(ctx *gin.Context) {
					handler.GetMnemonic(ctx)
				})
				util.POST("/hd-wallet-address", func(ctx *gin.Context) {
					handler.GetHdWalletAddress(ctx)
				})
			}
			r.Run()
			return nil
		},
	}
}
