package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/craton-api/chain/server/utils"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

const (
	checksumLength = 4
	version        = byte(0x00)
)

func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)

	versionedPubHash := append([]byte{version}, pubHash...)
	checksum := checksum(versionedPubHash)

	fullHash := append(versionedPubHash, checksum...)
	address := utils.Base58Encode(fullHash)

	return address
}

func checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))

	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	isValid := bytes.Equal(actualChecksum, targetChecksum)

	return isValid
}
