package _func

import (
	"fmt"
	"log"

	"nekoserver/middleware/data"

	"github.com/jmoiron/sqlx"
)

var AppDatabase = make(map[string]data.Database)

func AssignAppDataBaseList(adbl map[string]data.Database) {
	AppDatabase = adbl
}

func AssignDatabaseFromList(nameList []string) {
	if AppDatabase == nil {
		AppDatabase = map[string]data.Database{}
	}
	fmt.Println(AppDatabase)

	for k, v := range AppDatabase {
		for _, w := range nameList {
			if k == w {
				AssignMySQL(k, v)
			}
		}
	}
}

// Assign specific mysql
func AssignMySQL(name string, database data.Database) {
	var err error
	mysql, err := sqlx.Connect(database.Driver, database.Source)

	if err != nil {
		fmt.Println(database.Source)
		fmt.Println(err)
		log.Fatal(err)
		panic(err)
	}

	mysql.SetMaxOpenConns(database.MaxOpen)
	mysql.SetMaxIdleConns(database.MaxIdle)

	unsetDB(name)
	AppDatabase[name] = data.Database{
		DB:      mysql,
		Driver:  database.Driver,
		Source:  database.Source,
		MaxOpen: database.MaxOpen,
		MaxIdle: database.MaxIdle,
		Name:    database.Name,
	}
}

// unsetDB
func unsetDB(name string) {
	if _, exists := AppDatabase[name]; exists && AppDatabase[name].DB != nil {
		if AppDatabase[name].Driver == "mysql" {
			AppDatabase[name].DB.(*sqlx.DB).Close()
		}
	}
}