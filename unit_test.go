// Author:  fengguangpu
// Date:    2014/09/22

package gomapper

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
	"time"
)

const (
	host = "localhost"
	port = 3306
	user = "root"
	pass = ""
	db   = "test"

	sql_drop   = "drop table if exists t"
	sql_create = "create table t(" +
		"id bigint auto_increment primary key," +
		"first_name varchar(64) default NULL," +
		"last_name varchar(64) default NULL," +
		"email_verified bool default false," +
		"created_at datetime default '0000-00-00 00:00:00'" +
		")engine=innodb charset=utf8"

	xmlFile = "example/sqlmap.xml"
)

type Record struct {
	Id            int64
	FirstName     string
	LastName      string
	CreatedAt     time.Time
	EmailVerified bool
}

type RecordXx struct {
	Id            int64
	FirstNameXx   string
	LastNameXx    string
	CreatedAt     time.Time
	EmailVerified bool
}

func GetMySQLConnection() (*sql.DB, error) {
	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", user, pass, host, port, db, "charset=utf8&parseTime=True")
	return sql.Open("mysql", ds)
}

func TestFormatSqlAndVars(t *testing.T) {
	fmt.Println("\n---------- TestFormatSqlAndVars ----------")
	sql := "SELECT id, a, b, c FROM t WHERE id=#{Id} and d=#{Data}"
	s, vars, err := FormatSqlAndVars(sql)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal()
	}
	fmt.Printf("formated sql: %s\n", s)
	fmt.Printf("variables: %v, len: %d\n", vars, len(vars))
}

func TestNewSqlMapper(t *testing.T) {
	fmt.Println("\n---------- TestNewSqlMapper ----------")

	sqlMap, err := NewSqlMapByFile(xmlFile)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal()
	}

	for k, v := range sqlMap.Sqls {
		fmt.Printf("%s:\t -> \t%v\n", k, v)
	}
}

func TestParseQueryArgs(t *testing.T) {
	fmt.Println("\n---------- TestParseQueryArgs ----------")
	names := []string{"Id", "FirstName", "CreatedAt", "EmailVerified"}
	record := Record{Id: 100, FirstName: "Lisa", CreatedAt: time.Now(), EmailVerified: true}

	args0, err := ParseQueryArgs(names, record.Id, record.FirstName, record.CreatedAt, record.EmailVerified)
	if err != nil {
		t.Fatal()
	}
	fmt.Printf("Args0: %v\n", args0)

	args1, err := ParseQueryArgs(names, record)
	if err != nil {
		t.Fatal()
	}
	fmt.Printf("Args1: %v\n", args1)
}

