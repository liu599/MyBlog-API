package models

import (
	"database/sql"
	"fmt"
	"html"
	"html/template"
	"time"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"

	"gopkg.in/mgo.v2/bson"
)

func CreatePost(p data.Post) (error, string) {

	id := bson.NewObjectId().Hex()

	statement := fmt.Sprintf("INSERT INTO post (poid, ptitle, slug, created, modified, author, template, category, password, status, body) VALUES('%s', '%s', '%s', '%d', '%d', '%s', '%s', '%s', '%s')",
		p.Id, p.PTitle, p.Slug, p.CreatedAt, p.ModifiedAt, p.Author, p.Category, p.Password, p.Body)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, ""
	}

	_, err = db.Exec(statement)

	if err != nil {
		return err, ""
	}

	return nil, id


}

func UpdatePost(post data.Post) error {
	timestamp := time.Now().Unix()

	post.ModifiedAt = timestamp

	post.Body = template.HTMLEscapeString(post.Body)

	statement := fmt.Sprintf("UPDATE post SET ptitle='%s', slug='%s', category='%s', author='%s', body='%s', password='%s', createdAt='%d', modifiedAt='%d' WHERE poid=%s",
		post.PTitle, post.Slug, post.Category, post.Author, post.Body, post.Password, post.CreatedAt, post.ModifiedAt, post.Id)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err
	}

	_, err = db.Exec(statement)

	if err != nil {
		return err
	}

	return err
}

func DeletePost(id string) error {
	statement := fmt.Sprintf("DELETE FROM post WHERE poid=%d", id)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err
	}
	_, err = db.Exec(statement)
	return err
}

func PostsFetchChronology() (error, []string) {
	statement := fmt.Sprintf("select DATE_FORMAT(FROM_UNIXTIME(`createdAt`), '%s') from post ORDER BY `createdAt` ASC", "%Y%m")
	fmt.Println(statement)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, nil
	}
	rows, err := db.Query(statement)

	var ret []string
	for rows.Next() {
		var p string
		if err = rows.Scan(&p); err != nil {
			return err, nil
		}
		ret = append(ret, p)
	}
	ret = _func.ArrayFilter(ret)
	if err != nil {
		fmt.Println("Failure to get chronology")
		return err, nil
	}
	return nil, ret
}


func PostsFetchTotalNumber() (error, int) {

	var countNumber int

	statement := fmt.Sprintf("SELECT COUNT(pid) FROM post")

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, -25252
	}

	err = db.QueryRow(statement).Scan(&countNumber)

	if err != nil {
		return err, -25252
	}

	//fmt.Println(countNumber)

	return  nil, countNumber
}

func PostsFetchTotalNumberByCategory(id string) (error, int) {
	var countNumber int

	statement := fmt.Sprintf("SELECT COUNT(pid) FROM post where category='%s'", id)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, -25252
	}

	err = db.QueryRow(statement).Scan(&countNumber)

	//fmt.Println(countNumber)

	return nil, countNumber
}

func PostsFetchAllWithPageNumber(start, count int) (error, []data.Post) {
	statement := fmt.Sprintf("SELECT * FROM post left join category on post.category=category.id ORDER BY `createdAt` DESC LIMIT %d OFFSET %d", count, start)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err, []data.Post{}
	}
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
		return err, []data.Post{}
	}
	posts := []data.Post{}
	for rows.Next() {
		var p data.Post
		var cc data.Category
		var nulString sql.NullString
		if err = rows.Scan(&p.PID, &p.Id, &p.Author, &p.Category, &p.Body, &p.PTitle, &p.Slug, &p.Password, &p.CreatedAt, &p.ModifiedAt, &cc.CID, &cc.Id, &cc.CName, &cc.CLink, &nulString); err != nil {
			return err, nil
		}
		p.Body = html.UnescapeString(p.Body)
		p.Category = cc.CName
		posts = append(posts, p)
	}
	return nil, posts
}

func PostsFetchCategoryWithPageNumber(start, count int, cid string) (error, []data.Post) {
	statement := fmt.Sprintf("SELECT * FROM post left join category on post.category=category.id WHERE category='%s' ORDER BY `createdAt` DESC LIMIT %d OFFSET %d", cid, count, start)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err, []data.Post{}
	}
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
		return err, []data.Post{}
	}
	posts := []data.Post{}
	for rows.Next() {
		var p data.Post
		var cc data.Category
		var nulString sql.NullString
		if err = rows.Scan(&p.PID, &p.Id, &p.Author, &p.Category, &p.Body, &p.PTitle, &p.Slug, &p.Password, &p.CreatedAt, &p.ModifiedAt, &cc.CID, &cc.Id, &cc.CName, &cc.CLink, &nulString); err != nil {
			return err, nil
		}
		p.Body = html.UnescapeString(p.Body)
		p.Category = cc.CName
		posts = append(posts, p)
	}
	return nil, posts
}

func PostFetchOne(id string) (error, data.Post) {
	statement := fmt.Sprintf("SELECT * FROM post left join category on post.category=category.id WHERE poid='%s'", id)
	fmt.Println(statement)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, data.Post{}
	}
	var p data.Post
	var cc data.Category
	var nulString sql.NullString
	err = db.QueryRow(statement).Scan(&p.PID, &p.Id, &p.Author, &p.Category, &p.Body, &p.PTitle, &p.Slug, &p.Password, &p.CreatedAt, &p.ModifiedAt, &cc.CID, &cc.Id, &cc.CName, &cc.CLink, &nulString)
	p.Category = cc.CName
	if err != nil {
		return err, data.Post{}
	}
	return err, p
}

func fetchAllPosts() (error, []data.Post) {
	statement := fmt.Sprintf("SELECT poid, created FROM post ORDER BY `createdAt` DESC")

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, nil
	}

	rows, err := db.Query(statement)
	if err != nil {
		return err, nil
	}

	defer rows.Close()

	posts := []data.Post{}

	for rows.Next() {
		var p data.Post
		if err:= rows.Scan(&p.Id, &p.CreatedAt); err != nil {
			return err, nil
		}
		posts = append(posts, p)
	}
	return nil, posts
}