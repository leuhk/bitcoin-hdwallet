package helpers

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/tyler-smith/go-bip32"
)

func TestDeriveParamsFromPath(t *testing.T) {
	path := "m / 44' / 0' / 0' / 0 / 0"
	var walletHelper WalletHelper = NewWalletHelper()
	var expectedCoinType uint32 = 2147483648

	result, _ := walletHelper.DeriveParamsFromPath(path)
	if result.CoinType != expectedCoinType {
		t.Errorf("Test failed: input: %s  expected: %d received: %d ", path, expectedCoinType, result.CoinType)
	}
}
func TestDeriveExtendedKeys(t *testing.T) {
	var walletHelper WalletHelper = NewWalletHelper()
	var seed string = "3fdf3c7c40ef678dd8950caac27f8006e27fdfab5e379ff7e3cef34a0226df830a49ca85476e7873e096ca6127d365995f6f135c71c27e8efe6cd1c497f6003f"
	var path string = "m / 44' / 0' / 0' / 0 / 0"
	decodeSeed, _ := hex.DecodeString(seed)
	master, _ := bip32.NewMasterKey(decodeSeed)
	params, _ := walletHelper.DeriveParamsFromPath(path)
	xPrvKey, XPubKey, _ := walletHelper.DeriveExtendedKeys(master, params)

	var expectedPrvKey string = "xprvA319vcXCEmKZe8ens22j5pGfY6WR2yEeBnPPMPE2CDpn4JaoXyYjsHWyDeDbXFXDWwuJAgbJve2772PRfVrY6jFUBj43JDbXMJ5EZQYKDhM"
	var expectedPubKey string = "xpub6FzWL84658srrcjFy3ZjSxDQ68LuSRxVZ1Jz9mddkZMkw6ux5WrzR5qT4wSsnG7zpfQFrAeQDeoRzec8xXy5FRz8ZDewDG3NV8nDFNjYrjZ"

	if xPrvKey.String() != expectedPrvKey {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedPrvKey, xPrvKey.String())
	}
	if XPubKey.PublicKey().String() != expectedPubKey {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedPubKey, XPubKey.PublicKey().String())
	}
}

func TestDeriveAddress(t *testing.T) {
	var walletHelper WalletHelper = NewWalletHelper()
	var seed string = "3fdf3c7c40ef678dd8950caac27f8006e27fdfab5e379ff7e3cef34a0226df830a49ca85476e7873e096ca6127d365995f6f135c71c27e8efe6cd1c497f6003f"
	var path string = "m / 44' / 0' / 0' / 0 / 0"
	decodeSeed, _ := hex.DecodeString(seed)
	master, _ := bip32.NewMasterKey(decodeSeed)
	params, _ := walletHelper.DeriveParamsFromPath(path)
	xPrvKey, _, _ := walletHelper.DeriveExtendedKeys(master, params)
	prvKey := walletHelper.DerivePrivateKeyFromBytes(xPrvKey.Key)
	wif, p2pkhAddress, segwitBech32, segwitNested, _ := walletHelper.DeriveAddress(prvKey)

	var expectedWif string = "L1dbCB2GPDDzxNJJ7s7h7a84NSXViBNwJeLWTS2gJMtdapfkhmg8"
	var expectedP2pkhAddress string = "1AKdhnB63swG2XpSuuXWP8MqP596zRJJCz"
	var expectedSegwitBech32 string = "bc1qvcl5rm6vz7eaxed9zrcftax0ywsc29y8zgj876"
	var expectedSegwitNested string = "3MTm6vsDYfyQCSTeex9ZHWYd9zkUetUn7c"

	if wif != expectedWif {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedWif, wif)
	}
	if p2pkhAddress != expectedP2pkhAddress {
		t.Errorf("Test failed:  expected: %s received: %s ", p2pkhAddress, expectedP2pkhAddress)
	}
	if segwitBech32 != expectedSegwitBech32 {
		t.Errorf("Test failed:  expected: %s received: %s ", segwitBech32, expectedSegwitBech32)
	}
	if segwitNested != expectedSegwitNested {
		t.Errorf("Test failed:  expected: %s received: %s ", segwitNested, expectedSegwitNested)
	}
}

func TestDeriveOpcodes(t *testing.T) {
	var walletHelper WalletHelper = NewWalletHelper()
	var n int8 = 2

	result, _ := walletHelper.DeriveOpcodes(n)

	if result != txscript.OP_2 {
		t.Errorf("Test failed:  expected: %d received: %d ", txscript.OP_2, result)
	}
}

func TestGenerateMultisignatureRedeemHash(t *testing.T) {
	var walletHelper WalletHelper = NewWalletHelper()
	var wif []string
	wif = append(wif, "cQpHXfs91s5eR9PWXui6qo2xjoJb2X3VdUspwKXe4A8Dybvut2rL")
	wif = append(wif, "cVgxEkRBtnfvd41ssd4PCsiemahAHidFrLWYoDBMNojUeME8dojZ")
	var n int8 = 1
	var m int8 = 2

	numOfPubKeys, _ := walletHelper.DeriveOpcodes(n)
	minSignature, _ := walletHelper.DeriveOpcodes(m)

	publicKeys, _ := walletHelper.DerivePubKeyFromWif(wif)

	result, _ := walletHelper.GenerateMultisignatureRedeemHash(numOfPubKeys, minSignature, publicKeys)

	var expectedResult string = "3MqSiHLbK6M8YUL8sXKiULeiRSvckJV74h"

	if result != expectedResult {
		t.Errorf("Test failed:  expected: %s received: %s ", expectedResult, result)
	}
}
