## How to add a new source

#### Package Source
Any new source has to be included in the package source and implement the interface Sourcer:

// Sourcer is where the data is coming from (mysql, mongodb, api, postgresql, etc)
type Sourcer interface {
	GetData(string, int, int, string) ([]map[string]interface{}, error) //tableName, offset, limit, orderby
	GetQueryData(string, int, int, string, []string) ([]map[string]interface{}, error)
	GetTables() ([]string, error)
}

#### Include the source into source.New function
