package memeduck_test

import (
	"fmt"

	"github.com/genkami/memeduck"
)

func ExampleUpdate() {
	query, _ := memeduck.Update("user").
		Set(memeduck.Ident("position"), "BOTTOM LEFT").
		Set(memeduck.Ident("immortal"), true).
		Where(
			memeduck.Eq(memeduck.Ident("color"), "orange"),
			memeduck.Eq(memeduck.Ident("manager"), true),
		).
		SQL()
	fmt.Println(query)
	// Output: UPDATE user SET position = "BOTTOM LEFT", immortal = TRUE WHERE color = "orange" AND manager = TRUE
}

func ExampleUpdate_multipleWhere() {
	query, _ := memeduck.Update("user").
		Set(memeduck.Ident("race"), "gorilla").
		Where(memeduck.Eq(memeduck.Ident("race"), "angel")).
		Where(memeduck.Ge(memeduck.Ident("grip_strength_kg"), 50)).
		SQL()
	fmt.Println(query)
	// Output: UPDATE user SET race = "gorilla" WHERE race = "angel" AND grip_strength_kg >= 50
}

func ExampleUpdate_queryParameter() {
	query, _ := memeduck.Update("user").
		Set(memeduck.Ident("age"), memeduck.Param("age")).
		Where(memeduck.Eq(memeduck.Ident("shark"), true)).
		SQL()
	fmt.Println(query)
	// Output: UPDATE user SET age = @age WHERE shark = TRUE
}
