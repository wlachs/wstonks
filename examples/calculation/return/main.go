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

// main example function for calculating asset return
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

	returnPerAsset, err := worthCtx.GetAssetReturnMap()
	if err != nil {
		log.Fatalln(err)
	}

	for a, assetReturn := range returnPerAsset {
		w, _ := assetReturn.Float32()
		log.Printf("%s: %f\n", a.Id, w)
	}
}
