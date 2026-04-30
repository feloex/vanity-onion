//go:build !js

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "golang.org/x/crypto/ed25519"
	_ "golang.org/x/crypto/sha3"
)

var (
	totalAttempts int
	startTime     time.Time
	targetLength  int
)

func main() {
	targetPrefix := ""
	count := 1

	if len(os.Args) > 1 {
		targetPrefix = os.Args[1]
	}
	targetLength = len(targetPrefix)

	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n < 1 {
			fmt.Fprintln(os.Stderr, "usage: ./vanity-onion [prefix] [amount]")
			os.Exit(1)
		}
		count = n
	}

	for i := 0; i < count; i++ {
		totalAttempts = 0
		startTime = time.Now()
		onion, privateKey, publicKey := GenerateVanityOnion(targetPrefix, logProgress)
		if err := SaveOnionKeys(onion, privateKey, publicKey); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func logProgress() {
	totalAttempts += 10000
	elapsed := time.Since(startTime).Seconds()
	hashrate, effort := CalculateStats(totalAttempts, targetLength, elapsed)
	fmt.Printf("\r%.2fh/s effort: %.2f%%", hashrate, effort)
} //attempts: %d  ... ,totalAttempts

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
