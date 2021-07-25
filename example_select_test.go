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

func ExampleSelect_queryScalarSubquery() {
	subQueryStmt := memeduck.Select("user_status", []string{"state"}).
		Where(memeduck.Eq(memeduck.Ident("user_id"), memeduck.Param("id")))
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		SubQuery(memeduck.ScalarSubQuery(subQueryStmt).As("state")).
		Where(memeduck.Eq(memeduck.Ident("id"), memeduck.Param("id"))).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at, (SELECT state FROM user_status WHERE user_id = @id) AS state FROM user WHERE id = @id
}

func ExampleSelect_queryArraySubquery() {
	subQueryStmt := memeduck.Select("user_item", []string{"item_id", "count"}).
		Where(memeduck.Eq(memeduck.Ident("user_id"), memeduck.Param("id"))).
		AsStruct()
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		SubQuery(memeduck.ArraySubQuery(subQueryStmt).As("user_item")).
		Where(memeduck.Eq(memeduck.Ident("id"), memeduck.Param("id"))).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at, ARRAY(SELECT AS STRUCT item_id, count FROM user_item WHERE user_id = @id) AS user_item FROM user WHERE id = @id
}

func ExampleSelect_queryMultiSubQuery() {
	itemStmt := memeduck.Select("user_item", []string{"item_id", "count"}).
		Where(memeduck.Eq(memeduck.Ident("user_id"), "user-id")).
		AsStruct()
	// user has one status
	statusStmt := memeduck.Select("user_status", []string{"state"}).
		Where(memeduck.Eq(memeduck.Ident("user_id"), "user-id")).
		AsStruct()
	query, _ := memeduck.Select("user", []string{"name"}).
		SubQuery(
			memeduck.ArraySubQuery(itemStmt).As("user_item"),
			memeduck.ArraySubQuery(statusStmt).As("user_status"),
		).
		Where(memeduck.Eq(memeduck.Ident("user_id"), "user-id")).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, ARRAY(SELECT AS STRUCT item_id, count FROM user_item WHERE user_id = "user-id") AS user_item, ARRAY(SELECT AS STRUCT state FROM user_status WHERE user_id = "user-id") AS user_status FROM user WHERE user_id = "user-id"
}
