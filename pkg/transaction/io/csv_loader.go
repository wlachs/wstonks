package io

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/ioutils"
	"github.com/wlachs/wstonks/pkg/transaction"
	"log"
	"strconv"
	"time"
)

// TxCsvLoader implements the TransactionLoader interface to allow importing context data from a CSV file.
type TxCsvLoader struct {
	Path string
}

// Load tries to parse the CSV file at Path and loads the data to the context's transaction history.
func (l TxCsvLoader) Load(ctx *transaction.Context) error {
	t, err := parseCsv(l.Path)
	if err != nil {
		return err
	}

	return ctx.AddTransactions(t)
}

// parseCsv reads the CSV file at the given path and tries to convert it to a transaction.Tx slice.
func parseCsv(path string) ([]transaction.Tx, error) {
	fileContent, err := ioutils.ReadCsvFile(path)
	if err != nil {
		return nil, err
	}

	tradeHistory := make([]transaction.Tx, 0, len(fileContent))
	if len(fileContent) == 0 {
		log.Println("the CSV file is empty")
		return tradeHistory, nil
	}

	for _, row := range fileContent {
		tradeEvent, rowErr := readCsvRow(row)
		if rowErr != nil {
			return nil, rowErr
		}

		tradeHistory = append(tradeHistory, tradeEvent)
	}

	return tradeHistory, nil
}

// readCsvRow converts a single entry of the CSV file to a transaction.Tx object.
func readCsvRow(row []string) (transaction.Tx, error) {
	// Timestamp
	ts, err := parseTimestamp(row[0])
	if err != nil {
		return transaction.Tx{}, err
	}

	// Asset ID
	assetId, err := parseAssetId(row[1])
	if err != nil {
		return transaction.Tx{}, err
	}

	// Transaction type enum
	tradeType, err := parseTradeType(row[2])
	if err != nil {
		return transaction.Tx{}, err
	}

	// Order quantity
	quantity, err := ioutils.ParseRat(row[3])
	if err != nil {
		return transaction.Tx{}, err
	}

	// Unit price
	unitPrice, err := ioutils.ParseRat(row[4])
	if err != nil {
		return transaction.Tx{}, err
	}

	return transaction.Tx{
		Position: transaction.Position{
			Timestamp: ts,
			Asset:     &transaction.Asset{Id: assetId},
			Quantity:  quantity,
			UnitPrice: unitPrice,
		},
		Type: tradeType,
	}, nil
}

// parseTimestamp reads a raw timestamp string and converts it to time.Time.
func parseTimestamp(ts string) (time.Time, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	tm := time.UnixMilli(i)
	return tm, nil
}

// parseAssetId validates the asset ID
func parseAssetId(s string) (string, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("missing asset ID")
	}

	return s, nil
}

// parseTradeType converts the context type string to transaction.TxType
func parseTradeType(tt string) (transaction.TxType, error) {
	switch tt {
	case "BUY":
		return transaction.BUY, nil
	case "SELL":
		return transaction.SELL, nil
	case "DIVIDEND":
		return transaction.DIVIDEND, nil
	default:
		return -1, fmt.Errorf("unsupported trade type")
	}
}
