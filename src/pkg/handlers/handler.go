package handlers

import (
	"fmt"

	"btcwallet.com/src/pkg/managers"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetMnemonic(ctx *gin.Context)
	GetHdWalletAddress(ctx *gin.Context)
}

type handler struct {
	manager managers.Manager
}

func (h *handler) GetMnemonic(ctx *gin.Context) {

	mnemonic, err := h.manager.GenerateMnemonic()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(404, gin.H{
			"message": "Unable to generate mnemonic",
		})
	}
	ctx.JSON(200, gin.H{
		"mnemonic": mnemonic,
	})
}

func (h *handler) GetHdWalletAddress(ctx *gin.Context) {

}

func NewHandler(manager managers.Manager) Handler {
	return &handler{
		manager,
	}
}
