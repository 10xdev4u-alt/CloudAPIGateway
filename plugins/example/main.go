package main

//go:wasmexport greet
func greet() {
	println("Hello from Wasm Plugin!")
}

func main() {}
