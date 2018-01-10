package main

import (
	"database/sql"
	"fmt"
	"time"
)

type comment struct {
	COID int `json:"coid"`
	PID int `json:"pid"`
	Created int64 `json:"created"`
	Author string `json:"author"`
	Url string `json:"url"`
	Ip string `json:"ip"`
	Body string `json:"body"`
	Parent int `json:"parent"`
	Status string `json:"status"`
	Mail string `json:"mail"`
}

func (co *comment) getComment(db *sql.DB) error {
	fmt.Printf("\nSELECT pid, created, author, url, ip, body, parent, status, mail FROM comments_nh WHERE coid=%d\n", co.COID)
	statement := fmt.Sprintf("SELECT pid, created, author, url, ip, body, parent, status, mail FROM comments_nh WHERE coid=%d", co.COID)
	//return errors.New("Not implemented")
	return db.QueryRow(statement).Scan(&co.PID, &co.Created, &co.Author, &co.Url, &co.Ip, &co.Body, &co.Parent, &co.Status, &co.Mail)
}


func (co *comment) updateComment(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE comments_nh SET created='%d', author='%s', url='%s', ip='%s', body='%s', parent='%d', status='%s', mail='%s' WHERE coid=%d",
		co.Created, co.Author, co.Url, co.Ip, co.Body, co.Parent, co.Status, co.Mail, co.COID)
	_, err := db.Exec(statement)
	return err
}

func (co *comment) deleteComment(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM comments_nh WHERE coid=%d", co.COID)
	_, err := db.Exec(statement)
	return err
}


func (co *comment) createComment(db *sql.DB) error {

	timestamp := time.Now().Unix()

	co.Created = timestamp

	fmt.Printf("\n%d\n", co.COID)

	// TODO: Validate whether the post exists

	statement := fmt.Sprintf("INSERT INTO comments_nh (coid, pid, created, author, url, ip, body, parent, status, mail) VALUES('%d', '%d', '%d', '%s', '%s', '%s', '%s', '%d', '%s', '%s')",
		co.COID, co.PID, co.Created, co.Author, co.Url, co.Ip, co.Body, co.Parent, co.Status, co.Mail)

	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	return nil
}

func getComments(db *sql.DB, pid, start, count int) ([]comment, error) {

	fmt.Printf("\n%d %d %d\n", pid, start, count)

	statement := fmt.Sprintf("SELECT * FROM comments_nh WHERE pid=%d LIMIT %d OFFSET %d", pid, count, start)

	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []comment{}

	for rows.Next() {
		var co comment
		if err := rows.Scan(&co.PID, &co.COID, &co.Created, &co.Author, &co.Url, &co.Ip, &co.Body, &co.Parent, &co.Status, &co.Mail); err != nil {
			return nil, err
		}
		comments = append(comments, co)
	}

	return comments, nil
}