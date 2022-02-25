package employee

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	t "github.com/anish-kmr/db/Type"
)

type T t.T

type Employee struct {
	id    int
	Name  string
	Email string
	Role  string
}

// func SetDB(db *sql.DB) {
// 	res, err := db.Exec(`
// 		CREATE TABLE IF NOT EXISTS EMPLOYEE(
// 			id int primary key AUTO_INCREMENT,
// 			name VARCHAR(20) ,
// 			email VARCHAR(50),
// 			role VARCHAR(20)
// 		)`)
// 	if err != nil {
// 		fmt.Println("[EMPLOYEE] : Error Creating Table ")
// 		return
// 	}
// 	a, _ := res.RowsAffected()
// 	if a == 0 {
// 		fmt.Println("[EMPLOYEE] : Table `EMPLOYEE` Already Exists ")
// 	} else {
// 		fmt.Println("[EMPLOYEE] : Table `EMPLOYEE` Created ")
// 	}
// }

func GetByID(db *sql.DB, id int) (Employee, error) {
	row := db.QueryRow("SELECT * FROM EMPLOYEE WHERE ID=?", id)

	emp := Employee{}

	if err := row.Scan(&emp.id, &emp.Name, &emp.Email, &emp.Role); err != nil {
		fmt.Println("[EMPLOYEE] : No Rows Found")
		return emp, sql.ErrNoRows

	}

	return emp, nil

}

func (e *Employee) Save(db *sql.DB) (sql.Result, error) {

	if e.Name == "" || e.Email == "" || e.Role == "" {
		fmt.Printf("[EMPLOYEE] : empty fields %v %T \n", e, errors.New("empty"))
		return nil, errors.New("empty")
	}
	stmt, err := db.Prepare("INSERT INTO EMPLOYEE(name, email, role) VALUES (?,?,?)")
	if err != nil {
		fmt.Println("[EMPLOYEE] : Error Inserting to `Employee`")
		return nil, errors.New("fail to prepare")
	}
	if err == nil {
		defer stmt.Close()
	}

	res, err := stmt.Exec(e.Name, e.Email, e.Role)
	if err != nil {
		fmt.Println("[EMPLOYEE] : Failed to Save Employee To DB", err)
		return nil, errors.New("fail to execute")
	}

	fmt.Println("[EMPLOYEE] : Employee Saved To DB")
	return res, nil
}

func (e *Employee) Update(db *sql.DB) (sql.Result, error) {
	if e.Name == "" || e.Email == "" || e.Role == "" {
		fmt.Printf("[EMPLOYEE] : empty fields %v %T \n", e, errors.New("empty"))
		return nil, errors.New("empty")
	}
	stmt, err := db.Prepare("UPDATE EMPLOYEE SET name=? , email=?, role=? WHERE id=?")
	if err != nil {
		fmt.Println("[EMPLOYEE] : Error Updating `Employee` ", e)
		return nil, errors.New("fail to prepare")
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.Name, e.Email, e.Role, e.id)
	if err != nil {
		fmt.Println("[EMPLOYEE] : Failed to Update Employee")
		return nil, sql.ErrNoRows

	}
	fmt.Println("[EMPLOYEE] : Employee Updated To DB")
	return res, nil

}

func (e *Employee) Delete(db *sql.DB) (sql.Result, error) {
	stmt, err := db.Prepare("DELETE FROM EMPLOYEE WHERE id=?")
	if err != nil {
		fmt.Println("[EMPLOYEE] : Error Deleting `Employee` ", e)
		return nil, errors.New("fail to prepare")
	}
	if err == nil {
		defer stmt.Close()
	}
	res, err := stmt.Exec(e.id)
	if err != nil {
		fmt.Println("[EMPLOYEE] : Failed to Update Employee")
		return nil, sql.ErrNoRows

	}
	fmt.Println("[EMPLOYEE] : Employee Deleted")
	return res, nil

}

// func GetAll(db *sql.DB) []Employee {
// 	rows, err := db.Query(`
// 		SELECT * FROM EMPLOYEE
// 	`)
// 	if err != nil {
// 		fmt.Println("[EMPLOYEE] : Error fetching from `Employee`")
// 		return nil
// 	}
// 	employees := []Employee{}
// 	for rows.Next() {
// 		emp := Employee{}
// 		if err := rows.Scan(&emp.id, &emp.Name, &emp.Email, &emp.Role); err == nil {
// 			employees = append(employees, emp)
// 		}
// 	}
// 	return employees
// }

// func FindAll(db *sql.DB, condition map[t.T]t.T) []Employee {
// 	c := ""
// 	i := 1
// 	for k, v := range condition {
// 		c += fmt.Sprintf("%v=%q", k, v)
// 		if i != len(condition) {
// 			c += " AND "
// 		}
// 		i += 1
// 	}

// 	qry := fmt.Sprintf("SELECT * FROM EMPLOYEE where %v", c)
// 	rows, err := db.Query(qry)
// 	if err != nil {
// 		fmt.Println("[EMPLOYEE] : Error Fetching `Employee`")
// 		return nil
// 	}
// 	emps := []Employee{}
// 	for rows.Next() {
// 		emp := Employee{}
// 		rows.Scan(&emp.id, &emp.Name, &emp.Email, &emp.Role)
// 		emps = append(emps, emp)
// 	}

// 	return emps
// }

// func UpdateAll(db *sql.DB, condition map[t.T]t.T, updates map[t.T]t.T) int64 {
// 	c := ""
// 	i := 1
// 	for k, v := range condition {
// 		c += fmt.Sprintf("%v=%q", k, v)
// 		if i != len(condition) {
// 			c += " AND "
// 		}
// 		i += 1
// 	}
// 	u := ""
// 	i = 1
// 	for k, v := range updates {
// 		u += fmt.Sprintf("%v=%q", k, v)
// 		if i != len(updates) {
// 			u += ", "
// 		}
// 		i += 1
// 	}
// 	qry := fmt.Sprintf("UPDATE EMPLOYEE SET %v where %v", u, c)
// 	res, err := db.Exec(qry)
// 	if err != nil {
// 		fmt.Println("[EMPLOYEE] : Error Updating `Employee`")
// 		return 0
// 	}
// 	n, _ := res.RowsAffected()
// 	fmt.Println("[EMPLOYEE] : ", n, " Employee Updated Successfully")
// 	return n
// }

// func DeleteAll(db *sql.DB, condition map[t.T]t.T) int64 {
// 	c := ""
// 	i := 1
// 	for k, v := range condition {
// 		c += fmt.Sprintf("%v=%q", k, v)
// 		if i != len(condition) {
// 			c += " AND "
// 		}
// 		i += 1
// 	}

// 	qry := fmt.Sprintf("DELETE FROM EMPLOYEE where %v", c)
// 	fmt.Println("DELETE QRY :", qry)
// 	res, err := db.Exec(qry)
// 	if err != nil {
// 		fmt.Println("[EMPLOYEE] : Error Deleting from `Employee`")
// 		return 0
// 	}
// 	n, _ := res.RowsAffected()
// 	fmt.Println("[EMPLOYEE] : ", n, " Employee Deleted")
// 	return n

// }

func (e Employee) String() string {
	return fmt.Sprintf("%10.50s %10.50s %20.50s %10.50s", strconv.Itoa(e.id), e.Name, e.Email, e.Role)
}
