package memeduck_test

import (
	"fmt"

	"github.com/genkami/memeduck"
)

func ExampleSelect() {
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).Where(
		memeduck.Eq(memeduck.Ident("race"), "Phoenix"),
		memeduck.Eq(memeduck.Ident("work_at"), "KFP"),
	).SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at FROM user WHERE race = "Phoenix" AND work_at = "KFP"
}

func ExampleSelect_orderBy() {
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		Where(memeduck.IsNotNull(memeduck.Ident("age"))).
		OrderBy("subscribers", memeduck.ASC).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at FROM user WHERE age IS NOT NULL ORDER BY subscribers ASC
}

func ExampleSelect_limit() {
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		Where(memeduck.Eq(memeduck.Ident("likes"), "alcohol")).
		Limit(10).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at FROM user WHERE likes = "alcohol" LIMIT 10
}

func ExampleSelect_limitOffset() {
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		Where(memeduck.Eq(memeduck.Ident("good_at"), "cooking")).
		LimitOffset(10, 3).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at FROM user WHERE good_at = "cooking" LIMIT 10 OFFSET 3
}

func ExampleSelect_multipleWhere() {
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		Where(memeduck.Eq(memeduck.Ident("job"), "detective")).
		Where(memeduck.Eq(memeduck.Ident("defective"), true)).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at FROM user WHERE job = "detective" AND defective = TRUE
}

func ExampleSelect_queryParameter() {
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		Where(memeduck.Gt(memeduck.Ident("age"), memeduck.Param("age"))).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at FROM user WHERE age > @age
}
