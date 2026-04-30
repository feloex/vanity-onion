//go:build js && wasm

package main

import (
	"archive/zip"
	"bytes"
	"syscall/js"

	_ "golang.org/x/crypto/sha3"
)

func generateVanityOnionJS(this js.Value, args []js.Value) any {
	targetPrefix := ""
	if len(args) > 0 {
		targetPrefix = args[0].String()
	}

	onion, privateKey, publicKey := GenerateVanityOnion(targetPrefix, func() {
		js.Global().Call("postMessage", map[string]interface{}{
			"type":   "progress",
			"hashes": 10000,
		})
	})

	return map[string]interface{}{
		"onion":      onion,
		"privateKey": privateKey,
		"publicKey":  publicKey,
	}
}

func downloadKeysJS(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return nil
	}
	onion := args[0].String()
	privateKey := args[1].String()
	publicKey := args[2].String()

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	hostname, priv, pub, _ := GetExpandedSecrets(onion, privateKey, publicKey)

	f1, _ := zipWriter.Create(onion + "/hostname")
	f1.Write(hostname)
	f2, _ := zipWriter.Create(onion + "/hs_ed25519_secret_key")
	f2.Write(priv)
	f3, _ := zipWriter.Create(onion + "/hs_ed25519_public_key")
	f3.Write(pub)

	zipWriter.Close()

	zipBytes := buffer.Bytes()
	jsZipArray := js.Global().Get("Uint8Array").New(len(zipBytes))
	js.CopyBytesToJS(jsZipArray, zipBytes)

	return jsZipArray
}

func calculateStatsJS(this js.Value, args []js.Value) any {
	attempts := args[0].Int()
	prefixLen := args[1].Int()
	elapsedSeconds := args[2].Float()

	hashrate, effort := CalculateStats(attempts, prefixLen, elapsedSeconds)

	return map[string]interface{}{
		"hashrate": hashrate,
		"effort":   effort,
	}
}

func main() {
	js.Global().Set("downloadKeys", js.FuncOf(downloadKeysJS))
	js.Global().Set("generateVanityOnion", js.FuncOf(generateVanityOnionJS))
	js.Global().Set("calculateStats", js.FuncOf(calculateStatsJS))
	<-make(chan struct{})
}
