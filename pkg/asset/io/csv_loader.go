package io

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"github.com/wlachs/wstonks/pkg/ioutils"
	"log"
)

// LiveAssetCsvLoader implements the LiveAssetLoader interface to allow importing context data from a CSV file.
type LiveAssetCsvLoader struct {
	Path string
}

// Load tries to parse the CSV file at Path and loads the data into the context.
func (l LiveAssetCsvLoader) Load(ctx *asset.Context) error {
	a, err := parseCsv(l.Path)
	if err != nil {
		return err
	}

	return ctx.AddAssets(a)
}

// parseCsv reads the CSV file at the given path and tries to convert it to a asset.Asset slice.
func parseCsv(path string) ([]*asset.Asset, error) {
	fileContent, err := ioutils.ReadCsvFile(path)
	if err != nil {
		return nil, err
	}

	assets := make([]*asset.Asset, 0, len(fileContent))
	if len(fileContent) == 0 {
		log.Println("the CSV file is empty")
		return assets, nil
	}

	for _, row := range fileContent {
		a, rowErr := readCsvRow(row)
		if rowErr != nil {
			return nil, rowErr
		}

		assets = append(assets, a)
	}

	return assets, nil
}

// readCsvRow converts a single entry of the CSV file to a model.Tx object.
func readCsvRow(row []string) (*asset.Asset, error) {
	assetId, err := parseAssetId(row[0])
	if err != nil {
		return nil, err
	}

	unitPrice, err := ioutils.ParseRat(row[1])
	if err != nil {
		return nil, err
	}

	return &asset.Asset{
		Id:        assetId,
		UnitPrice: unitPrice,
	}, nil
}

// parseAssetId validates the asset ID
func parseAssetId(s string) (string, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("missing asset ID")
	}

	return s, nil
}
