package domains

import (
	"database/sql"
	"fmt"
)

// Test runs some mock actions
func Test(db *sql.DB) {
	// adding domains
	domainList := []Domain{
		Domain{
			Name: "name1",
			Data: "{}",
		},
		Domain{
			Name: "name2",
			Data: "{}",
		},
	}
	for _, d := range domainList {
		err := Add(db, &d)
		if err != nil {
			fmt.Println(err)
		}
	}

	// testing domain funcs
	UpdateData(db, 1, `{"updated": true}`)
}