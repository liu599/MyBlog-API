package models

import (
	"bytes"
	"fmt"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
)

func CheckUser(usr data.User) (error, bool) {
	var uk data.User
	statement := fmt.Sprintf("SELECT * FROM user WHERE name='%s'", usr.Name)
	//fmt.Println(usr.Password, usr.Name)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, false
	}
	err = db.QueryRow(statement).Scan(&uk.UID, &uk.USID, &uk.Name, &uk.Password, &uk.Mail, &uk.CreatedAt, &uk.LoggedAt)
	if err != nil || uk.Name != usr.Name || !bytes.Equal(usr.Password, uk.Password) {
		//fmt.Println("Cannot find user")
		return err, false
	}
	return nil, true
}