func TestSelectDmlTx(t *testing.T) {
	fmt.Println("\n---------- TestSelectDml ----------")
	db, err := GetMySQLConnection()
	if err != nil {
		fmt.Printf("connect to MySQL failed: %s\n", err.Error())
		t.Fatal()
	}

	_, err = db.Exec(sql_drop)
	if err != nil {
		fmt.Printf("drop table failed :%s\n", err.Error())
		t.Fatal()
	}

	// create table for test
	_, err = db.Exec(sql_create)
	if err != nil {
		fmt.Printf("create table failed :%s\n", err.Error())
		t.Fatal()
	}

	// create instance of GoMapper
	//mapper, err := NewGoMapperByFile(db, xmlFile)
	mapper, err := NewGoMapper(db, []byte(xmlsqls))
	if err != nil {
		fmt.Printf("NewGoMapper failed: %s\n", err.Error())
		t.Fatal()
	}
	defer mapper.Close()

	mapper.SetLogger(log.Printf)

	fmt.Println("---- Inset test")

	// insert record, normal case
	res, err := mapper.Insert("insertStmt", "Lilei", "LL", true)
	if err != nil {
		fmt.Printf("Insert failed: %s\n", err.Error())
		t.Fatal()
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("LastInsertId() faild: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Insert successed, last insert id: %d\n", id)

	// insert record, pass in struct, should success
	res, err = mapper.Insert("insertStmt", Record{FirstName: "David", LastName: "YY", EmailVerified: false})
	if err != nil {
		fmt.Printf("Insert failed: %s\n", err.Error())
		t.Fatal()
	}
	id, err = res.LastInsertId()
	if err != nil {
		fmt.Printf("LastInsertId() faild: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Insert successed, last insert id: %d\n", id)

	// insert record, pass in struct, should fail
	res, err = mapper.Insert("insertStmt", RecordXx{FirstNameXx: "David", LastNameXx: "YY", EmailVerified: true})
	if err == nil {
		fmt.Printf("insert should fail, but successed\n")
		t.Fatal()
	} else {
		fmt.Printf("Expected failure! error msg: %s\n", err.Error())
	}

	fmt.Println("---- Get test")

	// test Get() and Scan(), normal case
	record := Record{Id: 1, EmailVerified: true}
	err = mapper.Get("selectNamesByIdEmail", record.Id, record.EmailVerified).Scan(&record.FirstName, &record.LastName, &record.CreatedAt)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Record: %v\n", record)

	// test Get() and Scan(), pass in struct
	record = Record{Id: 2, EmailVerified: false}
	err = mapper.Get("selectNamesByIdEmail", record).Scan(&record)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Record: %v\n", record)

	fmt.Println("---- Select test")

	// test Select() and Scan(), normal case
	rows, err := mapper.Select("selectAll")
	if err != nil {
		fmt.Printf("GoMapper Select failed: %s\n", err.Error())
		t.Fatal()
	}

	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec.Id, &rec.FirstName, &rec.LastName, &rec.EmailVerified, &rec.CreatedAt)
		if err != nil {
			fmt.Printf("GoMapper Scan failed: %s\n", err.Error())
			t.Fatal()
		}
		fmt.Printf("Record: %v\n", rec)
	}
	rows.Close()

	// test Select() and Scan(), pass in struct
	rows, err = mapper.Select("selectAll")
	if err != nil {
		fmt.Printf("GoMapper Select failed: %s\n", err.Error())
		t.Fatal()
	}

	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec)
		if err != nil {
			fmt.Printf("GoMapper Scan failed: %s\n", err.Error())
			t.Fatal()
		}
		fmt.Printf("Record: %v\n", rec)
	}
	rows.Close()

	fmt.Println("---- Update test")

	// test update, normal case
	res, err = mapper.Update("updateById", "First", "Last", false, 1)
	if err != nil {
		fmt.Printf("GoMapper Update failed: %s\n", err.Error())
		t.Fatal()
	}

	var first, last string
	err = mapper.Get("selectNamesById", 1).Scan(&first, &last)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Record: First(%s), Last(%s)\n", first, last)

	// test update, pass in struct
	res, err = mapper.Update("updateById", Record{Id: 1, FirstName: "Lilei", LastName: "LL"})
	if err != nil {
		fmt.Printf("GoMapper Update failed: %s\n", err.Error())
		t.Fatal()
	}

	var rec Record
	err = mapper.Get("selectNamesById", 1).Scan(&rec)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Record: First(%s), Last(%s)\n", rec.FirstName, rec.LastName)

	fmt.Println("---- Delete test")

	// test delete, normal case
	res, err = mapper.Delete("deleteById", 1)
	if err != nil {
		fmt.Printf("GoMapper Delete failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Delete(id = 1) successed\n")

	// test delete, pass in struct
	res, err = mapper.Delete("deleteById", Record{Id: 2})
	if err != nil {
		fmt.Printf("GoMapper Delete failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Delete(id = 2) successed\n")

	// get number of records after deletes
	var cnt int
	err = mapper.Get("selectCnt").Scan(&cnt)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("There are %d records, should be 0\n", cnt)

	fmt.Println("---- Transaction rollback test")

	// transaction rollback test
	tx, err := mapper.Begin()
	if err != nil {
		t.Fatal("mapper.Begin failed: %s\n", err.Error())
	}

	res, err = tx.Insert("insertStmt", "Lilei", "LL", true)
	if err != nil {
		t.Fatal("tx.Insert failed: %s\n", err.Error())
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal("tx.Rollback failed: %s\n", err.Error())
	}

	// get number of records after deletes
	err = mapper.Get("selectCnt").Scan(&cnt)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("There are %d records after rollback, should be 0\n", cnt)

	fmt.Println("---- Transaction commit test")

	// transaction commit test
	tx, err = mapper.Begin()
	if err != nil {
		t.Fatal("mapper.Begin failed: %s\n", err.Error())
	}

	res, err = tx.Insert("insertStmt", Record{FirstName: "David", LastName: "YY", EmailVerified: false})
	if err != nil {
		t.Fatal("tx.Insert failed: %s\n", err.Error())
	}

	// test Get() and Scan(), normal case
	record = Record{Id: 4, EmailVerified: false}
	err = tx.Get("selectNamesByIdEmail", record.Id, record.EmailVerified).Scan(&record.FirstName, &record.LastName, &record.CreatedAt)
	if err != nil {
		fmt.Printf("GoMapperTx Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Record: %v\n", record)

	// test Get() and Scan(), pass in struct
	record = Record{Id: 4, EmailVerified: false}
	err = tx.Get("selectNamesByIdEmail", record).Scan(&record)
	if err != nil {
		fmt.Printf("GoMapperTx Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("Record: %v\n", record)

	// test Select() and Scan(), normal case
	rows, err = tx.Select("selectAll")
	if err != nil {
		fmt.Printf("GoMapperTx Select failed: %s\n", err.Error())
		t.Fatal()
	}

	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec.Id, &rec.FirstName, &rec.LastName, &rec.EmailVerified, &rec.CreatedAt)
		if err != nil {
			fmt.Printf("GoMapperTx Scan failed: %s\n", err.Error())
			t.Fatal()
		}
		fmt.Printf("Record: %v\n", rec)
	}
	rows.Close()

	// test Select() and Scan(), pass in struct
	rows, err = tx.Select("selectAll")
	if err != nil {
		fmt.Printf("GoMapperTx Select failed: %s\n", err.Error())
		t.Fatal()
	}

	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec)
		if err != nil {
			fmt.Printf("GoMapperTx Scan failed: %s\n", err.Error())
			t.Fatal()
		}
		fmt.Printf("Record: %v\n", rec)
	}
	rows.Close()

	err = tx.Commit()
	if err != nil {
		t.Fatal("tx.Commit failed: %s\n", err.Error())
	}

	// get number of records after deletes
	err = mapper.Get("selectCnt").Scan(&cnt)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		t.Fatal()
	}
	fmt.Printf("There are %d records after commit, should be 1\n", cnt)
}
