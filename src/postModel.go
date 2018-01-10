package main

import (
	"database/sql"
	"fmt"
	"time"
	"html/template"
	"html"
)

type post struct {
	PID int `json:"pid"`
	PTitle string `json:"title"`
	Slug string `json:"slug"`
	Category int `json:"category"`
	Template int `json:"template"`
	Status string `json:"status"`
	Author int `json:"author"`
	Body string `json:"body"`
	Password string `json:"password"`
	Created int64 `json:"created"`
	Modified int64 `json:"modified"`
}

type pager struct {
	Page int `json:"page"`
	Count int `json:"count"`
	TotalNumber int `json:"total"`
}


func (p *post) getPost(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT ptitle, slug, created, modified, author, template, category, password, status, body FROM posts_nh WHERE pid=%d", p.PID)
	//return errors.New("Not implemented")
	return db.QueryRow(statement).Scan(&p.PTitle, &p.Slug, &p.Created, &p.Modified, &p.Author, &p.Template, &p.Category, &p.Password, &p.Status, &p.Body)
}

func (p *post) updatePost(db *sql.DB) error {

	timestamp := time.Now().Unix()

	p.Modified = timestamp

	p.Body = template.HTMLEscapeString(p.Body)

	statement := fmt.Sprintf("UPDATE posts_nh SET ptitle='%s', slug='%s', category='%d', template='%d', status='%s', author='%d', body='%s', password='%s', created='%d', modified='%d' WHERE pid=%d",
		p.PTitle, p.Slug, p.Category, p.Template, p.Status, p.Author, p.Body, p.Password, p.Created, p.Modified, p.PID)

	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	return err
}

func (p *post) deletePost(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM posts_nh WHERE pid=%d", p.PID)
	_, err := db.Exec(statement)
	return err
}

func (p *post) createPost(db *sql.DB) error {

	timestamp := time.Now().Unix()

	p.Created = timestamp

	p.Modified = timestamp

	p.Body = template.HTMLEscapeString(p.Body)

	statement := fmt.Sprintf("INSERT INTO posts_nh (pid, ptitle, slug, created, modified, author, template, category, password, status, body) VALUES('%d', '%s', '%s', '%d', '%d', '%d', '%d', '%d', '%s', '%s', '%s')",
		p.PID, p.PTitle, p.Slug, p.Created, p.Modified, p.Author, p.Template, p.Category, p.Password, p.Status, p.Body)

	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	return err
}

func getPostsCountNumber (db *sql.DB) (int, error) {

	var countNumber int

	statement := fmt.Sprintf("SELECT COUNT(pid) FROM posts_nh")

	db.QueryRow(statement).Scan(&countNumber)

	fmt.Println(countNumber)

	return countNumber, nil
}

func getPostsCreateTime (db *sql.DB) ([]post, error) {
	statement := fmt.Sprintf("SELECT pid, created FROM posts_nh ORDER BY `created` DESC")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	postCreated := []post{}

	for rows.Next() {
		var p post
		if err:= rows.Scan(&p.PID, &p.Created); err != nil {
			return nil, err
		}
		postCreated = append(postCreated, p)
	}
	return postCreated, nil
}

func getPostsCountNumberByCategory (db *sql.DB, cid int) (int, error) {

	var countNumber int

	statement := fmt.Sprintf("SELECT COUNT(pid) FROM posts_nh where category=%d", cid)

	db.QueryRow(statement).Scan(&countNumber)

	fmt.Println(countNumber)

	return countNumber, nil
}

func getPosts(db *sql.DB, start, count int) ([]post, error) {
	statement := fmt.Sprintf("SELECT * FROM posts_nh ORDER BY `created` DESC LIMIT %d OFFSET %d", count, start)

	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []post{}

	for rows.Next() {
		var p post
		if err:= rows.Scan(&p.PID, &p.PTitle, &p.Slug, &p.Created, &p.Modified, &p.Author, &p.Template, &p.Category, &p.Password, &p.Status, &p.Body); err != nil {
			return nil, err
		}
		p.Body = html.UnescapeString(p.Body)
		posts = append(posts, p)
	}

	return posts, nil
}

func getPostsByCategory (db *sql.DB, category, start, count int) ([]post, error) {
	statement := fmt.Sprintf("SELECT pid, ptitle, slug, created, modified, author, template, password, status, body FROM posts_nh WHERE category=%d ORDER BY `created` DESC LIMIT %d OFFSET %d", category, count, start)

	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []post{}

	for rows.Next() {
		var p post
		if err:= rows.Scan(&p.PID, &p.PTitle, &p.Slug, &p.Created, &p.Modified, &p.Author, &p.Template, &p.Password, &p.Status, &p.Body); err != nil {
			return nil, err
		}
		p.Category = category
		p.Body = html.UnescapeString(p.Body)
		posts = append(posts, p)
	}

	return posts, nil
}