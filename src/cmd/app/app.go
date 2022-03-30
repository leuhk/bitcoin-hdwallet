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
				util.POST("/mnemonic", func(ctx *gin.Context) {
					handler.GenerateMnemonic(ctx)
				})
				util.POST("/hd-wallet", func(ctx *gin.Context) {
					handler.GenerateHdWallet(ctx)
				})
				util.POST("multi-sig/p2sh", func(ctx *gin.Context) {
					handler.GenerateMultisignature(ctx)
				})
			}
			r.Run()
			return nil
		},
	}
}
