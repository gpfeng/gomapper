// Author:  fengguangpu
// Date:    2014/09/22

package gomapper

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (gm *GoMapper) Close() error {
	db, ok := gm.DB.(*sql.DB)
	if ok && db != nil {
		return db.Close()
	} else {
		return errors.New("Invalid *sql.DB instance")
	}
}

func (gm *GoMapper) Begin() (*GoMapperTx, error) {
	db, ok := gm.DB.(*sql.DB)
	if ok && db != nil {
		var err error
		gmtx := new(GoMapperTx)
		gmtx.sqlMap = gm.sqlMap
		gmtx.logFunc = gm.logFunc
		gmtx.DB, err = db.Begin()
		return gmtx, err
	} else {
		return nil, errors.New("Invalid *sql.DB instance")
	}
}

func (gmtx *GoMapperTx) Commit() error {
	tx, ok := gmtx.DB.(*sql.Tx)
	if ok && tx != nil {
		return tx.Commit()
	}
	return errors.New("Invalid *sql.Tx instance")
}

func (gmtx *GoMapperTx) Rollback() error {
	tx, ok := gmtx.DB.(*sql.Tx)
	if ok && tx != nil {
		return tx.Rollback()
	}
	return errors.New("Invalid *sql.Tx instance")
}

// underlying implementation of Insert/Update/Delete
func (m *Mapper) dml(id string, args ...interface{}) (sql.Result, error) {
	element, err := m.sqlMap.Get(id)
	if err != nil {
		return nil, err
	}

	var start int64
	if m.logFunc != nil {
		start = time.Now().UnixNano()
	}

	sqlArgs, err := ParseQueryArgs(element.Vars, args...)
	if err != nil {
		return nil, err
	}

	if m.logFunc != nil && start != 0 {
		end := time.Now().UnixNano()
		elaspedMs := float64(end-start) / float64(1000000)
		sqlFormat := sqlVarsRegexp.ReplaceAllString(element.Sql, "'%v'")
		sqlCommand := fmt.Sprintf(sqlFormat, sqlArgs...)
		m.logFunc("\033[36;1m[%.2fms]\033[0m %s\n", elaspedMs, sqlCommand)
	}

	result, err := m.DB.Exec(element.Sql, sqlArgs...)

	return result, err
}

// args can be Struct
func (m *Mapper) Insert(id string, args ...interface{}) (sql.Result, error) {
	return m.dml(id, args...)
}

// args can be Struct
func (m *Mapper) Update(id string, args ...interface{}) (sql.Result, error) {
	return m.dml(id, args...)
}

// args can be Struct
func (m *Mapper) Delete(id string, args ...interface{}) (sql.Result, error) {
	return m.dml(id, args...)
}
