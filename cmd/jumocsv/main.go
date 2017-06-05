package main

import (
	"os"
	"log"
	"strings"
	"github.com/antonaudition/jumocsv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no input given")
	}
	input := os.Args[1]
	if !strings.HasSuffix(input, ".csv") {
		log.Fatalf("unexpected input: %q", input)
	}
	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("exiting: %v", err)
	}
	defer f.Close()

	out, err := os.Create("Output.csv")
	if err != nil {
		log.Fatalf("exiting: %v", err)
	}
	defer out.Close()

	c := jumocsv.NewCollator(f)
	err = c.Collate(out)
	if err != nil {
		log.Fatalf("exiting: %v", err)
	}
}
