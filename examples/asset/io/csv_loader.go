package main

import (
	"github.com/wlachs/wstonks/pkg/asset"
	"github.com/wlachs/wstonks/pkg/asset/io"
	"log"
	"os"
)

// main example function for using the CSV-parser
func main() {
	if len(os.Args) < 2 {
		log.Fatalln("missing file path arg")
	}

	path := os.Args[1]
	ctx := asset.Context{}
	csv := io.LiveAssetCsvLoader{Path: path}
	err := csv.Load(&ctx)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(ctx)
}
