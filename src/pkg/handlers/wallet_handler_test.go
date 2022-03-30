package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"btcwallet.com/src/pkg/helpers"
	"btcwallet.com/src/pkg/managers"
	"github.com/gin-gonic/gin"
)

func TestGenerateMnemonic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var walletHelper helpers.WalletHelper = helpers.NewWalletHelper()
	var walletManager managers.WalletManager = managers.NewWalletManager(walletHelper)
	var walletHandler WalletHandler = NewWalletHandler(walletManager)
	var url string = "/util/mnemonic"

	r := gin.Default()
	r.POST(url, walletHandler.GenerateMnemonic)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}

func TestGenerateHdWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var walletHelper helpers.WalletHelper = helpers.NewWalletHelper()
	var walletManager managers.WalletManager = managers.NewWalletManager(walletHelper)
	var walletHandler WalletHandler = NewWalletHandler(walletManager)
	var url string = "/util/hd-wallet"
	body := &HdWallet{
		Path: "m / 44' / 0' / 0' / 0 / 0",
		Seed: "3fdf3c7c40ef678dd8950caac27f8006e27fdfab5e379ff7e3cef34a0226df830a49ca85476e7873e096ca6127d365995f6f135c71c27e8efe6cd1c497f6003f",
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	r := gin.Default()
	r.POST(url, walletHandler.GenerateHdWallet)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

}
func TestGenerateMultisignature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var walletHelper helpers.WalletHelper = helpers.NewWalletHelper()
	var walletManager managers.WalletManager = managers.NewWalletManager(walletHelper)
	var walletHandler WalletHandler = NewWalletHandler(walletManager)
	var url string = "/util/multi-sig-p2sh"
	var wif []string
	wif = append(wif, "cQpHXfs91s5eR9PWXui6qo2xjoJb2X3VdUspwKXe4A8Dybvut2rL")
	wif = append(wif, "cVgxEkRBtnfvd41ssd4PCsiemahAHidFrLWYoDBMNojUeME8dojZ")
	body := &Multisignature{
		N:   1,
		M:   2,
		Wif: wif,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	r := gin.Default()
	r.POST(url, walletHandler.GenerateMultisignature)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

}
