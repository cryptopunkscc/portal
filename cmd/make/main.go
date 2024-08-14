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
	install.Run(parseArgs(os.Args[1:]))
}

func parseArgs(args []string) (m Make, a []string) {
	if len(args) == 0 {
		m = All
		a = []string{"linux", "windows"}
	}
	if len(args) > 0 {
		m = ParseMake(args[:1])
	}
	if len(args) > 1 {
		a = args[1:]
	}
	return
}
