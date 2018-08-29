package eosiodb

import (
	"errors"
	"github.com/eosspark/eos-go/db/storm"
	"path/filepath"
)

type Database struct {
	db   *storm.DB
	path string
	file string
	rw   bool // XXX read only or read write
}

func (db *Database) checkState() error {
	if !db.rw {
		return errors.New("read only")
	}
	return nil
}

func (db *Database) Insert(data interface{}) error {
	err := db.checkState()
	if err != nil {
		return err
	}
	return db.db.Save(data)
}

// all index
func (db *Database) ByIndex(fieldName string, to interface{}) error {
	return db.db.AllByIndex(fieldName, to)
}

func (db *Database) All(data interface{}) error {
	return db.db.All(data)
}

// uinque index
func (db *Database) Find(fieldName string, value interface{}, to interface{}) error {
	return db.db.One(fieldName, value, to)
}

// index==name
func (db *Database) Get(fieldName string, fieldValue interface{}, to interface{}) error {
	return db.db.Find(fieldName, fieldValue, to)
}

func (db *Database) Remover(item interface{}) error {
	err := db.checkState()
	if err != nil {
		return err
	}
	return db.db.DeleteStruct(item) // 	db.db.DeleteStruct ?
}

func (db *Database) UpdateField(data interface{}, fieldName string, value interface{}) error {
	err := db.checkState()
	if err != nil {
		return err
	}
	return db.db.UpdateField(data, fieldName, value)
}

func (db *Database) Update(old interface{}, fn func(interface{}) error) error {
	err := db.checkState()
	if err != nil {
		return err
	}
	err = fn(old)
	if err != nil {
		return err
	}
	return db.update(old)
}

func (db *Database) update(item interface{}) error {
	return db.db.Update(item)
}

func (db *Database) Close() {
	db.db.Close()
}

func NewDatabase(path string, fname string, rw bool /*read and  write*/) (*Database, error) {
	dir := filepath.Join(path, fname)
	db, err := storm.Open(dir)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:   db,
		path: path,
		file: fname,
		rw:   rw,
	}, nil
}

/////////////////////////////////////////////////////// Session //////////////////////////////////////////////////////////

type Session struct {
	node storm.Node
	db   *Database
}

func (db *Database) Start_Session() (*Session, error) {
	node, err := db.db.Begin(db.rw)
	if err != nil {
		return nil, err
	}
	return &Session{node: node, db: db}, nil
}

func (session *Session) Reset_Session() (*Session, error) {
	return session.db.Start_Session()
}

func (session *Session) Undo() error {
	return session.node.Rollback()
}

func (session *Session) Commit() error {
	return session.node.Commit()
}

func (session *Session) Get(fieldName string, fieldValue interface{}, to interface{}) error {
	return session.node.Find(fieldName, fieldValue, to)
}

func (session *Session) Find(fieldName string, value interface{}, to interface{}) error {
	return session.node.One(fieldName, value, to)
}

func (session *Session) Insert(data interface{}) error {
	return session.node.Save(data)
}

func (session *Session) Remover(item interface{}) error {
	return session.node.DeleteStruct(item)
}

func (session *Session) ByIndex(fieldName string, to interface{}) error {
	return session.node.AllByIndex(fieldName, to)
}

func (session *Session) All(data interface{}) error {
	return session.node.All(data)
}

func (session *Session) UpdateField(data interface{}, fieldName string, value interface{}) error {
	return session.node.UpdateField(data, fieldName, value)
}

func (session *Session) Update(old interface{}, fn func(interface{}) error) error {
	err := fn(old)
	if err != nil {
		return err
	}
	return session.update(old)
}

func (session *Session) update(item interface{}) error {
	return session.node.Update(item)
}
