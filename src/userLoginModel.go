package main

import (
	"database/sql"
	"fmt"
	"time"
)

type user struct {
	UID int `json:"uid"`
	Name string `json:"name"`
	Password string `json:"password"`
	Mail string `json:"mail"`
	Created int64 `json:"created"`
	Logged int64 `json:"logged"`
}

func (u *user) getUser(db *sql.DB) error {
	fmt.Printf("\nSELECT name, mail, created, logged FROM user_login_nh WHERE uid=%d\n", u.UID)
	statement := fmt.Sprintf("SELECT name, mail, created, logged FROM user_login_nh WHERE uid=%d", u.UID)
	//return errors.New("Not implemented")
	u.Password = "************"
	return db.QueryRow(statement).Scan(&u.Name, &u.Mail, &u.Created, &u.Logged)
}


func (ui *userinfo) updateUser(db *sql.DB) error {
	var u user = ui.UserLogin
	var ud userDetail = ui.UserDetail
	statement := fmt.Sprintf("UPDATE user_login_nh SET name='%s', mail='%s' WHERE uid=%d",
		u.Name, u.Mail, u.UID)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	statement = fmt.Sprintf("UPDATE users_nh SET nick='%s', url='%s', avatar='%d', intro='%s' WHERE uid=%d",
		ud.Nick, ud.Url, ud.Avatar, ud.Intro, ud.UID)
	_, err = db.Exec(statement)
	return err
}

func (u *user) logginUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE user_login_nh SET logged=%d WHERE uid=%d",
		u.Logged, u.UID)
	_, err := db.Exec(statement)
	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM user_login_nh WHERE uid=%d", u.UID)
	_, err := db.Exec(statement)
	return err
}


func (u *user) createUser(db *sql.DB) error {

	timestamp := time.Now().Unix()

	u.Created = timestamp

	u.Logged = timestamp

	statement := fmt.Sprintf("INSERT INTO user_login_nh (name, password, mail, created, logged) VALUES('%s', '%s', '%s', '%d', '%d')",
		u.Name, u.Password, u.Mail, u.Created, u.Logged)

	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.UID)

	if err != nil {
		return err
	}
	// key reference table users_nh
	statement_users_info := fmt.Sprintf("INSERT INTO users_nh (uid, groupid, nick, url, avatar, intro) VALUES('%d', '%d', '%s', '%s', '%d', '%s')",
		u.UID, 25252, "Default", "http://bang-dream.com", 1, "This is my info")

	_, err1 := db.Exec(statement_users_info)

	if err1 != nil {
		return err1
	}

	tm := time.Unix(timestamp, 0).Format("2017-04-01 03:04:05 PM")

	// key reference table user_login_auth_nh
	statement_users_login_authorization := fmt.Sprintf("INSERT INTO user_login_auth_nh (uid, auth_code, authorized) VALUES('%d', '%s', '%d')",
		u.UID, MD5(tm), 0)

	_, err2 := db.Exec(statement_users_login_authorization)

	if err2 != nil {
		return err2
	}

	return nil
}


func getUsers(db *sql.DB, start, count int) ([]user, error) {
	statement := fmt.Sprintf("SELECT uid, name, mail, created, logged FROM user_login_nh LIMIT %d OFFSET %d",
		count, start)

	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.UID, &u.Name, &u.Mail, &u.Created, &u.Logged); err != nil {
			return nil, err
		}
		u.Password = "************"
		users = append(users, u)
	}

	return users, nil
}

