package fiddler

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
)

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
	fmt.Println("Asset: ", a.ID, a.URL)
}

// Transform get the data from sd
func (a Asset) Transform(sd *sql.Rows, table configuration.Table) ([]Transformer, error) {
	var asset Asset
	var assets []Transformer

	// Needs to check on what actually my sd (table) fields are based on the configuration for Table Asset.
	if table.Fields["Raw"] != "" {
		raws := strings.Split(table.Fields["Raw"], ",") // all the raw fields. Needs to see how to incorporate them into the scan columns
		fmt.Println(raws)
	}

	// work around to parse the dateTime values
	var createDate string
	var updateDate string
	var vendorID string

	for sd.Next() {

		// Scaning the columns
		err := sd.Scan(&asset.ID, &vendorID, &asset.SourceID, &asset.URL, &createDate, &updateDate)

		if err != nil {
			return nil, scanError{error: err}
		}
		// asset.CreateDate, _ = time.Parse("2006-01-02", createDate) // To Do: I need to see how to dinamically discover what is the dateTime layout
		// asset.UpdateDate, _ = time.Parse("2006-01-02", updateDate)

		// Resize array
		n := len(assets)
		if len(assets) == cap(assets) {
			// Assets is full and we must expand
			// Double the size and add 1
			newAssets := make([]Transformer, len(assets), 2*len(assets)+1)
			copy(newAssets, assets)
			assets = newAssets
		}
		assets = assets[0 : n+1]
		assets[n] = asset
	}

	return assets, nil
}
