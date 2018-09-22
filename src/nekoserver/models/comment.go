package models

import (
	"fmt"
	"time"

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

func createComment() error {

	timestamp := time.Now().Unix()

	var co data.Comment

	co.CreatedAt = timestamp

	co.ModifiedAt = timestamp

	// TODO: Validate whether the post exists

	statement := fmt.Sprintf("INSERT INTO comment (id, pid, createdAt, ModifiedAt, author, url, ip, body, prid, mail) VALUES('%s', '%s', '%d', '%d', '%s', '%s', '%s', '%s', '%s', '%s')",
		 co.Id, co.PID, co.CreatedAt, co.ModifiedAt, co.Author, co.Url, co.Ip, co.Body, co.Prid, co.Mail)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err
	}

	_, err = db.Exec(statement)

	return err
}

func FetchComments(id string, start, count int) (error, []data.Comment) {
	statement := fmt.Sprintf("SELECT * FROM comment WHERE pid=%s LIMIT %d OFFSET %d", id, count, start)

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
		if err := rows.Scan(&co.PID, &co.Id, &co.CreatedAt, &co.ModifiedAt, &co.Author, &co.Url, &co.Ip, &co.Body, &co.Prid, &co.Mail); err != nil {
			return err, nil
		}
		comments = append(comments, co)
	}

	return nil, comments
}