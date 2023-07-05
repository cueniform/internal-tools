package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cueniform/internal-tools/ratchet"
)

func main() {

	// Usage: ratchet [imported cue file from json] [provider_address] [output directory]
	//
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [terraform-provider-schema.json] [provider_address]\n", os.Args[0])
		os.Exit(1)
	}
	JSONData, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	CUESchema := ratchet.EmitEntities(os.Args[2], JSONData)
	fmt.Println(CUESchema)
}
