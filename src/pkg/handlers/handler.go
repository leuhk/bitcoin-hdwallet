package handlers

import (
	"fmt"

	"btcwallet.com/src/pkg/managers"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GenerateMnemonic(ctx *gin.Context)
	GenerateHdWallet(ctx *gin.Context)
}

type handler struct {
	manager managers.Manager
}

type HdWallet struct {
	Path string `form:"path" json:"path" binding:"required"`
	Seed string `form:"seed" json:"seed" binding:"required"`
}

type Mnemonic struct {
	Passphrase string `form:"passphrase" json:"passphrase"`
}

func (h *handler) GenerateMnemonic(ctx *gin.Context) {
	var json Mnemonic

	if err := ctx.ShouldBindJSON(&json); err != nil {
		fmt.Println(err)
		ctx.JSON(422, gin.H{"error": "Unprocessable Entity"})
		return
	}
	mnemonic, seed, err := h.manager.GenerateMnemonic(json.Passphrase)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(404, gin.H{
			"error": "Unable to generate mnemonic",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"BIP39Mnemonic": mnemonic,
		"BIP39Seed":     seed,
	})
}

func (h *handler) GenerateHdWallet(ctx *gin.Context) {
	var json HdWallet

	if err := ctx.ShouldBindJSON(&json); err != nil {
		fmt.Println(err)
		ctx.JSON(422, gin.H{"error": "Unprocessable Entity"})
		return
	}

	xPrvKey, xPubKey, rootKey, err := h.manager.GenerateHdWallet(json.Seed, json.Path)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(404, gin.H{
			"error": "Unable to generate HD wallet",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"BIP32ExtendedPublicKey":  xPubKey,
		"BIP32ExtendedPrivateKey": xPrvKey,
		"BIP32RootKey":            rootKey,
	})

}

func NewHandler(manager managers.Manager) Handler {
	return &handler{
		manager,
	}
}
