// Author:  fengguangpu
// Date:    2014/09/24

// edit the next line to change package name
package main

// Define sqlmapper in the raw string format, instead of xml file to make app deployment easier.
// Read the following rules before adding element.
//
// Rules:
// The ROOT node name MUST be "sqlmap", DO Not change!
// Only four types of SQL node is supported: select/insert/update/delete, others will be ignored.
// The type specified by the SQL node MUST match the SQL statement, parser will find mismatch and report error.
//
// Passed in variables must defined in the format of "#{VarName}", passed in struct should have field named "VarName".
// Receiver struct MUST have the fields whose names are the "UpperCamel" case of the target table fields.
//  For example:
//      SELECT id, phone_number, email_addr as email FROM t WHERE name=#{Name}
//  Passed in struct must have a field named "Name"
//  Receiver struct must have fields named "Id", "PhoneNmuber" and "Email"
//
// Try to avoid using back quote(`) in SQL statements, as it is forbiddend in raw string in Golang.
// Split the previous raw string into 2 and add contents in the middle if back quote must be used.
//  For example:
//      `...
//      SELECT `key` FROM `t` ...`
//  Should be converted as:
//      `...
//      SELECT ` + "`key`" + ` FROM ` + "`t`" + `...`
//
//  Change name of the const raw string as you like.

const sqls =`
<?xml version="1.0" encoding="utf-8"?>
<sqlmap>
    <select id="selectAll">
        select id, first_name, last_name, email_verified, created_at from t
    </select>
    <select id="selectAllById">
        SELECT * FROM ` + "`t`" +` WHERE id=#{Id}
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
`
