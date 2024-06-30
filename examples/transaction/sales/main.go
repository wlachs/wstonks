package main

import (
	"github.com/wlachs/wstonks/pkg/transaction"
	"github.com/wlachs/wstonks/pkg/transaction/io"
	"log"
	"os"
)

// main example function for calculating realozed profits and losses
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

	profit := ctx.GetRealizedProfit()
	p, _ := profit.Float32()
	log.Printf("Overall profit\t:\t %f\n", p)

	loss := ctx.GetRealizedLoss()
	l, _ := loss.Float32()
	log.Printf("Overall loss\t:\t %f\n", l)
}
