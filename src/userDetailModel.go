package main

import (
	"database/sql"
	"fmt"
)

type userDetail struct {
	UID int `json:"uid"`
	Groupid int `json:"groupid"`
	Nick string `json:"nick"`
	Url string `json:"url"`
	Avatar int `json:"avatar"`
	Intro string `json:"intro"`
}

func (ud *userDetail) getUserDetail(db *sql.DB) error {
	// fmt.Printf("\nSELECT * FROM users_nh WHERE uid=%d\n", ud.UID)
	statement := fmt.Sprintf("SELECT groupid, nick, url, avatar, intro FROM users_nh WHERE uid=%d", ud.UID)
	//return errors.New("Not implemented")
	return db.QueryRow(statement).Scan(&ud.Groupid, &ud.Nick, &ud.Url, &ud.Avatar, &ud.Intro)
}


func getUserDetails(db *sql.DB, start, count int) ([]userDetail, error) {
	statement := fmt.Sprintf("SELECT uid, groupid, nick, url, avatar, intro FROM users_nh LIMIT %d OFFSET %d",
		count, start)

	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userDetails := []userDetail{}

	for rows.Next() {
		var u userDetail
		if err := rows.Scan(&u.UID, &u.Groupid, &u.Nick, &u.Url, &u.Avatar, &u.Intro); err != nil {
			return nil, err
		}
		userDetails = append(userDetails, u)
	}

	return userDetails, nil
}