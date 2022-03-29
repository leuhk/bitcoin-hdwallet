package managers

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Manager interface {
	GenerateMnemonic(passPhrase string) (mnemonic string, seed string, err error)
	GenerateHdWallet(seed string, path string) (extPrvKey string, extPubKey string, rootKey string, wif string, p2pkhAddress string, segwitBech32 string, segwitNested string, err error)
}

type manager struct {
}

type BIP44Params struct {
	Purpose      uint32 `json:"purpose"`
	CoinType     uint32 `json:"coinType"`
	Account      uint32 `json:"account"`
	Change       uint32 `json:"change"`
	AddressIndex uint32 `json:"addressIndex"`
}

const bitSize = 256

func (m *manager) GenerateMnemonic(passPhrase string) (mnemonic string, seed string, err error) {
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

func (m *manager) GenerateHdWallet(seed string, path string) (extPrvKey string, extPubKey string, rootKey string, wif string, p2pkhAddress string, segwitBech32 string, segwitNested string, err error) {
	decodeSeed, err := hex.DecodeString(seed)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	master, err := bip32.NewMasterKey(decodeSeed)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	rootKey = master.String()

	params, err := deriveParamsFromPath(path)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}

	xPrvKey, xPubKey, err := getExtendedKeys(master, params)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	extPrvKey, extPubKey = xPrvKey.String(), xPubKey.PublicKey().String()

	prvKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), xPrvKey.Key)

	wif, p2pkhAddress, segwitBech32, segwitNested, err = getAddress(prvKey)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}

	return extPrvKey, extPubKey, rootKey, wif, p2pkhAddress, segwitBech32, segwitNested, nil
}

func getAddress(prvKey *btcec.PrivateKey) (wif string, p2pkhAddress string, segwitBech32 string, segwitNested string, err error) {
	// generate the wif(wallet import format) string
	btcwif, err := btcutil.NewWIF(prvKey, &chaincfg.MainNetParams, true)
	if err != nil {
		return "", "", "", "", err
	}
	wif = btcwif.String()

	// generate a normal p2pkh address
	serializedPubKey := btcwif.SerializePubKey()
	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	p2pkhAddress = addressPubKey.EncodeAddress()

	// generate a normal p2wkh address from the pubkey hash
	witnessProg := btcutil.Hash160(serializedPubKey)
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	segwitBech32 = addressWitnessPubKeyHash.EncodeAddress()

	serializedScript, err := txscript.PayToAddrScript(addressWitnessPubKeyHash)
	if err != nil {
		return "", "", "", "", err
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	segwitNested = addressScriptHash.EncodeAddress()

	return wif, p2pkhAddress, segwitBech32, segwitNested, nil
}

func deriveParamsFromPath(path string) (*BIP44Params, error) {
	// Split path into params
	spl := strings.Split(path, "/")

	// Trim white spaces
	for i := range spl {
		spl[i] = strings.TrimSpace(spl[i])
	}

	if len(spl) != 6 {
		return nil, fmt.Errorf("path length is wrong. Expected 6, got %d", len(spl))
	}

	purpose, err := hardenedInt(spl[1])
	if err != nil {
		return nil, err
	}
	coinType, err := hardenedInt(spl[2])
	if err != nil {
		return nil, err
	}
	account, err := hardenedInt(spl[3])
	if err != nil {
		return nil, err
	}
	change, err := hardenedInt(spl[4])
	if err != nil {
		return nil, err
	}
	addressIdx, err := hardenedInt(spl[5])
	if err != nil {
		return nil, err
	}

	if spl[0] != "m" {
		return nil, fmt.Errorf("Invalid path")
	}

	// Validate path values
	if spl[1] != "44'" {
		return nil, fmt.Errorf("first field in path must be 44', got %s", spl[0])
	}

	if !isHardened(spl[2]) || !isHardened(spl[3]) {
		return nil,
			fmt.Errorf("second and third field in path must be hardened (ie. contain the suffix ', got %s and %s", spl[2], spl[3])
	}

	if isHardened(spl[4]) || isHardened(spl[5]) {
		return nil,
			fmt.Errorf("fourth and fifth field in path must not be hardened (ie. not contain the suffix ', got %s and %s", spl[3], spl[4])
	}

	if !(change == 0 || change == 1) {
		return nil, fmt.Errorf("change field can only be 0 or 1")
	}

	return &BIP44Params{
		Purpose:      purpose,
		CoinType:     coinType,
		Account:      account,
		Change:       change,
		AddressIndex: addressIdx,
	}, nil

}

func hardenedInt(field string) (uint32, error) {
	hasSuffix := strings.HasSuffix(field, "'")
	field = strings.TrimSuffix(field, "'")
	i, err := strconv.Atoi(field)

	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("fields must not be negative. got %d", i)
	}
	if hasSuffix {
		hardenedInt := bip32.FirstHardenedChild + uint32(i)
		return hardenedInt, nil
	}
	return uint32(i), nil
}

func isHardened(field string) bool {
	return strings.HasSuffix(field, "'")
}

func getExtendedKeys(master *bip32.Key, params *BIP44Params) (xPrvKey *bip32.Key, xPubKey *bip32.Key, err error) {
	child, err := master.NewChildKey(params.Purpose)
	if err != nil {
		return nil, nil, err
	}
	child, err = child.NewChildKey(params.CoinType)
	if err != nil {
		return nil, nil, err
	}

	child, err = child.NewChildKey(params.Account)
	if err != nil {
		return nil, nil, err
	}

	child, err = child.NewChildKey(params.Change)
	if err != nil {
		return nil, nil, err
	}
	child, err = child.NewChildKey(params.AddressIndex)
	if err != nil {
		return nil, nil, err
	}
	return child, child.PublicKey(), nil
}

func NewManager() Manager {
	return &manager{}
}
