package chain

import (
	"github.com/eosspark/eos-go/common"
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
	if !common.Empty(m.MSession) {
		m.MSession.Squash()
	}
}

func (m *MaybeSession) Push() {
	if !common.Empty(m.MSession) {
		//m.MSession.Push()
	}
}

func (m *MaybeSession) Undo() {
	if !common.Empty(m.MSession) {
		m.MSession.Undo()
	}
}
