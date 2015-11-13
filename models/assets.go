package models

import (
	"database/sql"
	"fmt"
)

// Taxonomy has information on taxonomy needed for the asset
type Taxonomy struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// Asset has the articles
type Asset struct {
	ID       string     `json:"id" bson:"_id"`
	URL      string     `json:"url" bson:"url"`
	Taxonomy []Taxonomy `json:"taxonomy" bson:"taxonomy"`
}

// Print only print information about the comment
func (a Asset) Print() {
	fmt.Println("Asset: ", a.ID, a.URL)
}

// Transform get the data from sd
func (a Asset) Transform(sd *sql.Rows) ([]Model, error) {
	var asset Asset
	var assets []Model

	for sd.Next() {
		err := sd.Scan(&asset.ID, &asset.URL)
		if err != nil {
			return nil, scanError{error: err}
		}

		n := len(assets)
		if len(assets) == cap(assets) {
			// Comments is full and we must expand
			// Double the size and add 1
			newAssets := make([]Model, len(assets), 2*len(assets)+1)
			copy(newAssets, assets)
			assets = newAssets
		}
		assets = assets[0 : n+1]
		assets[n] = asset
	}

	return assets, nil
}
