package main

import (
	"github.com/andiogenes/dvach/pkg/ui"
	"log"
)

func main() {
	if err := ui.InitTUI(); err != nil {
		log.Fatal(err)
	}
}
