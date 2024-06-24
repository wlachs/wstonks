package main

import (
	"github.com/wlachs/wstonks/pkg/transaction"
	"github.com/wlachs/wstonks/pkg/transaction/io"
	"log"
	"os"
)

// main example function for using the CSV-parser
func main() {
	if len(os.Args) < 2 {
		log.Fatalln("missing file path arg")
	}

	path := os.Args[1]
	ctx := transaction.Context{}
	csv := io.TxCsvLoader{Path: path}
	err := csv.Load(&ctx)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(ctx)
}
