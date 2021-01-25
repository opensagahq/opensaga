package repositories

func NewStmt(query string, args ...interface{}) *Stmt {
	return &Stmt{
		query: query,
		args:  args,
	}
}

func (stmt *Stmt) Query() string {
	return stmt.query
}

func (stmt *Stmt) Args() []interface{} {
	return stmt.args
}

type Stmt struct {
	query string
	args  []interface{}
}
