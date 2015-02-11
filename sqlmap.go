// Author:  fengguangpu
// Date:    2014/09/22

package gomapper

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type SqlType int

const (
	SQL_TYPE_BEGIN = iota
	SQL_SELECT
	SQL_INSERT
	SQL_UPDATE
	SQL_DELETE
	SQL_TYPE_END
)

func (st SqlType) String() string {
	switch st {
	case SQL_SELECT:
		return "SELECT"
	case SQL_INSERT:
		return "INSERT"
	case SQL_UPDATE:
		return "UPDATE"
	case SQL_DELETE:
		return "DELETE"
	default:
		return "Not supported sql type"
	}
}

type XmlSqlNode struct {
	Id  string `xml:"id,attr"`
	Sql string `xml:",chardata"`
}

type XmlSqls struct {
	XMLName xml.Name     `xml:"sqlmap"`
	Selects []XmlSqlNode `xml:"select"`
	Inserts []XmlSqlNode `xml:"insert"`
	Updates []XmlSqlNode `xml:"update"`
	Deletes []XmlSqlNode `xml:"delete"`
}

type SqlElement struct {
	Id   string   // unique name
	Sql  string   // sql statement
	Type SqlType  // type of statement: insert/update/delete/select
	Vars []string // names of variables that needed to be passed
}

type SqlMap struct {
	Sqls map[string]SqlElement
}

func (sm *SqlMap) InitMap() {
	sm.Sqls = make(map[string]SqlElement)
}

func (sm *SqlMap) Get(id string) (*SqlElement, error) {
	if val, ok := sm.Sqls[id]; ok {
		return &val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Sql statement specified by '%s' does not exist", id))
	}
}

// parse a XmlSqlNode and save to SqlMapper
// SqlType must be checked
func (sm *SqlMap) Add(node *XmlSqlNode, t SqlType) error {
	id := strings.Trim(node.Id, " ")

	// parse and check sql type
	sql := strings.Trim(strings.Replace(node.Sql, "\n", " ", -1), " ")
	err := CheckSqlType(sql, t)
	if err != nil {
		return err
	}

	// find variables and format sql statement
	sql, vars, err := FormatSqlAndVars(sql)
	if err != nil {
		return err
	}

	// add to SqlMapper
	sm.Sqls[id] = SqlElement{Id: id, Sql: sql, Type: t, Vars: vars}

	return nil
}

// sql must be trimed!
func CheckSqlType(sql string, t SqlType) error {
	// remove the leading tab, leading spaces are trimed already
	sql = strings.TrimLeft(sql, "\t")
	// check the first 10 bytes
	begin := strings.ToUpper(string([]byte(sql)[0:10]))
	if strings.HasPrefix(begin, t.String()) {
		return nil
	} else {
		return errors.New(fmt.Sprintf("[%s] is not %s statement", sql, t.String()))
	}
}

// find all variables defined as '#{Name}' with '?' in sql statement
// replace the variable names with '?'
func FormatSqlAndVars(sql string) (string, []string, error) {
	str := ""
	s := sql

	// 32 is enough for most statements
	vars := make([]string, 0, 32)

	for {
		if s == "" {
			break
		}
		bytes := []byte(s)

		i, j := strings.Index(s, "#{"), strings.IndexByte(s, '}')
		if i == -1 && j == -1 {
			str += string(bytes)
			return str, vars, nil
		} else if i != -1 && j != -1 && i < j {
			// replace name with ?
			str += string(bytes[:i]) + "?"
			// val name is between '#{' and '}', trim leading and trailing spaces
			vars = append(vars, strings.Trim(string(bytes[i+2:j]), " "))
			// repeat on the remaining string
			s = string(bytes[j+1:])
		} else {
			return str, vars, errors.New(fmt.Sprintf("Unmatched '#{' and '}' in %s", sql))
		}
	}
	return str, vars, nil
}

func NewSqlMapByFile(xmlFilePath string) (*SqlMap, error) {
	file, err := os.Open(xmlFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	xmlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return NewSqlMap(xmlBytes)
}

func NewSqlMap(xmlBytes []byte) (*SqlMap, error) {
	var sqls XmlSqls
	err := xml.Unmarshal(xmlBytes, &sqls)
	if err != nil {
		return nil, err
	}

	var mapper SqlMap
	mapper.InitMap()

	for _, v := range sqls.Selects {
		err = mapper.Add(&v, SQL_SELECT)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range sqls.Inserts {
		err = mapper.Add(&v, SQL_INSERT)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range sqls.Updates {
		err = mapper.Add(&v, SQL_UPDATE)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range sqls.Deletes {
		err = mapper.Add(&v, SQL_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return &mapper, nil
}
