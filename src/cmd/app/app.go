package app

import (
	"btcwallet.com/src/pkg/handlers"
	"btcwallet.com/src/pkg/helpers"
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
				walletHelper  helpers.WalletHelper   = helpers.NewWalletHelper()
				walletManager managers.WalletManager = managers.NewWalletManager(walletHelper)
				walletHandler handlers.WalletHandler = handlers.NewWalletHandler(walletManager)
			)
			util := r.Group("/util")
			{
				util.POST("/mnemonic", func(ctx *gin.Context) {
					walletHandler.GenerateMnemonic(ctx)
				})
				util.POST("/hd-wallet", func(ctx *gin.Context) {
					walletHandler.GenerateHdWallet(ctx)
				})
				util.POST("multi-sig-p2sh", func(ctx *gin.Context) {
					walletHandler.GenerateMultisignature(ctx)
				})
			}
			r.Run()
			return nil
		},
	}
}
