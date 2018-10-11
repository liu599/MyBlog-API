package models

import (
	"fmt"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
)

func UserFetch(name string) (error, data.User) {
	var uk data.User
	statement := fmt.Sprintf("SELECT * FROM user WHERE name='%s'", name)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, data.User{}
	}

	err = db.QueryRow(statement).Scan(&uk.UID, &uk.USID, &uk.Name, &uk.Password, &uk.Mail, &uk.CreatedAt, &uk.LoggedAt)
    fmt.Println(uk.Password)

	return nil, uk
}