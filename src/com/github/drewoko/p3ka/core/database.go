package core

import (
	"log"
	"time"
	"database/sql"

	_ "gopkg.in/mattn/go-sqlite3.v1"
)

const MESSAGE_MAIN_QUERY string = "SELECT id, name, origurl as url, mature, source FROM messages";

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
	self.db.Exec("ALTER TABLE messages ADD COLUMN source text DEFAULT 'peka2tv';")
	self.db.Exec("CREATE TABLE ggChannels (id integer not null primary key, channel text not null)")
}

func (self *DataBase) IsExists(author string, url string) bool {
	s := self.db.QueryRow("SELECT COUNT(id) as cnt FROM messages WHERE origurl=?", url)
	var cnt int
	s.Scan(&cnt)
	return cnt > 0
}

func (self *DataBase) AddRow(messageid int64, author string, url string, mature bool, source string, deleted bool) {

	_, err := self.db.Exec(
		"INSERT INTO messages (messageid, date, name, origurl, mature, source, deleted) VALUES (?, ?, ?, ?, ?, ?, ?)",
			int(messageid), int(time.Now().Unix()), author, url, mature, source, deleted)

	if(err != nil) {
		self.processError(err)
	}
}

func (self *DataBase) AddRowGGChannel(channel string) {

	s := self.db.QueryRow("SELECT count(id) as cnt FROM ggChannels WHERE channel = ?", channel);
	var cnt int
	s.Scan(&cnt)

	if cnt == 0 {
		_, err := self.db.Exec(
			"INSERT INTO ggChannels (channel) VALUES (?)", channel)

		if(err != nil) {
			self.processError(err)
		}
	}
}

func (self *DataBase) GetGGChannels() []string {
	s, err := self.db.Query("SELECT channel FROM ggChannels")

	if(err != nil) {
		self.processError(err)
	}
	defer s.Close()
	var channels []string
	for s.Next() {
		var channel string
		err = s.Scan(&channel)

		if err != nil {
			log.Println("Failed to parse channel")
			continue
		}

		channels = append(channels, channel)
	}

	return channels
}

func (self *DataBase) GetMessageById(id int64) RowMap {
	s := self.db.QueryRow(MESSAGE_MAIN_QUERY + " where deleted = 0 and id=?", int(id));
	return self.SingleMessageScan(s)
}

func (self *DataBase) GetMessageByIdAndSource(id int64, source string) RowMap {
	s := self.db.QueryRow(MESSAGE_MAIN_QUERY + " where deleted = 0 and id=? and source=?", int(id), source);
	return self.SingleMessageScan(s)
}

func (self *DataBase) GetMessageByMessageId(id int64) RowMap {
	s := self.db.QueryRow(MESSAGE_MAIN_QUERY + " where deleted = 0 and messageid=?", int(id))
	return self.SingleMessageScan(s)
}

func (self *DataBase) SetDeleted(id interface{}) {
	self.db.Exec("UPDATE messages set deleted = 1 where id=?", id)
}

func (self *DataBase) SetDeletedByUser(user interface{}) {
	self.db.Exec("UPDATE messages set deleted = 1 where name=?", user)
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

func (self *DataBase) GetLastUserById(limit int, start int, id int) []RowMap {

	firstImage := self.GetMessageById(int64(id));

	resp := append([]RowMap{}, firstImage)

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " where deleted=0 and name=? and id!=? order by id desc limit ?,?", firstImage["name"], id, start, limit)

	if(err != nil) {
		self.processError(err)
	}

	for _, row := range self.MultipleMessageScan(s) {
		resp = append(resp, row)
	}

	return resp
}

func (self *DataBase) GetLastUser(limit int, start int, user string) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " where deleted=0 and name=? order by id desc limit ?,?", user, start, limit)

	if(err != nil) {
		self.processError(err)
	}

	return self.MultipleMessageScan(s)
}

func (self *DataBase) getRandom(limit int) []RowMap {

	s, err := self.db.Query(
		MESSAGE_MAIN_QUERY + " WHERE deleted = 0 ORDER BY RANDOM() LIMIT ?", limit)

	if(err != nil) {
		self.processError(err)
	}

	return self.MultipleMessageScan(s)
}

func (self *DataBase) getTopUsers(limit int, exclude []string) []RowMap {

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
	var source string

	s.Scan(&id, &name, &url, &mature, &source)

	row["id"] = id
	row["name"] = name
	row["url"] = url
	row["mature"] = mature
	row["source"] = source

	return row
}

func (self *DataBase) SingleMessageScan(s *sql.Row) RowMap {
	row := make(RowMap)

	var id int
	var name string
	var url string
	var mature int
	var source string

	s.Scan(&id, &name, &url, &mature, &source)

	row["id"] = id
	row["name"] = name
	row["url"] = url
	row["mature"] = mature
	row["source"] = source

	return row
}

func (self *DataBase) processError(err error) {
	log.Println("Failed to execute query. Reason: ", err)
}

func (self *DataBase) Close() {
	self.db.Close();
}