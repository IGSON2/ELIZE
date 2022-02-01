package main

import (
	"elizebch/explorer"
	"elizebch/restapi"
)

func main() {
	go explorer.Start()
	restapi.Start()
}
