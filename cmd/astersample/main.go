package main

import (
	"flag"
	"log"

	"github.com/yardbirdsax/aster"
)

func main() {
	dir := flag.String("d", "", "The directory to parse")
	flag.Parse()

	if dir == nil || *dir == "" {
		log.Fatal("must provide a value for the directory to be parsed")
	}

	results, err := aster.FromDirectory(*dir).MatchComment("aster:")
	if err != nil {
		log.Fatal("could not retrieve results: %w", err)
	}

	for _, r := range results {
		log.Println("result found:")
		log.Printf("\tname: %s, type: %s", r.Name, r.Type)
		log.Printf("\tComments:\n\t\t\t%s", r.Comments)
		log.Println("------------------------------------")
	}
}
