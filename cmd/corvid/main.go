package main

import (
	"github.com/CartConnoisseur/corvid/srv"
)

func main() {
	srv.Start()
	select {}
}
