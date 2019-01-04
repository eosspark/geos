package chain

import (
	"github.com/eosspark/eos-go/database"
)

type MaybeSession struct {
	MSession *database.Session
	Valid    bool
}

func NewMaybeSession(db database.DataBase) *MaybeSession {
	s := MaybeSession{}
	s.MSession = db.StartSession()
	return &s
}

func NewMaybeSession2() *MaybeSession {
	s := MaybeSession{}
	return &s
}

func (m *MaybeSession) Squash() {
	if m.MSession != nil {
		m.MSession.Squash()
	}
}

func (m *MaybeSession) Push() {
	if m.MSession != nil {
		m.MSession.Push()
	}
}

func (m *MaybeSession) Undo() {
	if m.MSession != nil {
		m.MSession.Undo()
	}
}
