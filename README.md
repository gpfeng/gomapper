# Gomapper
A lightweight ORM of Golang for MySQL

# Introduction

suppose we have table 't' defined as

```
mysql> show create table t\G
*************************** 1. row ***************************
   Table: t
Create Table: CREATE TABLE `t` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(64) DEFAULT NULL,
  `last_name` varchar(64) DEFAULT NULL,
  `email_verified` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB
```

we can use this struct to map a row in database

```
type Record struct {
	Id            int64
	FirstName     string
	LastName      string
	CreatedAt     time.Time
	EmailVerified bool
}
```
before reading/writing data from/to database using this struct, we should define sql mapping in xml first:

```
<?xml version="1.0" encoding="utf-8"?>
<sqlmap>
    <select id="selectAll">
        select id, first_name, last_name, email_verified, created_at from t
    </select>
    <select id="selectAllById">
        SELECT * FROM `t` WHERE id=#{Id}
    </select>
    <select id="selectNamesById">
        SELECT first_name, last_name FROM t WHERE id=#{Id}
    </select>
    <select id="selectNamesByIdEmail">
        SELECT first_name, last_name, created_at FROM t WHERE id=#{Id} and email_verified=#{EmailVerified}
    </select>
    <select id="selectCnt">
        SELECT COUNT(*) FROM t
    </select>
    <insert id="insertStmt">
        INSERT INTO t(first_name, last_name, email_verified, created_at)
        VALUES(#{FirstName}, #{LastName}, #{EmailVerified}, NOW())
    </insert>
    <update id="updateById">
        UPDATE t SET
        first_name=#{FirstName},
        last_name=#{LastName},
        email_verified=#{EmailVerified}
        WHERE id=#{Id}
    </update>
    <delete id="deleteById">
        DELETE FROM t WHERE id=#{Id}
    </delete>
</sqlmap>
```

## New gommaper instance
```	
ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", user, pass, host, port, db, "charset=utf8&parseTime=True")
db, err := sql.Open("mysql", ds)
if err != nil {
	fmt.Printf("connect to MySQL failed: %s\n", err.Error())
	return
}
mapper, err := NewGoMapperByFile(db, xmlFile)
..

```

## Query one row

```
record = Record{Id: 2, EmailVerified: false}
err = mapper.Get("selectNamesByIdEmail", record).Scan(&record)
```
or

```
var first, last string
err = mapper.Get("selectNamesById", 1).Scan(&first, &last)
```

## Query multi rows

```
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
```
or

```
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

```

## DML(insert/update/delete)
### Insert:
```
res, err = mapper.Insert("insertStmt", Record{FirstName: "David", LastName: "YY", EmailVerified: false})
if err != nil {
	fmt.Printf("Insert failed: %s\n", err.Error())
	return
}
id, err = res.LastInsertId()
```
or

```
res, err := mapper.Insert("insertStmt", "Lilei", "LL", true)
if err != nil {
	fmt.Printf("Insert failed: %s\n", err.Error())
	return
}
id, err := res.LastInsertId()
```
### Update:
```
res, err = mapper.Update("updateById", Record{Id: 1, FirstName: "Lilei", LastName: "LL"})
```
or

```
res, err = mapper.Update("updateById", "First", "Last", false, 1)
```
### Delete:
```
res, err = mapper.Delete("deleteById", Record{Id: 2})
```
or

```
res, err = mapper.Delete("deleteById", 1)
```

## Transaction support
```
tx, err := mapper.Begin()
..
tx.Get(..).Scan(..)
tx.Select(..)
tx.Insert(..)
tx.Update(..)
tx.Delete(..)
..
tx.Commint()/tx.Rollback()
```

## Excution log
```
mapper.SetLogger(log.Printf)
```
logs will be printed like this:

```
..
2015/02/11 13:44:59 [0.00ms] DELETE FROM t WHERE id='1'
2015/02/11 13:44:59 [0.01ms] DELETE FROM t WHERE id='2'
2015/02/11 13:44:59 [17.21ms] SELECT COUNT(*) FROM t
2015/02/11 13:44:59 [0.00ms] INSERT INTO t(first_name, last_name, email_verified, created_at)         VALUES('Lilei', 'LL', 'true', NOW())
2015/02/11 13:44:59 [0.24ms] SELECT COUNT(*) FROM t
..
```

## API
see mapper.go

## Example
see unit_test.go or example/example.go
