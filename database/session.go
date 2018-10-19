package database


/////////////////////////////////////////////////////// Session  //////////////////////////////////////////////////////////
type Session struct {
	db      	DataBase
	apply   	bool
	revision 	int64
}

func (session *Session) Commit() {
	if !session.apply {
		// log ?
		return
	}
	//	version := session.version
	//	session.db.commit(version)
	session.apply = false
}

func (session *Session) Squash() {
	if !session.apply {
		return
	}
	//	session.db.squash()
	session.apply = false
}

func (session *Session) Undo() {
	if !session.apply {
		return
	}
	//	session.db.undo()
	session.apply = false
}
