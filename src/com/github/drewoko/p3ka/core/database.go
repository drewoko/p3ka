package core

import (
	"log"
	"time"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const MESSAGE_MAIN_QUERY string = "SELECT id, name, origurl as url, mature FROM messages";

type RowMap map[string]interface{}

type DataBase struct {
	db *sql.DB
}

func (self *DataBase) Init(path string) *DataBase {

	var err error
	self.db, err = sql.Open("sqlite3", path)

	if(err != nil) {
		log.Fatal("Failed to open database. Reason: ", err)
	}

	self.createTable()

	return self
}

func (self *DataBase) createTable() {
	self.db.Exec("CREATE TABLE messages (id integer not null primary key, date int, messageid int, name text, origurl text, deleted int DEFAULT 0)")
	self.db.Exec("ALTER TABLE messages ADD COLUMN mature int DEFAULT 0;")
}

func (self *DataBase) IsExists(author string, url string) bool {
	s := self.db.QueryRow("SELECT COUNT(id) as cnt FROM messages WHERE origurl=?", url)
	var cnt int
	s.Scan(&cnt)
	return cnt > 0
}

func (self *DataBase) AddRow(messageid int64, author string, url string, mature bool) {

	_, err := self.db.Exec(
		"INSERT INTO messages (messageid, date, name, origurl, mature) VALUES (?, ?, ?, ?, ?)",
			int(messageid), int(time.Now().Unix()), author, url, mature)

	if(err != nil) {
		self.processError(err)
	}
}

func (self *DataBase) GetMessageById(id int64) RowMap {
	s := self.db.QueryRow(MESSAGE_MAIN_QUERY + " where deleted = 0 and id=?", int(id));
	return self.SingleMessageScan(s)
}

func (self *DataBase) GetMessageByMessageId(id int64) RowMap {
	s := self.db.QueryRow(MESSAGE_MAIN_QUERY + " where deleted = 0 and messageid=?", int(id))
	return self.SingleMessageScan(s)
}

func (self *DataBase) SetDeleted(id interface{}) {
	self.db.Exec("UPDATE messages set deleted = 1 where id=?", id)
}

func (self *DataBase) GetAll() []RowMap {
	s, err := self.db.Query(MESSAGE_MAIN_QUERY + " where deleted=0");

	if(err != nil) {
		self.processError(err)
	}

	return self.MultipleMessageScan(s)
}

func (self *DataBase) GetLast(limit int, start int) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " where deleted=0 order by id desc limit ?,?", start, limit)

	if(err != nil) {
		self.processError(err)
	}

	return self.MultipleMessageScan(s)
}

func (self *DataBase) get_last_user(limit int, start int, user string) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " where deleted=0 and name=? order by id desc limit ?,?", user, start, limit)

	if(err != nil) {
		self.processError(err)
	}

	return self.MultipleMessageScan(s)
}

func (self *DataBase) get_random(limit int) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " WHERE deleted = 0 ORDER BY RANDOM() LIMIT ?", limit)

	if(err != nil) {
		self.processError(err)
	}

	return self.MultipleMessageScan(s)
}

func (self *DataBase) get_top_users(limit int, exclude []string) []RowMap {

	var rows []RowMap

	s, err := self.db.Query(
		"select count(*) as cnt, name from messages where deleted = 0 group by name order by cnt desc LIMIT ?", limit)

	if(err != nil) {
		self.processError(err)
	}

	for s.Next() {
		row := make(RowMap)

		var cnt int
		var name string

		s.Scan(&cnt, &name)

		row["cnt"] = cnt
		row["name"] = name

		if(!ContainsString(exclude, row["name"].(string))) {
			rows = append(rows, row)
		}
	}

	return rows
}

func (self *DataBase) MultipleMessageScan(s *sql.Rows) []RowMap {
	var rows []RowMap

	for s.Next() {
		rows = append(rows, self.MultipleItemMessageScan(s))
	}

	s.Close()

	return rows
}

func (self *DataBase) MultipleItemMessageScan(s *sql.Rows) RowMap {
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

func (self *DataBase) SingleMessageScan(s *sql.Row) RowMap {
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

func (self *DataBase) processError(err error) {
	log.Println("Failed to execute query. Reason: ", err)
}

func (self *DataBase) Close() {
	self.db.Close();
}