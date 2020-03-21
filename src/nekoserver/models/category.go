package models
<<<<<<< HEAD
=======

import (
	"database/sql"
	"fmt"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"

	"gopkg.in/mgo.v2/bson"
)

func FindCategory(p data.Category) (error, bool) {
	statement := fmt.Sprintf("select count(cid) from category where id = '%s'", p.Id)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err, false
	}

	res, err := db.Exec(statement)
	num, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Fail to execute sql query %v", err)
		return err, false
	}

	return nil, num > 0
}

func CategoryCreate(category data.Category) (error, string) {


	id := bson.NewObjectId().Hex()

	statement := fmt.Sprintf("INSERT INTO category (id, cname, clink, cinfo) VALUES ('%s', '%s', '%s', '%s')",
		id, category.CName, category.CLink, category.CInfo)

	fmt.Println(statement)

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

func CategoryUpdate(category data.Category) (error) {
	statement := fmt.Sprintf("UPDATE category SET cname='%s', cinfo='%s' WHERE id='%s'",
		category.CName, category.CInfo,  category.Id)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err
	}
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func FetchCategoryList() (error, []data.Category) {
	statement := fmt.Sprintf("SELECT * FROM `category`")
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
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
		if err := rows.Scan(&cat.CID, &cat.Id, &cat.CName, &cat.CLink, &nulString); err != nil {
			return err, []data.Category{}
		}
		if nulString.Valid {
			cat.CInfo = nulString.String
		} else {
			cat.CInfo = ""
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
		return err, data.Category{}
	}
	err = db.QueryRow(statement).Scan(&category.CID, &category.Id, &category.CName, &category.CInfo, &category.CLink)
	if err != nil {
		return err, data.Category{}
	}
	return nil, category
}

func DeleteCategory(id string) data.Error {
	_, eflag := FindCategory(data.Category{Id: id})
	if eflag != true {
		fmt.Println("Cannot find category")
		return data.Error{
			Code: "401",
			Message: "Fail to find the category, probably wrong id!",
		}
	}
	statement := fmt.Sprintf("DELETE FROM category WHERE id='%s'", id)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return data.Error{
			Code: "502",
			Message: "Error Database Connection",
		}
	}
	_, err = db.Exec(statement)
	if err != nil {
		return data.Error{
			Code: "502",
			Message: "Error Database Connection",
		}
	}

	return data.Error{}
}
>>>>>>> nekohandserverv1/master
