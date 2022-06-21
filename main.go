package main

import (
	"elize/cli"
	"elize/database"
)

func main() {
	defer database.Close()
	database.InitDB()
	cli.Start()
}
