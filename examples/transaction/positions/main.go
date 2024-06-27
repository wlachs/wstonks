package main

import (
	"github.com/wlachs/wstonks/pkg/transaction"
	"github.com/wlachs/wstonks/pkg/transaction/io"
	"log"
	"os"
)

// pos is a readable position struct
type pos struct {
	timestamp string
	quantity  float32
	unitPrice float32
}

// main example function for calculating the chronological asset position summary
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

	assets := ctx.GetAssetPositionSliceMap()
	for asset, positions := range assets {
		ps := make([]pos, 0, len(positions))
		for _, position := range positions {
			qty, _ := position.Quantity.Float32()
			up, _ := position.UnitPrice.Float32()
			p := pos{
				timestamp: position.Timestamp.String(),
				quantity:  qty,
				unitPrice: up,
			}
			ps = append(ps, p)
		}
		log.Printf("%s: %v\n", asset.Id, ps)
	}
}
