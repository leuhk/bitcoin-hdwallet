package managers

import (
	"fmt"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Manager interface {
	GenerateMnemonic() (mnemonic string, err error)
}

type manager struct {
}

const bitSize = 256

func (m *manager) GenerateMnemonic() (mnemonic string, err error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}
	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	// // Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "Secret Passphrase")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	fmt.Println("Mnemonic: ", mnemonic)
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)
	return mnemonic, nil

}

func NewManager() Manager {
	return &manager{}
}
