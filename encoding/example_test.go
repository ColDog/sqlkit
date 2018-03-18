package encoding_test

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/coldog/sqlkit/encoding"
)

func Example() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	type user struct {
		ID int `db:"id"`
	}

	_, err = db.Exec("create table users (id int primary key)")
	if err != nil {
		panic(err)
	}

	usr := user{ID: 1}

	cols, vals, err := encoding.Marshal(usr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("cols: %v, vals: %v\n", cols, vals)

	_, err = db.Exec(
		"insert into users ("+strings.Join(cols, ",")+") values "+"(?)", vals...,
	)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select * from users")
	if err != nil {
		panic(err)
	}

	usrs := []user{}
	err = encoding.Unmarshal(&usrs, rows)
	if err != nil {
		panic(err)
	}
	fmt.Printf("query: %+v\n", usrs)

	// Output:
	// cols: [id], vals: [1]
	// query: [{ID:1}]
}
