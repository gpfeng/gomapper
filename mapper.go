// Author:  fengguangpu
// Date:    2014/09/22

package gomapper

import "database/sql"
import "regexp"

// thread-safe, can be used by cocurrent go routines
var sqlVarsRegexp = regexp.MustCompile(`(\$\d+)|\?`)

// TODO: Prepare will be added in the future
type DbDriver interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// DB can be *sql.DB and *sql.Tx
type Mapper struct {
	DB      DbDriver
	sqlMap  *SqlMap
	logFunc func(format string, args ...interface{})
}

type GoMapper struct {
	Mapper
}

type GoMapperTx struct {
	Mapper
}

// API definitions

// Create an instance of *GoMapper*, defined in mapper.go
//func NewGoMapperByFile(db *sql.DB, xmlFilePath string) (*GoMapper, error)
//func NewGoMapper(db *sql.DB, xmlBytes []byte) (*GoMapper, error)

// Exported field of *GoMapper*
// DB: passed in by caller as an instance of *sql.DB, can be usesd directly to do begin/commit/rollback

// Get one row from db, defined in select.go
//func (r *gomapper.Row) Scan(dest ...interface{}) error
//func (m *Mapper) Get(id string, args ...interface{}) (row *GoMapper.Row)

// Select multi rows from db, defined in select.go
//func (rs *gomapper.Rows) Close() error
//func (rs *gomapper.Rows) Columns() ([]string, error)
//func (rs *gomapper.Rows) Err() error
//func (rs *gomapper.Rows) Next() bool
//func (rs *gomapper.Rows) Scan(dest ...interface{}) error
//func (m *Mapper) Select(id string, args ...interface{}) (rows *GoMapper.Rows, error)

// DML wrapper, defined in dml.go
//func (gm *GoMapper) Close() err error
//func (gm *GoMapper) Begin() (gmtx *GoMapperTx, err error)
//func (gmtx *GoMapperTx) Commit() error
//func (gmtx *GoMapperTx) Rollback() error
//func (m *Mapper) Insert(id string, args ...interface{}) (sql.Result, error)
//func (m *Mapper) Update(id string, args ...interface{}) (sql.Result, error)
//func (m *Mapper) Delete(id string, args ...interface{}) (sql.Result, error)

func NewGoMapperByFile(db *sql.DB, xmlFilePath string) (*GoMapper, error) {
	sqlMap, err := NewSqlMapByFile(xmlFilePath)
	if err != nil {
		return nil, err
	}

	m := new(GoMapper)
	m.DB, m.sqlMap = db, sqlMap
	return m, nil
}

func NewGoMapper(db *sql.DB, xmlBytes []byte) (*GoMapper, error) {
	sqlMap, err := NewSqlMap(xmlBytes)
	if err != nil {
		return nil, err
	}

	m := new(GoMapper)
	m.DB, m.sqlMap = db, sqlMap
	return m, nil
}

// set logger for gommaper
// added by fengguangpu on 2014-12-16
func (gm *GoMapper) SetLogger(logFunc func(format string, args ...interface{})) {
	gm.logFunc = logFunc
}
