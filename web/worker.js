importScripts("wasm_exec.js")

const go = new Go()
let wasmReady = WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance)
    //console.log("loaded wasm")
});

onmessage = async function(e) {
    await wasmReady
    
    if (e.data.action == "downloadKeys") {
        const zipArray = downloadKeys(e.data.onion, e.data.privateKey, e.data.publicKey);
        postMessage({type: "zip", zipArray: zipArray})
        return
    }

    //console.log("generating onion for prefix:", e.data.prefix)
    const result = generateVanityOnion(e.data.prefix)
    postMessage({type: "keys", result: result})
};