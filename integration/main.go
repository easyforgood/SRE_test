package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Test Failed. ", err)
	}
}
func checkCount(rs *sql.Rows, n int, info string) {

	if rs.Next() {
		var c int
		err := rs.Scan(&c)
		checkErr(err)
		if c != n {
			log.Fatalln("Test Failed. ", info)
		}

	}
}

func main() {
	// initial
	db, err := sql.Open("mysql", "root@tcp(localhost:4000)/test")
	defer db.Close()
	checkErr(err)

	_, err = db.Exec("drop table  if exists x")
	checkErr(err)
	_, err = db.Exec("create table x (id int primary key, c int)")
	checkErr(err)
	_, err = db.Exec(`insert into x values(1, 1);`)
	checkErr(err)

	// begin test
	tx, err := db.Begin()
	checkErr(err)

	rs, err := tx.Query(`select count(*) from x `)
	checkErr(err)
	checkCount(rs, 1, "Transaction Normal Read Error.")
	rs.Close()

	_, err = db.Exec(`insert into x values(2, 2);`)
	checkErr(err)
	rs, err = db.Query(`select count(*) from x `)
	checkErr(err)
	checkCount(rs, 2, "Normal Read Error")
	rs.Close()

	rs, err = tx.Query(`select count(*) from x`)
	checkErr(err)
	checkCount(rs, 1, "Phantom Read When Select")
	rs.Close()

	result, err := tx.Exec("update x set c=3 where id=2")
	rowCnt, err := result.RowsAffected()
	checkErr(err)
	if rowCnt >= 1 {
		log.Fatalln("Test Failed. Phantom Read When Update")
	}

	result, err = tx.Exec("delete from x where id=2")
	rowCnt, err = result.RowsAffected()
	checkErr(err)
	if rowCnt >= 1 {
		log.Fatalln("Test Failed. Phantom Read When Delete")
	}

	tx.Commit()
	log.Println("Test Pass")
	return
}
