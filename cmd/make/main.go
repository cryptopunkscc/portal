package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Portal make...")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	install := NewInstall(wd)
	install.Run(ParseArgs(os.Args[1:2]), os.Args[2:])
}
