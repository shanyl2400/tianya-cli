package main

import (
	"shanyl2400/tianya/client"
	"shanyl2400/tianya/cui"
)

func main() {
	client := client.NewClient()
	ui := cui.NewCUI(client)

	ui.Start()
}
