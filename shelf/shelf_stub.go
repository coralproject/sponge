/*
Package shelf is the service layer to get data into our databases. This is a stub file until we have the real one.
*/
package shelf

import (
	"fmt"

	"github.com/coralproject/sponge/models"
)

// Add imports data into the collection collection in shelf
func Add(modelName string, data []models.Model) error {

	fmt.Printf("Push data from %s into shelf.", modelName)

	return nil
}
