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

func GenerateVanityOnion(TargetPrefix string, callbackFunction func(int)) (string, string, string) {
	count := 0
	for {
		privateKey, publicKey := RandomKeyPair()
		TryOnion := OnionFromPublicKey(publicKey)

		count++

		if strings.HasPrefix(TryOnion, TargetPrefix) {
			if callbackFunction != nil {
				callbackFunction(count)
			}
			fmt.Println("\nfound: " + TryOnion)
			return TryOnion, privateKey, publicKey
		}

		if count >= 10000 {
			if callbackFunction != nil {
				callbackFunction(count)
			}
			count = 0
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

func GetExpandedSecrets(onionAddress string, privateKeyHex string, publicKeyHex string) ([]byte, []byte, []byte, error) {
	publicBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, nil, nil, err
	}

	privateSeed, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, nil, nil, err
	}

	privateBytes := TorExpandedSecretFromSeed(privateSeed)

	hostname := []byte(onionAddress + "\n")

	publicKeyWithHeader := append([]byte("== ed25519v1-public: type0 ==\x00\x00\x00"), publicBytes...)
	privateKeyWithHeader := append([]byte("== ed25519v1-secret: type0 ==\x00\x00\x00"), privateBytes...)

	return hostname, privateKeyWithHeader, publicKeyWithHeader, nil
}
