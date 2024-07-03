package main

import (
	"github.com/wlachs/wstonks/pkg/asset"
	assetio "github.com/wlachs/wstonks/pkg/asset/io"
	"github.com/wlachs/wstonks/pkg/calculation"
	"github.com/wlachs/wstonks/pkg/transaction"
	txio "github.com/wlachs/wstonks/pkg/transaction/io"
	"log"
	"math/big"
	"os"
)

// main example function for calculating asset sales profit and loss
func main() {
	if len(os.Args) < 3 {
		log.Fatalln("missing file path arg(s)")
	}

	txCsvPath := os.Args[1]
	txCtx := transaction.Context{}
	txCsv := txio.TxCsvLoader{Path: txCsvPath}
	err := txCsv.Load(&txCtx)

	if err != nil {
		log.Fatalln(err)
	}

	assetCsvPath := os.Args[2]
	assetCtx := asset.Context{}
	assetCsv := assetio.LiveAssetCsvLoader{Path: assetCsvPath}
	err = assetCsv.Load(&assetCtx)

	if err != nil {
		log.Fatalln(err)
	}

	worthCtx := calculation.Context{
		AssetContext:       &assetCtx,
		TransactionContext: &txCtx,
	}

	profits, losses, err := worthCtx.GetMaxProfitAndLoss()
	if err != nil {
		log.Fatalln(err)
	}

	for a := range profits {
		p, _ := profits[a].Float32()
		l, _ := losses[a].Float32()
		log.Printf("%12s:\t%f\t-\t%f\n", a.Id, l, p)
	}

	log.Println("------------")

	m, err := worthCtx.GetSalesForReturn(big.NewRat(400, 1))
	if err != nil {
		log.Fatalln(err)
	}

	for a, sell := range m {
		s, _ := sell.Float32()
		log.Printf("%12s:\t%f\n", a.Id, s)
	}
}
