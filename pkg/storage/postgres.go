package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lock sync.Mutex

//sql

const AdminINS = "INSERT INTO admins (id, telegram_name, first_name, chat_id) VALUES ($1,$2,$3,$4)"
const AdminSEL = "SELECT * FROM admins"

const RequestINS = "INSERT INTO requests (user_tgname, user_first_name, user_chat_id, status) VALUES ($1,$2,$3,$4)"
const RequestSEL = "SELECT * FROM requests"

const FindAns = "SELECT * FROM temperature WHERE key = $1"

type DB struct {
	Sql           *sql.DB
	Stmt          *sql.Stmt
	adminBuffer   []Admin
	requestBuffer []Request
}

type Admin struct {
	Id           int
	TelegramName string
	FirstName    string
	ChatID       int64
}

type Request struct {
	Id           int
	TelegramName string
	FirstName    string
	ChatID       int64
	Status       *bool
}

type SQLBuilderArgs struct {
	Query       string
	Table       string
	Columns     []string
	WhereString string
}

type Temp struct {
	Key         string
	Temperature string
	Params      []string
	Logic       string
}

func NewDB(dbFile string) (*DB, error) {
	lock.Lock()
	defer lock.Unlock()
	sqlDB, err := sql.Open("postgres", dbFile)
	if err != nil {
		return nil, err
	}

	db := DB{
		Sql:           sqlDB,
		adminBuffer:   make([]Admin, 0, 5),
		requestBuffer: make([]Request, 0, 5),
	}
	time.Sleep(1 * time.Millisecond)
	return &db, nil
}

func (db *DB) ChangeStmt(rim string) error {
	stmt, err := db.Sql.Prepare(rim)
	if err != nil {
		return err
	}
	db.Stmt = stmt
	return nil
}

// add
func (db *DB) AddAdmin(admin Admin) error {
	if len(db.adminBuffer) == cap(db.adminBuffer) {
		return errors.New("tutor buffer is full")
	}

	db.adminBuffer = append(db.adminBuffer, admin)
	if len(db.adminBuffer) == cap(db.adminBuffer) {
		if err := db.FlushAdmin(); err != nil {
			return fmt.Errorf("unable to flush tutor: %w", err)
		}
	}

	return nil
}

func (db *DB) AddRequest(request Request) error {
	if len(db.requestBuffer) == cap(db.requestBuffer) {
		return errors.New("tutor buffer is full")
	}

	db.requestBuffer = append(db.requestBuffer, request)
	if len(db.requestBuffer) == cap(db.requestBuffer) {
		if err := db.FlushRequest(); err != nil {
			return fmt.Errorf("unable to flush tutor: %w", err)
		}
	}

	return nil
}

func (db *DB) GetAdmins() ([]Admin, error) {
	var Admins []Admin
	rows, err := db.Sql.Query(AdminSEL)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := Admin{}
		if err = rows.Scan(&t.Id, &t.TelegramName, &t.FirstName, &t.ChatID); err != nil {
			return nil, err
		}
		Admins = append(Admins, t)
	}
	return Admins, err
}
func (db *DB) GetRequests(wherestring string) ([]Request, error) {
	var Requests []Request
	rows, err := db.Sql.Query(RequestSEL + " " + wherestring)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := Request{}
		if err = rows.Scan(&t.Id, &t.TelegramName, &t.FirstName, &t.ChatID, &t.Status); err != nil {
			return nil, err
		}
		Requests = append(Requests, t)
	}
	return Requests, err
}

func (db *DB) FindAnw(rim string) (ret []Temp, err error) {
	if _, err = db.Sql.Query("SELECT * FROM temperature WHERE key = $1", rim); err != nil {
		return nil, err
	}
	rows, _ := db.Sql.Query("SELECT * FROM temperature WHERE key = $1", rim)
	defer rows.Close()

	temperatures := []Temp{}
	for rows.Next() {
		t := Temp{}
		err = rows.Scan(&t.Key, &t.Temperature, &t.Params, &t.Logic)
		if err != nil {
			fmt.Println(err)
			continue
		}
		temperatures = append(temperatures, t)
	}
	return ret, err
}

// flush
func (db *DB) FlushAdmin() error {
	if err := db.ChangeStmt(AdminINS); err != nil {
		return err
	}

	tx, err := db.Sql.Begin()
	if err != nil {
		return err
	}

	for _, admin := range db.adminBuffer {
		_, err = tx.Stmt(db.Stmt).Exec(admin.Id, admin.TelegramName, admin.FirstName, admin.ChatID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	db.adminBuffer = db.adminBuffer[:0]
	return tx.Commit()
}
func (db *DB) FlushRequest() error {
	if err := db.ChangeStmt(RequestINS); err != nil {
		return err
	}

	tx, err := db.Sql.Begin()
	if err != nil {
		return err
	}

	for _, request := range db.requestBuffer {
		_, err = tx.Stmt(db.Stmt).Exec(request.TelegramName, request.FirstName, request.ChatID, sql.NullBool{})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	db.requestBuffer = db.requestBuffer[:0]
	return tx.Commit()
}

func (db *DB) SQLBuilder(args SQLBuilderArgs) (string, error) {
	var result string
	switch args.Query {
	case "update":
		var columnArr []string
		var columnI int
		for columnI = range args.Columns {
			columnArr = append(columnArr, args.Columns[columnI]+"=$"+strconv.Itoa(columnI+1))
		}
		column := strings.Join(columnArr, ",")
		result = "UPDATE" + " " + args.Table + " " + "SET" + " " + column
		if args.WhereString != "" {
			result = result + " " + args.WhereString + "=$" + strconv.Itoa(columnI+2)
		}
	case "insert":
		column := strings.Join(args.Columns, ",")
		var indexes []string
		for index := range args.Columns {
			indexes = append(indexes, "$"+strconv.Itoa(index+1))
		}
		form := strings.Join(indexes, ",")
		//INSERT INTO tutor (ol_id, login, password, first_name, curator_id) VALUES ($1,$2,$3,$4,$5)
		result = "INSERT INTO" + " " + args.Table + " " + "(" + column + ")" + " " + "VALUES" + " " + "(" + form + ")"
	}
	return result, nil
}

// closers
func (db *DB) Close() error {
	defer func() {
		db.Stmt.Close()
		db.Sql.Close()
	}()

	if err := db.FlushAdmin(); err != nil {
		return err
	}

	return nil
}
