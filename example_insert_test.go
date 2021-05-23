package memeduck_test

import (
	"fmt"

	"github.com/genkami/memeduck"
)

func ExampleInsert() {
	query, _ := memeduck.Insert("user", []string{"name", "oshi_mark"}).Values([][]string{
		{"Subaru", ":ambulance:"},
		{"Watame", ":sheep:"},
	}).SQL()
	fmt.Println(query)
	// Output: INSERT INTO user (name, oshi_mark) VALUES ("Subaru", ":ambulance:"), ("Watame", ":sheep:")
}

func ExampleInsert_queryParameter() {
	query, _ := memeduck.Insert("user", []string{"name", "weight", "is_onion"}).Values([][]interface{}{
		{memeduck.Param("name"), memeduck.Param("weight"), memeduck.Param("is_onion")},
	}).SQL()
	fmt.Println(query)
	// Output: INSERT INTO user (name, weight, is_onion) VALUES (@name, @weight, @is_onion)
}
