package main

import (
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const MESSAGE_MAIN_QUERY string = "SELECT id, name, origurl as url, mature FROM messages";

type RowMap map[string]interface{}

type DataBase struct {
	db *sql.DB
}

func (self *DataBase) init(path string) *DataBase {

	var err error
	self.db, err = sql.Open("sqlite3", path)

	if(err != nil) {
		log.Fatal("Failed to open database. Reason: ", err)
	}

	self.create_table()

	return self
}

func (self *DataBase) create_table() {

	self.db.Query("CREATE TABLE messages (id integer not null primary key, date int, messageid int, name text, origurl text, deleted int DEFAULT 0)")
	self.db.Query("ALTER TABLE messages ADD COLUMN mature int DEFAULT 0;")
}

func (self *DataBase) is_exists(author string, url string) bool {
	s := self.db.QueryRow("SELECT COUNT(id) as cnt FROM messages WHERE origurl=?", url)
	var cnt int
	s.Scan(&cnt)
	return cnt > 0
}

func (self *DataBase) add_row(messageid int64, author string, url string, mature bool) {

	_, err := self.db.Exec(
		"INSERT INTO messages (messageid, date, name, origurl, mature) VALUES (?, ?, ?, ?, ?)",
			int(messageid), int(time.Now().Unix()), author, url, mature)

	if(err != nil) {
		self.process_error(err)
	}
}

func (self *DataBase) get_message_by_id(id int64) RowMap {
	s := self.db.QueryRow(MESSAGE_MAIN_QUERY + " where deleted = 0 and id=?", int(id));
	return self.single_message_scan(s)
}

func (self *DataBase) get_message_by_messageid(id int64) RowMap {
	s := self.db.QueryRow(MESSAGE_MAIN_QUERY + " where deleted = 0 and messageid=?", int(id))
	return self.single_message_scan(s)
}

func (self *DataBase) set_deleted(id interface{}) {
	self.db.Query("UPDATE messages set deleted = 1 where id=?", id)
}

func (self *DataBase) get_all() []RowMap {
	s, err := self.db.Query(MESSAGE_MAIN_QUERY + " where deleted=0");

	if(err != nil) {
		self.process_error(err)
	}

	return self.multiple_message_scan(s)
}

func (self *DataBase) get_last(limit int, start int) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " where deleted=0 order by id desc limit ?,?", start, limit)

	if(err != nil) {
		self.process_error(err)
	}

	return self.multiple_message_scan(s)
}

func (self *DataBase) get_last_user(limit int, start int, user string) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " where deleted=0 and name=? order by id desc limit ?,?", user, start, limit)

	if(err != nil) {
		self.process_error(err)
	}

	return self.multiple_message_scan(s)
}

func (self *DataBase) get_random(limit int) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " WHERE deleted = 0 ORDER BY RANDOM() LIMIT ?", limit)

	if(err != nil) {
		self.process_error(err)
	}

	return self.multiple_message_scan(s)
}

func (self *DataBase) get_top_users(limit int, exclude []string) []RowMap {

	var rows []RowMap

	s, err := self.db.Query(
		"select count(*) as cnt, name from messages where deleted = 0 group by name order by cnt desc LIMIT ?", limit)

	if(err != nil) {
		self.process_error(err)
	}

	for s.Next() {
		row := make(RowMap)

		var cnt int
		var name string

		s.Scan(&cnt, &name)

		row["cnt"] = cnt
		row["name"] = name

		if(!contains_string(exclude, row["name"].(string))) {
			rows = append(rows, row)
		}
	}

	return rows
}

func (self *DataBase) multiple_message_scan(s *sql.Rows) []RowMap {
	var rows []RowMap

	for s.Next() {
		rows = append(rows, self.multiple_item_message_scan(s))
	}

	return rows
}

func (self *DataBase) multiple_item_message_scan(s *sql.Rows) RowMap {
	row := make(RowMap)

	var id int
	var name string
	var url string
	var mature int

	s.Scan(&id, &name, &url, &mature)

	row["id"] = id
	row["name"] = name
	row["url"] = url
	row["mature"] = mature

	return row
}

func (self *DataBase) single_message_scan(s *sql.Row) RowMap {
	row := make(RowMap)

	var id int
	var name string
	var url string
	var mature int

	s.Scan(&id, &name, &url, &mature)

	row["id"] = id
	row["name"] = name
	row["url"] = url
	row["mature"] = mature

	return row
}

func (self *DataBase) process_error(err error) {
	log.Info("Failed to execute query. Reason: ", err)
}
