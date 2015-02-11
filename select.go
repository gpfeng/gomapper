// Author:  fengguangpu
// Date:    2014/09/22

package gomapper

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// for single row query
// all fields are not allowed to access by other packages
type Row struct {
	mapper *Mapper
	sqlId  string
	args   []interface{}
}

// used to convert table field name in database to struct field name in golang
// example: "user_name"  --> "UserName"
func SnakeToUpperCamel(s string) string {
	buf := bytes.NewBufferString("")
	for _, v := range strings.Split(s, "_") {
		if len(v) > 0 {
			buf.WriteString(strings.ToUpper(v[:1]))
			buf.WriteString(v[1:])
		}
	}
	return buf.String()
}

// extract arguments from struct by names of variables
func ParseQueryArgs(vars []string, args ...interface{}) ([]interface{}, error) {
	varsNum := len(vars)
	argsNum := len(args)
	if varsNum == 0 || argsNum == 0 || argsNum > 1 {
		return args, nil
	}

	// the first argument can tell the types of passed in arguments
	arg := args[0]
	queryArgs := make([]interface{}, varsNum)

	value := reflect.ValueOf(arg)
	kind := value.Kind()

	switch kind {
	case reflect.Struct:
		for i, name := range vars {
			field := value.FieldByName(name)
			if field.IsValid() {
				queryArgs[i] = field.Interface()
			} else {
				return queryArgs, errors.New(fmt.Sprintf("struct has no field '%s'", name))
			}
		}
	default:
		// primitive type
		return args, nil
	}

	return queryArgs, nil
}

func ScanToStruct(rows *sql.Rows, arg interface{}) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// slice of addresses of target fields
	result := make([]interface{}, len(columns))
	for i, name := range columns {
		// find the field in the struct
		goName := SnakeToUpperCamel(name)
		field := reflect.ValueOf(arg).Elem().FieldByName(goName)
		if field.IsValid() {
			if field.CanAddr() {
				// locate the address of field, will be used in scan
				result[i] = field.Addr().Interface()
			} else {
				return errors.New("CanAddr() is false in ScanToStruct")
			}
		} else {
			return errors.New("IsValid() is false in ScanToStruct")
		}
	}
	return rows.Scan(result...)
}

// Get one row
func (m *Mapper) Get(id string, args ...interface{}) *Row {
	return &Row{sqlId: id, args: args, mapper: m}
}

// *Scan* is the only method defined in interface *Scaner* in package *sql*
// only the following types are allowed to be passed
//  int64/float64/bool/[]byte/string/time.Time/nil
//
// We override this method to support *struct*
func (r *Row) Scan(dest ...interface{}) error {
	element, err := r.mapper.sqlMap.Get(r.sqlId)
	if err != nil {
		return err
	}

	sqlArgs, err := ParseQueryArgs(element.Vars, r.args...)
	if err != nil {
		return err
	}

	var start int64
	if r.mapper.logFunc != nil {
		start = time.Now().UnixNano()
	}

	// Use Query instead QueryRow here to get names of selected columns
	rows, err := r.mapper.DB.Query(element.Sql, sqlArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if r.mapper.logFunc != nil && start != 0 {
		end := time.Now().UnixNano()
		elaspedMs := float64(end-start) / float64(1000000)
		sqlFormat := sqlVarsRegexp.ReplaceAllString(element.Sql, "'%v'")
		sqlCommand := fmt.Sprintf(sqlFormat, sqlArgs...)
		r.mapper.logFunc("\033[36;1m[%.2fms]\033[0m %s\n", elaspedMs, sqlCommand)
	}

	if !rows.Next() {
		return sql.ErrNoRows
	}

	// the first argument can tell whether target receiver is struct or map
	arg := dest[0]
	kind := reflect.Indirect(reflect.ValueOf(arg)).Kind()

	switch kind {
	case reflect.Struct:
		err = ScanToStruct(rows, arg)
	default:
		err = rows.Scan(dest...)
	}
	return err
}

// for mutli rows query
type Rows struct {
	rows *sql.Rows // not allowed to access by other packages
}

func (rs *Rows) Close() error {
	return rs.rows.Close()
}

func (rs *Rows) Columns() ([]string, error) {
	return rs.rows.Columns()
}

func (rs *Rows) Err() error {
	return rs.rows.Err()
}

func (rs *Rows) Next() bool {
	return rs.rows.Next()
}

// Select multi rows
func (m *Mapper) Select(id string, args ...interface{}) (*Rows, error) {
	element, err := m.sqlMap.Get(id)
	if err != nil {
		return nil, err
	}

	sqlArgs, err := ParseQueryArgs(element.Vars, args...)
	if err != nil {
		return nil, err
	}
	rows, err := m.DB.Query(element.Sql, sqlArgs...)

	return &Rows{rows: rows}, err
}

// *Scan* is the only method defined in interface *Scaner* in package *sql*
// only the following types are allowed to be passed
//  int64/float64/bool/[]byte/string/time.Time/nil
//
// We override this method to support *struct*
func (rs *Rows) Scan(dest ...interface{}) error {
	var err error
	arg := dest[0]

	kind := reflect.Indirect(reflect.ValueOf(arg)).Kind()

	switch kind {
	case reflect.Struct:
		err = ScanToStruct(rs.rows, arg)
	default:
		err = rs.rows.Scan(dest...)
	}

	return err
}
