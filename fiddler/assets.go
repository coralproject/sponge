package fiddler

import (
	"fmt"
	"log"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
	"github.com/oleiade/reflections"
)

//* ASSETS *//

// Asset is embedding the comment package to extend it
type Asset struct {
	comment.Asset
}

//// Taxonomy has information on taxonomy needed for the asset
// type Taxonomy struct {
// 	Name  string `json:"name" bson:"name"`
// 	Value string `json:"value" bson:"value"`
// }
//
// // Asset has the articles
// type Asset struct {
// 	ID       string `json:"id" bson:"_id"`
// 	VendorID string `json:"vendorid" bson:"vendorid"`
// 	SourceID string `json:"sourceid" bson:"sourceid"`
// 	URL      string `json:"url" bson:"url"`
// 	//Taxonomy   []Taxonomy `json:"taxonomy" bson:"taxonomy"`
// 	CreateDate time.Time `json:"createdate" bson:"createdate"`
// 	UpdateDate time.Time `json:"updatedate" bson:"updatedate"`
// 	Raw        []string  `json:"raws" bson:"raws"`
// }

// Print only print information about the comment
func (a Asset) Print() {
	fmt.Println("Asset: ", a.AssetID, a.URL)
}

// Transform does the data transformation on the Asset
func (a Asset) Transform(sd []map[string]interface{}, table configuration.Table) ([]Transformer, error) {

	var asset Asset
	var assets []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coralField, f := range table.Fields {

			if f != "" { // convert field f with value value[f] into field coralField
				newValue := transformAssetField(f, value[f], coralField)

				if newValue != nil {
					err := reflections.SetField(&asset, coralField, newValue)
					if err != nil {
						log.Fatal(err)
						return nil, err
					}
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
		//comment.Raw = strings.Split(raws, ",")
		assets[n] = asset
	}

	return assets, nil
}

//Here we transform the record into what we want (based on the configuration)
// 1. convert types (values are all strings) into the struct
func transformAssetField(sourceField string, oldValue interface{}, coralField string) interface{} {

	var newValue interface{}

	const longForm = "2015-11-02 12:26:05"

	// Right now this the simpler thing to do. Needs to look more into reflect to do it
	// dinamically (if it is worth it) or merge this into the comment package
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
