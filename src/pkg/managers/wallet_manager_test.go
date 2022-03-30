package managers

import (
	"strings"
	"testing"

	"btcwallet.com/src/pkg/helpers"
)

func TestGenerateMnemonic(t *testing.T) {
	var walletHelper helpers.WalletHelper = helpers.NewWalletHelper()
	var walletManager WalletManager = NewWalletManager(walletHelper)
	var passPhrase string = ""

	menmonic, _, _ := walletManager.GenerateMnemonic(passPhrase)
	menmonicArr := strings.Split(menmonic, " ")

	if len(menmonicArr) != 24 {
		t.Errorf("Test failed:  expected: %d received: %d ", 24, len(menmonicArr))
	}
}

func TestGenerateHdWallet(t *testing.T) {
	var walletHelper helpers.WalletHelper = helpers.NewWalletHelper()
	var walletManager WalletManager = NewWalletManager(walletHelper)
	var seed string = "3fdf3c7c40ef678dd8950caac27f8006e27fdfab5e379ff7e3cef34a0226df830a49ca85476e7873e096ca6127d365995f6f135c71c27e8efe6cd1c497f6003f"
	var path string = "m / 44' / 0' / 0' / 0 / 0"
	var expectedExtPrvKey string = "xprvA319vcXCEmKZe8ens22j5pGfY6WR2yEeBnPPMPE2CDpn4JaoXyYjsHWyDeDbXFXDWwuJAgbJve2772PRfVrY6jFUBj43JDbXMJ5EZQYKDhM"
	var expectedExtPubKey string = "xpub6FzWL84658srrcjFy3ZjSxDQ68LuSRxVZ1Jz9mddkZMkw6ux5WrzR5qT4wSsnG7zpfQFrAeQDeoRzec8xXy5FRz8ZDewDG3NV8nDFNjYrjZ"
	var expectedRootKey string = "xprv9s21ZrQH143K2pnPh3AEko6pZTqmoyFW3Kt8heSpwhSzSfSJP3T9rFome7xNkvk9GaW7M91QEvkbP22z6HwhvFqTtuisH5hHPTu5xDBQRkG"
	var expectedWif string = "L1dbCB2GPDDzxNJJ7s7h7a84NSXViBNwJeLWTS2gJMtdapfkhmg8"
	var expectedP2pkhAddress string = "1AKdhnB63swG2XpSuuXWP8MqP596zRJJCz"
	var expectedSegwitBech32 string = "bc1qvcl5rm6vz7eaxed9zrcftax0ywsc29y8zgj876"
	var expectedSegwitNested string = "3MTm6vsDYfyQCSTeex9ZHWYd9zkUetUn7c"

	extPrvKey, extPubKey, rootKey, wif, p2pkhAddress, segwitBech32, segwitNested, _ := walletManager.GenerateHdWallet(seed, path)

	if extPrvKey != expectedExtPrvKey {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedExtPrvKey, extPrvKey)
	}
	if extPubKey != expectedExtPubKey {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedExtPubKey, extPubKey)
	}
	if rootKey != expectedRootKey {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedRootKey, rootKey)
	}
	if wif != expectedWif {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedWif, wif)
	}
	if p2pkhAddress != expectedP2pkhAddress {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedP2pkhAddress, p2pkhAddress)
	}
	if segwitBech32 != expectedSegwitBech32 {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedSegwitBech32, segwitBech32)
	}
	if segwitNested != expectedSegwitNested {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedSegwitNested, segwitNested)
	}
}

func TestGenerateMultisignature(t *testing.T) {
	var walletHelper helpers.WalletHelper = helpers.NewWalletHelper()
	var walletManager WalletManager = NewWalletManager(walletHelper)
	var wif []string
	wif = append(wif, "cQpHXfs91s5eR9PWXui6qo2xjoJb2X3VdUspwKXe4A8Dybvut2rL")
	wif = append(wif, "cVgxEkRBtnfvd41ssd4PCsiemahAHidFrLWYoDBMNojUeME8dojZ")
	var n int8 = 1
	var m int8 = 2
	var expectedAddress string = "3MqSiHLbK6M8YUL8sXKiULeiRSvckJV74h"

	address, _ := walletManager.GenerateMultisignature(n, m, wif)

	if address != expectedAddress {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedAddress, address)
	}
}
