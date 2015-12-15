package fiddler

import (
	"fmt"

	str "github.com/coralproject/sponge/strategy"
)

//* ASSETS *//

// Asset is embedding the model package to extend it
type Asset struct {
	//	model.Asset
	fields map[string]interface{}
}

// Print only print information about the model
func (a Asset) Print() {
	fmt.Println("Asset: ", a.fields["AssetID"], a.fields["URL"])
}

// Transform does the data transformation on the Asset
func (a Asset) Transform(sd []map[string]interface{}, table str.Table) ([]Transformer, error) {

	var asset Asset
	asset.fields = make(map[string]interface{})

	var assets []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coralField, f := range table.Fields {

			if f != "" { // convert field f with value value[f] into field coralField
				newValue := transformAssetField(f, value[f], coralField)

				if newValue != nil {
					asset.fields[coralField] = newValue
				}
			}
		}

		n := len(assets)
		if len(assets) == cap(assets) {
			// Comments is full and we must expand
			// Double the size and add 1
			newAssets := make([]Transformer, len(assets), 2*len(assets)+1)
			copy(newAssets, assets)
			assets = newAssets
		}
		assets = assets[0 : n+1]
		//model.Raw = strings.Split(raws, ",")
		assets[n] = asset
	}

	return assets, nil
}

//Here we transform the record into what we want (based on the configuration)
// 1. convert types (values are all strings) into the struct
func transformAssetField(sourceField string, oldValue interface{}, coralField string) interface{} {

	var newValue interface{}

	// Right now this the simpler thing to do. Needs to look more into reflect to do it
	// dinamically (if it is worth it) or merge this into the model package
	switch coralField {
	case "AssetID": //string
		newValue = oldValue
	case "SourceID": //string    `json:"parent_id" bson:"parent_d"`
		newValue = oldValue
	case "URL": //string    `json:"asset_id" bson:"asset_id"`
		newValue = oldValue
		// Taxonomies missing
	}

	return newValue
}
