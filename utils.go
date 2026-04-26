package main

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
)

func GenerateVanityOnion(TargetPrefix string) (string, string, string) {
	for {
		privateKey, publicKey := RandomKeyPair()
		TryOnion := OnionFromPublicKey(publicKey)
		if strings.HasPrefix(TryOnion, TargetPrefix) {
			fmt.Println(TryOnion)
			return TryOnion, privateKey, publicKey
		}
	}
}

func RandomKeyPair() (privateHex string, publicHex string) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(privateKey.Seed()), hex.EncodeToString(publicKey)
}

func TorExpandedSecretFromSeed(seed []byte) []byte {
	expanded := sha512.Sum512(seed)
	expanded[0] &= 248
	expanded[31] &= 63
	expanded[31] |= 64
	return expanded[:]
}

func OnionFromPublicKey(PublicKeyHex string) string {
	pubKeyBytes, err := hex.DecodeString(PublicKeyHex)
	if err != nil {
		panic(err)
	}
	version := byte(0x03)

	hash := crypto.SHA3_256.New()
	hash.Write([]byte(".onion checksum"))
	hash.Write(pubKeyBytes)
	hash.Write([]byte("\x03"))
	checksum := hash.Sum(nil)[:2]

	byteOnion := slices.Concat(pubKeyBytes, checksum, []byte{version})

	onion := string(base32.StdEncoding.EncodeToString(byteOnion))
	onion = strings.ToLower(onion)
	onion = strings.TrimRight(onion, "=")
	onion = onion + ".onion"

	return onion
}
