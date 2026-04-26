//go:build !js

package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	_ "golang.org/x/crypto/ed25519"
	_ "golang.org/x/crypto/sha3"
)

func main() {
	targetPrefix := ""
	count := 1

	if len(os.Args) > 1 {
		targetPrefix = os.Args[1]
	}

	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n < 1 {
			fmt.Fprintln(os.Stderr, "usage: ./vanity-onion [prefix] [amount]")
			os.Exit(1)
		}
		count = n
	}

	for i := 0; i < count; i++ {
		onion, privateKey, publicKey := GenerateVanityOnion(targetPrefix)
		if err := SaveOnionKeys(onion, privateKey, publicKey); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// Save onion credentials in a directory named after the onion address
func SaveOnionKeys(onionAddress string, privateKeyHex string, publicKeyHex string) error {
	dir := filepath.Join("keys", onionAddress)

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	hostname, privateKeyWithHeader, publicKeyWithHeader, err := GetExpandedSecrets(onionAddress, privateKeyHex, publicKeyHex)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "hostname"), hostname, 0600); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "hs_ed25519_public_key"), publicKeyWithHeader, 0600); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "hs_ed25519_secret_key"), privateKeyWithHeader, 0600); err != nil {
		return err
	}

	return nil
}

func GetExpandedSecrets(onionAddress string, privateKeyHex string, publicKeyHex string) ([]byte, []byte, []byte, error) {
	/* Keypair files: header + 3 byte padding + key
	pubHeader := []byte("== ed25519v1-public: type0 ==\x00\x00\x00")
	os.WriteFile(filepath.Join(dir, "hs_ed25519_public_key"),
		append(pubHeader, publicBytes...), 0600)*/

	/*privHeader := []byte("== ed25519v1-secret: type0 ==\x00\x00\x00")
	os.WriteFile(filepath.Join(dir, "hs_ed25519_secret_key"),
		append(privHeader, privateBytes...), 0600)*/

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

// notes

// 56 chars long
// private/public key pair
// ends with d for some reason
// ED25519 elliptic curve
// onion address based on the public key

// checksum:
//  bytes of ".onion checksum" + public key + bytes of version number ("\x03")
//  then put this byte array in SHA3-256 and take first 2 bytes

// address:
//  Public key + checksum + bytes(\x03) -> put in base32 and append ".onion"

//var version byte = '\x03'
