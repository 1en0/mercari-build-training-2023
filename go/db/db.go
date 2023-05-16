package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func DbInit() {
	db, err := sql.Open("sqlite3", "./db/mercari.sqlite")
	Db = db
	if err != nil {
		panic("database connection failed")
	}
}

func AddItem(name string, category string, image_filename string) error {
	id, err := getCategoryId(category)
	if err == sql.ErrNoRows {
		id, err = addCategory(category)
		if err != nil {
			return err
		}
	}
	stmt, err := Db.Prepare("INSERT INTO items (name, category_id, image_filename) values (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, id, image_filename)
	if err != nil {
		return err
	}
	return nil
}

func getCategoryId(category string) (int64, error) {
	var id int64
	err := Db.QueryRow("SELECT id FROM category WHERE NAME = ?", category).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func addCategory(category string) (int64, error) {
	stmt, err := Db.Prepare("INSERT INTO category (name) values (?)")
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(category)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}
