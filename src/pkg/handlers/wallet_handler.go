package handlers

import (
	"fmt"

	"btcwallet.com/src/pkg/managers"
	"github.com/gin-gonic/gin"
)

type WalletHandler interface {
	GenerateMnemonic(ctx *gin.Context)
	GenerateHdWallet(ctx *gin.Context)
	GenerateMultisignature(ctx *gin.Context)
}

type walletHandler struct {
	walletManager managers.WalletManager
}

type HdWallet struct {
	Path string `form:"path" json:"path" binding:"required"`
	Seed string `form:"seed" json:"seed" binding:"required"`
}

type Multisignature struct {
	N   int8     `form:"n" json:"n" binding:"required"`
	M   int8     `form:"m" json:"m" binding:"required"`
	Wif []string `form:"wif" json:"wif" binding:"required"`
}

func (wh *walletHandler) GenerateMultisignature(ctx *gin.Context) {
	var json Multisignature

	if err := ctx.ShouldBindJSON(&json); err != nil {
		fmt.Println(err)
		ctx.JSON(422, gin.H{"error": "Unprocessable Entity"})
		return
	}
	address, err := wh.walletManager.GenerateMultisignature(json.N, json.M, json.Wif)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(404, gin.H{
			"error": "Unable to generate address",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"address": address,
	})

}
func (wh *walletHandler) GenerateMnemonic(ctx *gin.Context) {
	var passphrase string = ""
	mnemonic, seed, err := wh.walletManager.GenerateMnemonic(passphrase)

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

func (wh *walletHandler) GenerateHdWallet(ctx *gin.Context) {
	var json HdWallet

	if err := ctx.ShouldBindJSON(&json); err != nil {
		fmt.Println(err)
		ctx.JSON(422, gin.H{"error": "Unprocessable Entity"})
		return
	}

	extPrvKey, extPubKey, rootKey, wif, p2pkhAddress, segwitBech32, segwitNested, err := wh.walletManager.GenerateHdWallet(json.Seed, json.Path)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(404, gin.H{
			"error": "Unable to generate HD wallet",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"bip32ExtendedPublicKey":  extPubKey,
		"bip32ExtendedPrivateKey": extPrvKey,
		"bip32RootKey":            rootKey,
		"WIF":                     wif,
		"p2pkhAddress":            p2pkhAddress,
		"segwitBech32":            segwitBech32,
		"segwitNested":            segwitNested,
	})
}

func NewWalletHandler(walletManager managers.WalletManager) WalletHandler {
	return &walletHandler{
		walletManager,
	}
}
