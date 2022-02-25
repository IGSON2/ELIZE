package main

import (
	"elizebch/cli"
	"elizebch/database"
)

func main() {
	defer database.Close()
	database.InitDB()
	cli.Start()
}
