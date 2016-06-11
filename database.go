package main

import (
	"time"

	"github.com/mxk/go-sqlite/sqlite3"
)

type DataBase struct {
	db *sqlite3.Conn
}

func (self *DataBase) init(path string) *DataBase {

	var err error

	self.db, err = sqlite3.Open(path)

	if(err != nil) {
		log.Fatal("Failed to open database. Reason: ", err)
	}

	self.create_table()

	return self
}

func (self *DataBase) create_table() {

	self.db.Exec("CREATE TABLE messages (id integer not null primary key, date int, messageid int, name text, origurl text, deleted int DEFAULT 0)")
	err := self.db.Exec("ALTER TABLE messages ADD COLUMN mature int DEFAULT 0;")

	log.Info("", err)
}

func (self *DataBase) is_exists(author string, url string) bool {

	for s, err := self.db.Query("SELECT * FROM messages WHERE origurl=?", url); err == nil; err = s.Next() {
		return true
	}

	return false
}

func (self *DataBase) add_row(messageid int64, author string, url string, mature bool) {

	self.db.Begin()
	err := self.db.Exec(
		"INSERT INTO messages (messageid, date, name, origurl, mature) VALUES (?, ?, ?, ?, ?)",
			int(messageid), int(time.Now().Unix()), author, url, mature)

	self.db.Commit()

	if(err != nil) {
		self.process_error(err)
	}
}

func (self *DataBase) get_message_by_id(id int64) sqlite3.RowMap {

	row := make(sqlite3.RowMap)
	for s, err := self.db.Query("SELECT * FROM messages where deleted = 0 and id=?", int(id)); err == nil; err = s.Next() {
		s.Scan(row)
	}

	return row
}

func (self *DataBase) get_message_by_messageid(id int64) sqlite3.RowMap {

	row := make(sqlite3.RowMap)
	for s, err := self.db.Query("SELECT * FROM messages where deleted = 0 and messageid=?", int(id)); err == nil; err = s.Next() {
		s.Scan(row)
	}

	return row
}

func (self *DataBase) set_deleted(id interface{}) {
	self.db.Exec("UPDATE messages set deleted = 1 where id=?", id)
}

func (self *DataBase) process_error(err error) {
	log.Info("Failed to execute query. Reason: ", err)
}

func (self *DataBase) get_all_outdated() []sqlite3.RowMap {

	last_time := int(time.Now().Unix())-15552000

	var rows []sqlite3.RowMap
	row := make(sqlite3.RowMap)

	for s, err := self.db.Query("SELECT * FROM messages where deleted=0 and date < ?", last_time); err == nil; err = s.Next() {
		s.Scan(row)

		rows = append(rows, row)
	}

	return rows
}


func (self *DataBase) get_all() []sqlite3.RowMap {

	var rows []sqlite3.RowMap

	for s, err := self.db.Query("SELECT * FROM messages where deleted=0"); err == nil; err = s.Next() {
		row := make(sqlite3.RowMap)
		s.Scan(row)
		rows = append(rows, row)
	}

	return rows
}

func (self *DataBase) get_last(limit int, start int) []sqlite3.RowMap {

	var rows []sqlite3.RowMap

	for s, err := self.db.Query(
		"SELECT name, origurl as url, mature FROM messages where deleted=0 order by id desc limit ?,?", start, limit); err == nil; err = s.Next() {
		row := make(sqlite3.RowMap)

		s.Scan(row)

		rows = append(rows, row)
	}

	return rows
}


func (self *DataBase) get_last_user(limit int, start int, user string) []sqlite3.RowMap {

	var rows []sqlite3.RowMap

	log.Info("GetAllUsers req. limit: %s, start: %s, user: %s", limit, start, user)

	for s, err := self.db.Query(
		"SELECT name, origurl as url, mature FROM messages where deleted=0 and name=? order by id desc limit ?,?", user, start, limit); err == nil; err = s.Next() {
		row := make(sqlite3.RowMap)

		s.Scan(row)

		rows = append(rows, row)
	}

	return rows
}

func (self *DataBase) get_random(limit int) []sqlite3.RowMap {

	var rows []sqlite3.RowMap

	for s, err := self.db.Query(
		"SELECT name, origurl as url, mature FROM messages WHERE deleted = 0 ORDER BY RANDOM() LIMIT ?", limit); err == nil; err = s.Next() {
		row := make(sqlite3.RowMap)

		s.Scan(row)

		rows = append(rows, row)
	}

	return rows
}

func (self *DataBase) get_top_users(limit int, exclude []string) []sqlite3.RowMap {

	var rows []sqlite3.RowMap

	for s, err := self.db.Query(
		"select count(*) as cnt, name from messages where deleted = 0 group by name order by cnt desc LIMIT ?", limit); err == nil; err = s.Next() {
		row := make(sqlite3.RowMap)
		s.Scan(row)

		if(!contains_string(exclude, row["name"].(string))) {
			rows = append(rows, row)
		}
	}

	return rows
}