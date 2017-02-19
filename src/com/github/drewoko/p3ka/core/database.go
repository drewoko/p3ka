package core

import (
	"database/sql"
	"log"
	"time"

	_ "gopkg.in/mattn/go-sqlite3.v1"
)

const MESSAGE_MAIN_QUERY string = "SELECT id, name, origurl as url, mature, source FROM messages"

const (
	ALL = iota
	PEKA2TV
	GOODGAME
)

type RowMap map[string]interface{}

type DataBase struct {
	db *sql.DB
}

func (d *DataBase) Init(path string) *DataBase {

	var err error
	d.db, err = sql.Open("sqlite3", path)

	if err != nil {
		log.Fatal("Failed to open database. Reason: ", err)
	}

	d.createTable()

	return d
}

func (d *DataBase) createTable() {
	d.db.Exec("CREATE TABLE messages (id integer not null primary key, date int, messageid int, name text, origurl text, deleted int DEFAULT 0)")
	d.db.Exec("ALTER TABLE messages ADD COLUMN mature int DEFAULT 0;")
	d.db.Exec("ALTER TABLE messages ADD COLUMN source text DEFAULT 'peka2tv';")
	d.db.Exec("CREATE TABLE ggChannels (id integer not null primary key, channel text not null)")
}

func (d *DataBase) IsExists(author string, url string) bool {
	s := d.db.QueryRow("SELECT COUNT(id) as cnt FROM messages WHERE origurl=?", url)
	var cnt int
	s.Scan(&cnt)
	return cnt > 0
}

func (d *DataBase) AddRow(messageid int64, author string, url string, mature bool, source string, deleted bool) {

	_, err := d.db.Exec(
		"INSERT INTO messages (messageid, date, name, origurl, mature, source, deleted) VALUES (?, ?, ?, ?, ?, ?, ?)",
		int(messageid), int(time.Now().Unix()), author, url, mature, source, deleted)

	if err != nil {
		d.processError(err)
	}
}

func (d *DataBase) AddRowGGChannel(channel string) {

	s := d.db.QueryRow("SELECT count(id) as cnt FROM ggChannels WHERE channel = ?", channel)
	var cnt int
	s.Scan(&cnt)

	if cnt == 0 {
		_, err := d.db.Exec(
			"INSERT INTO ggChannels (channel) VALUES (?)", channel)

		if err != nil {
			d.processError(err)
		}
	}
}

func (d *DataBase) GetGGChannels() []string {
	s, err := d.db.Query("SELECT channel FROM ggChannels")

	if err != nil {
		d.processError(err)
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

func (d *DataBase) GetMessageById(id int64) RowMap {
	s := d.db.QueryRow(MESSAGE_MAIN_QUERY+" where deleted = 0 and id=?", int(id))
	return d.SingleMessageScan(s)
}

func (d *DataBase) GetMessageByIdAndSource(id int64, source string) RowMap {
	s := d.db.QueryRow(MESSAGE_MAIN_QUERY+" where deleted = 0 and id=? and source=?", int(id), source)
	return d.SingleMessageScan(s)
}

func (d *DataBase) GetMessageByMessageId(id int64) RowMap {
	s := d.db.QueryRow(MESSAGE_MAIN_QUERY+" where deleted = 0 and messageid=?", int(id))
	return d.SingleMessageScan(s)
}

func (d *DataBase) SetDeleted(id interface{}) {
	d.db.Exec("UPDATE messages set deleted = 1 where id=?", id)
}

func (d *DataBase) SetDeletedByUser(user interface{}) {
	d.db.Exec("UPDATE messages set deleted = 1 where name=?", user)
}

func (d *DataBase) GetAll() []RowMap {
	s, err := d.db.Query(MESSAGE_MAIN_QUERY + " where deleted=0")

	if err != nil {
		d.processError(err)
	}

	return d.MultipleMessageScan(s)
}

func (d *DataBase) getFilterWhere(start bool, filter int) string {
	query := ""

	if filter == ALL {
		return query
	}

	if !start {
		query = " and "
	}

	if filter == PEKA2TV {
		query += "source='peka2tv'"
	} else {
		query += "source='goodgame'"
	}

	return query
}

func (d *DataBase) GetLast(limit int, start int, filter int) []RowMap {

	s, err := d.db.Query(
		MESSAGE_MAIN_QUERY+" where deleted=0 "+d.getFilterWhere(false, filter)+" order by id desc limit ?,?", start, limit)

	if err != nil {
		d.processError(err)
	}

	return d.MultipleMessageScan(s)
}

func (d *DataBase) GetLastUserById(limit int, start int, id int) []RowMap {

	firstImage := d.GetMessageById(int64(id))

	resp := []RowMap{}
	if firstImage["id"] != 0 {
		resp = append(resp, firstImage)
	}

	s, err := d.db.Query(
		MESSAGE_MAIN_QUERY+" where deleted=0 and name=? and id!=? order by id desc limit ?,?", firstImage["name"], id, start, limit)

	if err != nil {
		d.processError(err)
	}

	for _, row := range d.MultipleMessageScan(s) {
		resp = append(resp, row)
	}

	return resp
}

func (d *DataBase) GetLastUser(limit int, start int, user string) []RowMap {

	s, err := d.db.Query(
		MESSAGE_MAIN_QUERY+" where deleted=0 and name=? order by id desc limit ?,?", user, start, limit)

	if err != nil {
		d.processError(err)
	}

	return d.MultipleMessageScan(s)
}

func (d *DataBase) getRandom(limit int, filter int) []RowMap {

	s, err := d.db.Query(
		MESSAGE_MAIN_QUERY+" WHERE deleted = 0  "+d.getFilterWhere(false, filter)+"   ORDER BY RANDOM() LIMIT ?", limit)

	if err != nil {
		d.processError(err)
	}

	return d.MultipleMessageScan(s)
}

func (d *DataBase) getTop(limit int, exclude []string) []RowMap {

	var rows []RowMap

	s, err := d.db.Query(
		"select count(*) as cnt, name from messages where deleted = 0 group by name order by cnt desc LIMIT ?", limit)

	if err != nil {
		d.processError(err)
	}

	for s.Next() {
		row := make(RowMap)

		var cnt int
		var name string

		s.Scan(&cnt, &name)

		row["cnt"] = cnt
		row["name"] = name

		if !ContainsString(exclude, row["name"].(string)) {
			rows = append(rows, row)
		}
	}

	return rows
}

func (d *DataBase) getTopUsersBySource(limit int, source string, exclude []string) []RowMap {

	var rows []RowMap

	s, err := d.db.Query(
		"select count(*) as cnt, name from messages where deleted = 0 and source=? group by name order by cnt desc LIMIT ?", source, limit)

	if err != nil {
		d.processError(err)
	}

	for s.Next() {
		row := make(RowMap)

		var cnt int
		var name string

		s.Scan(&cnt, &name)

		row["cnt"] = cnt
		row["name"] = name

		if !ContainsString(exclude, row["name"].(string)) {
			rows = append(rows, row)
		}
	}

	return rows
}

func (d *DataBase) MultipleMessageScan(s *sql.Rows) []RowMap {
	var rows []RowMap

	for s.Next() {
		rows = append(rows, d.MultipleItemMessageScan(s))
	}

	s.Close()

	return rows
}

func (d *DataBase) MultipleItemMessageScan(s *sql.Rows) RowMap {
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

func (d *DataBase) SingleMessageScan(s *sql.Row) RowMap {
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

func (d *DataBase) processError(err error) {
	log.Println("Failed to execute query. Reason: ", err)
}

func (d *DataBase) Close() {
	d.db.Close()
}
