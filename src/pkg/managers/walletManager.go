package managers

import (
	"encoding/hex"

	"btcwallet.com/src/pkg/helpers"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type WalletManager interface {
	GenerateMnemonic(passPhrase string) (mnemonic string, seed string, err error)
	GenerateHdWallet(seed string, path string) (extPrvKey string, extPubKey string, rootKey string, wif string, p2pkhAddress string, segwitBech32 string, segwitNested string, err error)
	GenerateMultisignature(n int8, m int8, wif []string) (address string, err error)
}

type walletManager struct {
	walletHelper helpers.WalletHelper
}

type BIP44Params struct {
	Purpose      uint32 `json:"purpose"`
	CoinType     uint32 `json:"coinType"`
	Account      uint32 `json:"account"`
	Change       uint32 `json:"change"`
	AddressIndex uint32 `json:"addressIndex"`
}

const bitSize = 256

func (wm *walletManager) GenerateMultisignature(n int8, m int8, wif []string) (address string, err error) {

	publicKeys, err := wm.walletHelper.DerivePubKeyFromWif(wif)
	if err != nil {
		return "", nil
	}
	minSignature, err := wm.walletHelper.DeriveOpcodes(m)
	if err != nil {
		return "", err
	}
	numOfPubKeys, err := wm.walletHelper.DeriveOpcodes(n)
	if err != nil {
		return "", err
	}
	address, err = wm.walletHelper.GenerateMultisignatureRedeemHash(numOfPubKeys, minSignature, publicKeys)
	if err != nil {
		return "", err
	}
	return address, nil
}

func (wm *walletManager) GenerateMnemonic(passPhrase string) (mnemonic string, seed string, err error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", "", err
	}
	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return "", "", err
	}
	seed = hex.EncodeToString(bip39.NewSeed(mnemonic, passPhrase))
	return mnemonic, seed, nil
}

func (wm *walletManager) GenerateHdWallet(seed string, path string) (extPrvKey string, extPubKey string, rootKey string, wif string, p2pkhAddress string, segwitBech32 string, segwitNested string, err error) {
	decodeSeed, err := hex.DecodeString(seed)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	master, err := bip32.NewMasterKey(decodeSeed)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	rootKey = master.String()

	params, err := wm.walletHelper.DeriveParamsFromPath(path)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}

	xPrvKey, xPubKey, err := wm.walletHelper.DeriveExtendedKeys(master, params)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	extPrvKey, extPubKey = xPrvKey.String(), xPubKey.PublicKey().String()

	prvKey := wm.walletHelper.DerivePrivateKeyFromBytes(xPrvKey.Key)

	wif, p2pkhAddress, segwitBech32, segwitNested, err = wm.walletHelper.DeriveAddress(prvKey)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}

	return extPrvKey, extPubKey, rootKey, wif, p2pkhAddress, segwitBech32, segwitNested, nil
}

func NewWalletManager(walletHelper helpers.WalletHelper) WalletManager {
	return &walletManager{
		walletHelper,
	}
}
