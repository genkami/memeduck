package memeduck_test

import (
	"fmt"

	"github.com/genkami/memeduck"
)

func ExampleDelete() {
	query, _ := memeduck.Delete("user").Where(
		memeduck.Eq(memeduck.Ident("id"), 123),
		memeduck.Eq(memeduck.Ident("unused"), true),
	).SQL()
	fmt.Println(query)
	// Output: DELETE FROM user WHERE id = 123 AND unused = TRUE
}

func ExampleDelete_multipleWhere() {
	query, _ := memeduck.Delete("user").
		Where(memeduck.Eq(memeduck.Ident("id"), 123)).
		Where(memeduck.Eq(memeduck.Ident("unused"), true)).
		SQL()
	fmt.Println(query)
	// Output: DELETE FROM user WHERE id = 123 AND unused = TRUE
}

func ExampleDelete_queryParameter() {
	query, _ := memeduck.Delete("user").
		Where(memeduck.Eq(memeduck.Ident("id"), memeduck.Param("id"))).
		SQL()
	fmt.Println(query)
	// Output: DELETE FROM user WHERE id = @id
}
