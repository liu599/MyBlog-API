package models

import (
	"database/sql"
	"fmt"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"

	"gopkg.in/mgo.v2/bson"
)

func CreateCategory(category data.Category) (error, string) {


	id := bson.NewObjectId().Hex()

	statement := fmt.Sprintf("INSERT INTO category (id, cname, cinfo) VALUES ('%s', '%s', '%s')",
		id, category.CID, category.CInfo)

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

func UpdateCategory(category data.Category) (error, string) {
	statement := fmt.Sprintf("UPDATE category SET cname='%s', cinfo='%s' WHERE id='%s'",
		category.CName, category.CInfo, category.Id)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, ""
	}
	_, err = db.Exec(statement)
	if err != nil {
		return err, category.Id
	}
	return nil, category.Id
}

func FetchCategoryList() (error, []data.Category) {
	statement := fmt.Sprintf("SELECT * FROM `category`")
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, []data.Category{}
	}
	var categoryList []data.Category
	rows, err := db.Query(statement)
	if err != nil {
		return err, []data.Category{}
	}
	for rows.Next() {
		var cat data.Category
		var nulString sql.NullString
		if err := rows.Scan(&cat.CID, &cat.Id, &cat.CName, &nulString); err != nil {
			return err, []data.Category{}
		}
		if nulString.Valid {
			cat.CInfo = ""
		} else {
			cat.CInfo = nulString.String
		}
		categoryList = append(categoryList, cat)
	}
	return nil, categoryList
}

func FetchOneCategory(id string) (error, data.Category) {
	var category data.Category
	statement := fmt.Sprintf("SELECT * FROM `category` WHERE id='%s'", id)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err, data.Category{}
	}
	err = db.QueryRow(statement).Scan(&category.CID, &category.Id, &category.CName, &category.CInfo)
	if err != nil {
		return err, data.Category{}
	}
	return nil, category
}

func DeleteCategory(id string) error {
	// 删除一级分类
	statement := fmt.Sprintf("DELETE FROM category WHERE id='%s'", id)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return err
	}
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}