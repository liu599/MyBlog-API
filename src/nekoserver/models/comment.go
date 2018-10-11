package models

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
)

func deleteComment(id string) error {
	statement := fmt.Sprintf("DELETE FROM comment WHERE id=%d", id)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err
	}

	_, err = db.Exec(statement)

	return err
}

func CommentCreate(co data.Comment) error {

	timestamp := time.Now().Unix()

	co.COMID = bson.NewObjectId().Hex()

	co.CreatedAt = timestamp

	co.ModifiedAt = timestamp

	statement := fmt.Sprintf("INSERT INTO comment ( commentid, pid, author, mail, url, ip, prid, body, createdAt, ModifiedAt) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d')",
		 co.COMID, co.PID, co.Author, co.Mail, co.Url, co.Ip, co.Prid, co.Body, co.CreatedAt, co.ModifiedAt)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err
	}

	_, err = db.Exec(statement)

	return err
}

func CommentsFetch(id string) (error, []data.Comment) {
	statement := fmt.Sprintf("SELECT * FROM comment WHERE pid='%s'", id)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, nil
	}

	rows, err := db.Query(statement)

	defer rows.Close()

	comments := []data.Comment{}

	for rows.Next() {
		var co data.Comment
		if err := rows.Scan(&co.COID, &co.PID, &co.COMID, &co.Author, &co.Mail, &co.Url, &co.Ip, &co.Prid, &co.Body, &co.CreatedAt, &co.ModifiedAt); err != nil {
			return err, nil
		}
		comments = append(comments, co)
	}

	return nil, comments
}

func CommentsFetchNumber(id string) (error, int) {
	var countNumber int
	statement := fmt.Sprintf("SELECT COUNT(commentid) FROM comment WHERE pid='%s'", id)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, -25252
	}

	err = db.QueryRow(statement).Scan(&countNumber)

	return nil, countNumber
}