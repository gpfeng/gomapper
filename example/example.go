// Author:  fengguangpu
// Date:    2014/09/22

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gomapper"
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

	xmlFile = "sqlmap.xml"
)

type Record struct {
	Id            int64
	FirstName     string
	LastName      string
	CreatedAt     time.Time
	EmailVerified bool
}

func GetMySQLConnection() (*sql.DB, error) {
	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", user, pass, host, port, db, "charset=utf8&parseTime=True")
	return sql.Open("mysql", ds)
}

func main() {
	db, err := GetMySQLConnection()
	if err != nil {
		fmt.Printf("connect to MySQL failed: %s\n", err.Error())
		return
	}

	// create instance of GoMapper
	mapper, err := gomapper.NewGoMapperByFile(db, xmlFile)
	if err != nil {
		fmt.Printf("NewGoMapper failed: %s\n", err.Error())
		return
	}

	_, err = mapper.DB.Exec(sql_drop)
	if err != nil {
		fmt.Printf("drop table failed :%s\n", err.Error())
		return
	}

	// create table for test
	_, err = mapper.DB.Exec(sql_create)
	if err != nil {
		fmt.Printf("create table failed :%s\n", err.Error())
		return
	}

	fmt.Println("---- Insert test")

	// insert record, normal case
	res, err := mapper.Insert("insertStmt", "Lilei", "LL", true)
	if err != nil {
		fmt.Printf("Insert failed: %s\n", err.Error())
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("LastInsertId() faild: %s\n", err.Error())
		return
	}
	fmt.Printf("Insert successed, last insert id: %d\n", id)

	// insert record, pass in struct, should success
	res, err = mapper.Insert("insertStmt", Record{FirstName: "David", LastName: "YY", EmailVerified: false})
	if err != nil {
		fmt.Printf("Insert failed: %s\n", err.Error())
		return
	}
	id, err = res.LastInsertId()
	if err != nil {
		fmt.Printf("LastInsertId() faild: %s\n", err.Error())
		return
	}
	fmt.Printf("Insert successed, last insert id: %d\n", id)

	fmt.Println("---- Get test")

	// test Get() and Scan(), normal case
	record := Record{Id: 1, EmailVerified: true}
	err = mapper.Get("selectNamesByIdEmail", record.Id, record.EmailVerified).Scan(&record.FirstName, &record.LastName, &record.CreatedAt)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		return
	}
	fmt.Printf("Record: %v\n", record)

	// test Get() and Scan(), pass in struct
	record = Record{Id: 2, EmailVerified: false}
	err = mapper.Get("selectNamesByIdEmail", record).Scan(&record)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		return
	}
	fmt.Printf("Record: %v\n", record)

	fmt.Println("---- Select test")

	// test Select() and Scan(), normal case
	rows, err := mapper.Select("selectAll")
	if err != nil {
		fmt.Printf("GoMapper Select failed: %s\n", err.Error())
		return
	}

	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec.Id, &rec.FirstName, &rec.LastName, &rec.EmailVerified, &rec.CreatedAt)
		if err != nil {
			fmt.Printf("GoMapper Scan failed: %s\n", err.Error())
			return
		}
		fmt.Printf("Record: %v\n", rec)
	}
	rows.Close()

	// test Select() and Scan(), pass in struct
	rows, err = mapper.Select("selectAll")
	if err != nil {
		fmt.Printf("GoMapper Select failed: %s\n", err.Error())
		return
	}

	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec)
		if err != nil {
			fmt.Printf("GoMapper Scan failed: %s\n", err.Error())
			return
		}
		fmt.Printf("Record: %v\n", rec)
	}
	rows.Close()

	fmt.Println("---- Update test")

	// test update, normal case
	res, err = mapper.Update("updateById", "First", "Last", false, 1)
	if err != nil {
		fmt.Printf("GoMapper Update failed: %s\n", err.Error())
		return
	}

	var first, last string
	err = mapper.Get("selectNamesById", 1).Scan(&first, &last)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		return
	}
	fmt.Printf("Record: First(%s), Last(%s)\n", first, last)

	// test update, pass in struct
	res, err = mapper.Update("updateById", Record{Id: 1, FirstName: "Lilei", LastName: "LL"})
	if err != nil {
		fmt.Printf("GoMapper Update failed: %s\n", err.Error())
		return
	}

	var rec Record
	err = mapper.Get("selectNamesById", 1).Scan(&rec)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		return
	}
	fmt.Printf("Record: First(%s), Last(%s)\n", rec.FirstName, rec.LastName)

	fmt.Println("---- Delete test")

	// test delete, normal case
	res, err = mapper.Delete("deleteById", 1)
	if err != nil {
		fmt.Printf("GoMapper Delete failed: %s\n", err.Error())
		return
	}
	fmt.Printf("Delete(id = 1) successed\n")

	// test delete, pass in struct
	res, err = mapper.Delete("deleteById", Record{Id: 2})
	if err != nil {
		fmt.Printf("GoMapper Delete failed: %s\n", err.Error())
		return
	}
	fmt.Printf("Delete(id = 2) successed\n")

	// get number of records after deletes
	var cnt int
	err = mapper.Get("selectCnt").Scan(&cnt)
	if err != nil {
		fmt.Printf("GoMapper Get failed: %s\n", err.Error())
		return
	}
	fmt.Printf("There are %d records\n", cnt)
}
