package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func GetOneUserNameBy(age int64) (string, error) {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", errors.Wrap(err, "open db err")
	}

	var name string
	sqlTpl := "SELECT name FROM users WHERE age = $1"
	err = db.QueryRow(sqlTpl, age).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} else {
			return "", errors.Wrap(err, fmt.Sprintf("err by sql tpl \"%s\", args: %v", sqlTpl, age))
		}
	}
	return name, nil
}


func main() {
	var age int64 = 21
	name, err := GetOneUserNameBy(age)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if name == "" {
		log.Println("Not found any user age is", age)
		return
	}

}
