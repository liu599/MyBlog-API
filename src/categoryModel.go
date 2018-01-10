package main

import (
	"database/sql"
	"fmt"
)

type category struct {
	CID int `json:"cid"`
	CName string `json:"cname"`
	CInfo string `json:"cinfo"`
}

func (c *category) getCategoryByName(db *sql.DB) error {
	fmt.Printf("\nSELECT cid, cinfo FROM post_categories_nh WHERE cname='%s'\n", c.CName)
	statement := fmt.Sprintf("SELECT cid, cinfo FROM post_categories_nh WHERE cname='%s'", c.CName)
	//return errors.New("Not implemented")
	return db.QueryRow(statement).Scan(&c.CID, &c.CInfo)
}

func (c *category) getCategory(db *sql.DB) error {
	fmt.Printf("\nSELECT cname, cinfo FROM post_categories_nh WHERE cid=%d\n", c.CID)
	statement := fmt.Sprintf("SELECT cname, cinfo FROM post_categories_nh WHERE cid=%d", c.CID)
	//return errors.New("Not implemented")
	return db.QueryRow(statement).Scan(&c.CName, &c.CInfo)
}

func (c *category) updateCategory(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE post_categories_nh SET cname='%s', cinfo='%s' WHERE cid=%d",
		c.CName, c.CInfo, c.CID)
	_, err := db.Exec(statement)
	return err
}

func (c *category) deleteCategory(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM post_categories_nh WHERE cid=%d", c.CID)
	_, err := db.Exec(statement)
	return err
}


func (c *category) createCategory(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO post_categories_nh (cid, cname, cinfo) VALUES('%d', '%s', '%s')",
		c.CID, c.CName, c.CInfo)

	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&c.CID)

	if err != nil {
		return err
	}

	return nil
}


func getCategories(db *sql.DB, start, count int) ([]category, error) {
	statement := fmt.Sprintf("SELECT cid, cname, cinfo FROM post_categories_nh LIMIT %d OFFSET %d",
		count, start)

	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []category{}

	for rows.Next() {
		var c category
		if err := rows.Scan(&c.CID, &c.CName, &c.CInfo); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}