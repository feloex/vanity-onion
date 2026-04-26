//go:build js && wasm

package main

import (
	"syscall/js"

	_ "golang.org/x/crypto/sha3"
)

func generateVanityOnionJS(this js.Value, args []js.Value) any {
	targetPrefix := ""
	if len(args) > 0 {
		targetPrefix = args[0].String()
	}

	onion, privateKey, publicKey := GenerateVanityOnion(targetPrefix)

	return map[string]interface{}{
		"onion":      onion,
		"privateKey": privateKey,
		"publicKey":  publicKey,
	}
}

func main() {
	js.Global().Set("generateVanityOnion", js.FuncOf(generateVanityOnionJS))
	<-make(chan struct{})
}
