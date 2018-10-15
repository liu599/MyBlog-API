package models

import (
	"fmt"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
)

func TokenCheckUser(usr data.User) (error, bool) {
	var uk data.User
	statement := fmt.Sprintf("SELECT * FROM user WHERE name='%s'", usr.Name)
	//fmt.Println(usr.Password, usr.Name)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, false
	}
	err = db.QueryRow(statement).Scan(&uk.UID, &uk.USID, &uk.Name, &uk.Password, &uk.Mail, &uk.CreatedAt, &uk.LoggedAt)
	if err != nil || uk.Name != usr.Name || usr.Password != uk.Password {
		fmt.Println(uk.Password)
		return err, false
	}
	return nil, true
}