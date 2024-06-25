package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	log.Println("Starting install...")
	jobs := All
	if len(os.Args) > 1 {
		arg, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		jobs = Make(arg)
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	install := NewInstall(wd)
	install.Run(jobs)
}
