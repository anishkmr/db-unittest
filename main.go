package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/anish-kmr/db/employee"
	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "password"
	hostname = "127.0.0.1:3306"
	dbname   = "test"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)

}

func showEmployees(employees []employee.Employee) {
	fmt.Println("=============================================================")
	for _, emp := range employees {
		fmt.Println(emp)

	}
	fmt.Println("=============================================================")
}
func main() {
	db, err := sql.Open("mysql", dsn("test"))
	if err != nil {
		fmt.Println("Cannot connect to mysql database")
		return
	}

	defer db.Close()

	// employee.SetDB(db)
	for i := 0; i < 10; i++ {
		anish := employee.Employee{Name: "Anish" + strconv.Itoa(i), Email: "anish" + strconv.Itoa(i) + "@gmail.com", Role: "SDE1"}
		anish.Save(db)
	}
	// showEmployees(employee.GetAll(db))
	// showEmployees(employee.FindAll(db, map[t.T]t.T{"name": "Anish40"}))
	// employee.UpdateAll(db, map[t.T]t.T{"id": "24"}, map[t.T]t.T{"name": "Anish Kumar"})
	// employee.DeleteAll(db, map[t.T]t.T{"name": "Anish Kumar"})
	// showEmployees(employee.GetAll(db))

	anish, err := employee.GetByID(db, 34)
	if err == nil {
		fmt.Println("Error getting id ", err)
	}
	anish.Name = "UpdatedName"
	res, _ := anish.Update(db)
	l, _ := res.LastInsertId()
	m, _ := res.RowsAffected()
	fmt.Println(l, m)

}
