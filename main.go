package main

import (
	"eliz/blockchainn"
	"fmt"
)

func main() {
	elize := blockchainn.GetBlockchain()
	fmt.Println(elize.Blocks[0].Hash)
}
