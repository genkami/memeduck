package memeduck_test

import (
	"fmt"

	"github.com/genkami/memeduck"
)

type ExampleUserStruct struct {
	Name string `spanner:"UserName"`
	Papa string `spanner:"PapaName"`
}

func ExampleInsert_struct() {
	query, _ := memeduck.Insert("users", []string{"UserName", "PapaName"}).Values([]ExampleUserStruct{
		{Name: "Kiara", Papa: "huke"},
	}).SQL()
	fmt.Println(query)
	// Output: INSERT INTO users (UserName, PapaName) VALUES ("Kiara", "huke")
}
