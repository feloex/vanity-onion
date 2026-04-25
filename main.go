package main

import (
	"crypto"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/crypto/ed25519"
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

// Save onion credentials in a directory named after the onion adress
func SaveOnionKeys(onionAddress string, privateKeyHex string, publicKeyHex string) error {
	dir := filepath.Join("keys", onionAddress)

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	publicBytes, _ := hex.DecodeString(publicKeyHex)
	privateSeed, _ := hex.DecodeString(privateKeyHex)
	privateBytes := TorExpandedSecretFromSeed(privateSeed)

	// hostname file
	os.WriteFile(filepath.Join(dir, "hostname"), []byte(onionAddress+"\n"), 0600)

	// Keypair files: header + 3 byte padding + key
	pubHeader := []byte("== ed25519v1-public: type0 ==\x00\x00\x00")
	os.WriteFile(filepath.Join(dir, "hs_ed25519_public_key"),
		append(pubHeader, publicBytes...), 0600)

	privHeader := []byte("== ed25519v1-secret: type0 ==\x00\x00\x00")
	os.WriteFile(filepath.Join(dir, "hs_ed25519_secret_key"),
		append(privHeader, privateBytes...), 0600)

	return nil
}

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

// notes

// 56 chars long
// private/public key pair
// ends with d for some reason
// ED25519 elliptic curve
// onion adress based on the public key

// checksum:
//  bytes of ".onion checksum" + public key + bytes of version number ("\x03")
//  then put this byte array in SHA3-256 and take first 2 bytes

// adress:
//  Public key + checksum + bytes(\x03) -> put in base32 and append ".onion"

//var version byte = '\x03'
