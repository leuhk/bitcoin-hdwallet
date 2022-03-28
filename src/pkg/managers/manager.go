package managers

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Manager interface {
	GenerateMnemonic(passPhrase string) (mnemonic string, seed string, err error)
	GenerateHdWallet(seed string, path string) (xPrvKey string, xPubKey string, rootKey string, err error)
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

func (m *manager) GenerateHdWallet(seed string, path string) (xPrvKey string, xPubKey string, rootKey string, err error) {
	decodeSeed, err := hex.DecodeString(seed)
	if err != nil {
		return "", "", "", err
	}
	master, err := bip32.NewMasterKey(decodeSeed)
	if err != nil {
		return "", "", "", err
	}

	params, err := deriveParamsFromPath(path)
	if err != nil {
		return "", "", "", err
	}
	child, err := master.NewChildKey(params.Purpose)
	if err != nil {
		return "", "", "", err
	}
	child, err = child.NewChildKey(params.CoinType)
	if err != nil {
		return "", "", "", err
	}

	child, err = child.NewChildKey(params.Account)
	if err != nil {
		return "", "", "", err
	}

	child, err = child.NewChildKey(params.Change)
	if err != nil {
		return "", "", "", err
	}
	child, err = child.NewChildKey(params.AddressIndex)
	if err != nil {
		return "", "", "", err
	}

	xPrvKey, xPubKey, rootKey = child.String(), child.PublicKey().String(), master.String()

	return xPrvKey, xPubKey, rootKey, nil

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

func NewManager() Manager {
	return &manager{}
}
