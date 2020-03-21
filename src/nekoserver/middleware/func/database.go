package _func

import (
	"fmt"
	"log"

<<<<<<< HEAD
	"nekoserver/middleware/data"
	"v3/common"
=======

	"nekoserver/middleware/data"
>>>>>>> nekohandserverv1/master

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
<<<<<<< HEAD
	fmt.Println(AppDatabase)
=======
>>>>>>> nekohandserverv1/master

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
<<<<<<< HEAD
		fmt.Println(database.Source)
		fmt.Println(err)
=======
>>>>>>> nekohandserverv1/master
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

// GetDBconnection

func MySqlGetDB(connection string) (*sqlx.DB, error) {
	var err error
	conn, exists := AppDatabase[connection]
	if !exists || conn.DB.(*sqlx.DB).Ping() != nil {
		AssignDatabaseFromList([]string{"nekohand"})
		conn, exists = AppDatabase[connection]
		if !exists {
<<<<<<< HEAD
			err = fmt.Errorf(common.MSG_DATABASE_CONNECTION_NOT_EXISTS)
=======
			err = fmt.Errorf("%s", "Database Connection Not Exists")
			//err = fmt.Errorf(common.MSG_DATABASE_CONNECTION_NOT_EXISTS)
>>>>>>> nekohandserverv1/master
			return nil, err
		}
	}
	// return sqlx.NewDb(conn.DB.(*sqlx.DB), conn.Driver), err
	return conn.DB.(*sqlx.DB), err
}

