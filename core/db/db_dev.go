/*
NoFace Anonymous Chat Service 
Database Functionality
No auth or registration, just putting stuff in 
and taking stuff out
*/
package db

import (
	//"bufio"
	//"os"
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/*
Name: getDB
Params: dbType - for example, mysql
		login - [dbUser:dbPass@tcp(dbServer:port/dbName)]
Returns: SQL DB object
*/
func GetDB(dbType, login string) *sql.DB {
	db, err := sql.Open(dbType, login)
	if err != nil {
		fmt.Println("DB Handle Failed")
		log.Fatal(err)
	}
	return db
}

/*
Name: 	addUser
Params:	db - SQL DB object
		id - User ID to add
		disp - Display name for user
		hash - hashed password of user
		mfa - Google Auth token for user
Return: true if user was added successfully
*/
func AddUser(db sql.DB, id int, hash, disp, mfa string) bool {
	// Prepare the statement
	stmt, err := db.Prepare("INSERT INTO faces(id, disp, pw, mfa) VALUES(?,?,?,?)")
	if err != nil {
		fmt.Println("Statement Preparation Failed")
		log.Fatal(err)
		return false
	}

	// Try to execute the statement
	_, err = stmt.Exec(id, disp, hash, mfa)
	if err != nil {
		fmt.Println("Statement Execution Failed.")
		log.Fatal(err)
		return false
	}
	return true
}

/*
Name: 	delUser
Param:	db - SQL Database Object
		id - user id of the user to remove
		hash - hashed password of user to remove
Return:	true if user removed
*/
func DelUser(db sql.DB, id int) bool {
	// Prepare the statement
	stmt, err := db.Prepare("DELETE FROM faces WHERE id=?")
	if err != nil {
		fmt.Println("Statement Preparation Failed.")
		log.Fatal(err)
		return false
	}

	// Attempt removal of user
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Println("Statement Execution Failed.")
		log.Fatal(err)
		return false
	}
	return true
}

/*
Name: 	getHash
Param:	id - user id
Return: password hash as string
*/
func GetHash(db sql.DB, id int) string {
	var hash string //Placeholder for query result

	err := db.QueryRow("SELECT pw FROM faces WHERE id=?", id).Scan(&hash)
	if err != nil {
		fmt.Println("Login Failed.")
		return ""
	}
	return hash
}

/*
Name: 	checkUser
Param:	id - user id
Return: true if user exists
*/
func CheckUser(db sql.DB, id int) int {
	var exists int //Placeholder for query result
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM faces WHERE id=?)", id).Scan(&exists)
	if err != nil {
		fmt.Println("Query failed")
		log.Fatal(err)
	}
	return exists
}
