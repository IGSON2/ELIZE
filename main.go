package main

import (
	"elizebch/elizebch"
)

func main() {
	bch := elizebch.GetBlockchain()
	bch.AddBlock("Fourth")
	bch.AddBlock("Fifth")
	bch.AddBlock("Sixth")
}
