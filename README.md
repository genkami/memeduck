# memeduck
![ci status](https://github.com/genkami/memeduck/workflows/Test/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/genkami/memeduck.svg)](https://pkg.go.dev/github.com/genkami/memeduck)

![duck](./doc/img/memeduck.png)

The memeduck is a SQL query builder for Cloud Spanner. (named after [MakeNowJust/memefish](https://github.com/MakeNowJust/memefish))

## Examples
### Select
```go
	query, _ := memeduck.Select("user", []string{"name", "created_at"}).
		Where(memeduck.Eq(memeduck.Ident("good_at"), "cooking")).
		LimitOffset(10, 3).
		SQL()
	fmt.Println(query)
	// Output: SELECT name, created_at FROM user WHERE good_at = "cooking" LIMIT 10 OFFSET 3
```

### Insert
```go
type ExampleUserStruct struct {
	Name string `spanner:"UserName"`
	Papa string `spanner:"PapaName"`
}
```

```go
	query, _ := memeduck.Insert("users", []string{"UserName", "PapaName"}).Values([]ExampleUserStruct{
		{Name: "Kiara", Papa: "huke"},
	}).SQL()
	fmt.Println(query)
	// Output: INSERT INTO users (UserName, PapaName) VALUES ("Kiara", "huke")
```

### Update
```go
	query, _ := memeduck.Update("user").
		Set(memeduck.Ident("age"), memeduck.Param("age")).
		Where(memeduck.Eq(memeduck.Ident("shark"), true)).
		SQL()
	fmt.Println(query)
	// Output: UPDATE user SET age = @age WHERE shark = TRUE
```

### Delete
```go
	query, _ := memeduck.Delete("user").
		Where(memeduck.Eq(memeduck.Ident("id"), 123)).
		Where(memeduck.Eq(memeduck.Ident("unused"), true)).
		SQL()
	fmt.Println(query)
	// Output: DELETE FROM user WHERE id = 123 AND unused = TRUE
```

You can see more examples in the examples sectionss of the [documentation](https://pkg.go.dev/github.com/genkami/memeduck).

## License

Distributed under the Apache License, Version 2.0. See LICENSE for more information.
