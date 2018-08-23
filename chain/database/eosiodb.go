package database

import (
	"errors"
	"fmt"
	"github.com/asdine/storm"
	"log"
	"path/filepath"
)

type Database struct {
	db   *storm.DB
	node storm.Node
	path string
	file string
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (db *Database) Insert(data interface{}) error {
	if db.node == nil {
		return db.db.Save(data)
	}
	return db.node.Save(data)
}

// all index
func (db *Database) Get_index(fieldName string, to interface{}) error {
	return db.db.AllByIndex(fieldName, to)
}

// uinque index
func (db *Database) Find(fieldName string, value string, to interface{}) error {
	return db.db.One(fieldName, value, to)
}

// index
func (db *Database) Get(fieldName string, fieldValue interface{}, to interface{}) error {
	return db.db.Find(fieldName, fieldValue, to)
}

func (db *Database) Remover(item interface{}) error {
	if db.node == nil {
		return db.db.DeleteStruct(item) // 	db.db.DeleteStruct ?
	}
	return db.node.DeleteStruct(item) // 	db.db.DeleteStruct ?
}

func (db *Database) Update(item interface{}) error {
	if db.node == nil {
		return db.db.Update(item)
	}
	return db.node.Update(item)
}

//func (db *Database) Modify(old interface{}, new interface{}) error {
//node, err := db.db.Begin(true)
//if err != nil {
//return err
//}
//defer node.Rollback()
//err = node.Drop(old)
//if err != nil {
//return err
//}
//err = node.Save(new)
//if err != nil {
//return err
//}
//node.Commit()
//return nil
//}

func (db *Database) Start_Undo_Session() error {
	if db.node != nil {
		// log ?
		return errors.New(" Start_Undo_Session ")
	}
	node, err := db.db.Begin(true)
	db.node = node
	return err
}

func (db *Database) Commit() {
	db.node.Commit()
}

func (db *Database) Close() {
	if db.node != nil {
		db.Undo()
	}
	db.db.Close()
}

func (db *Database) Undo() {
	if db.node == nil {
		return
	}
	db.node.Rollback()
	db.node = nil
}

func NewDatabase(path string, fname string) *Database {
	dir := filepath.Join(path, fname)
	log.Println(dir)
	db, err := storm.Open(fname)
	CheckError(err)
	return &Database{
		db:   db,
		path: path,
		file: fname,
		node: nil,
	}
}

type Account struct {
	ID   int    `storm:"id,increment"`
	Name string `storm:"index"`
}

func (account *Account) To_String() {
	str := "  "
	fmt.Println(account.ID, str, account.Name)
}
