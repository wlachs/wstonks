package main

import (
	"github.com/wlachs/wstonks/pkg/asset"
	assetio "github.com/wlachs/wstonks/pkg/asset/io"
	"github.com/wlachs/wstonks/pkg/calculation"
	"github.com/wlachs/wstonks/pkg/transaction"
	txio "github.com/wlachs/wstonks/pkg/transaction/io"
	"log"
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

	profits, losses, err := worthCtx.GetMaxProfitAndLossForAssets()
	if err != nil {
		log.Fatalln(err)
	}

	for a := range profits {
		p, _ := profits[a].Float32()
		l, _ := losses[a].Float32()
		log.Printf("%s: %f - %f\n", a.Id, l, p)
	}
}
