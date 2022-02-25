package employee

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetByID(t *testing.T) {
	d, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error("Error Mocking the db", err)
	}
	defer d.Close()
	row := sqlmock.NewRows([]string{"id", "name", "email", "role"})
	tcs := []struct {
		desc   string
		id     int
		output Employee
		err    error
		mockQ  interface{}
	}{
		{
			desc:   "Success Case",
			id:     1,
			output: Employee{id: 1, Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			mockQ:  mock.ExpectQuery("SELECT * FROM EMPLOYEE WHERE ID=?").WithArgs(1).WillReturnRows(row.AddRow(1, "Anish1", "anish1@gmail.com", "SDE1")),
		},
		{
			desc:   "Fail Case",
			id:     2,
			output: Employee{id: 0, Name: "", Email: "", Role: ""},
			err:    sql.ErrNoRows,
			mockQ:  mock.ExpectQuery("SELECT * FROM EMPLOYEE WHERE ID=?").WithArgs(2).WillReturnError(sql.ErrNoRows),
		},
	}
	for i, tc := range tcs {
		fmt.Println("Test Case #", i, tc)
		out, err := GetByID(d, tc.id)

		if err != nil && !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\nExpected Error \n%v, \nGot Error \n%v", tc.err, err)
		}
		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("\nExpected Employee \n%v, \nGot \n%v", tc.output, out)
		}
	}
}

func TestSave(t *testing.T) {
	d, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error("Error Mocking the db", err)
	}
	defer d.Close()

	tcs := []struct {
		desc   string
		emp    Employee
		output sql.Result
		err    error
		mockQ  interface{}
	}{
		{
			desc:   "Success Case",
			emp:    Employee{Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: sqlmock.NewResult(1, 1),
			mockQ:  mock.ExpectPrepare("INSERT INTO EMPLOYEE(name, email, role) VALUES (?,?,?)").ExpectExec().WithArgs("Anish1", "anish1@gmail.com", "SDE1").WillReturnResult(sqlmock.NewResult(1, 1)),
		},
		{
			desc:   "Error Case",
			emp:    Employee{Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: nil,
			err:    errors.New("fail to execute"),
			mockQ:  mock.ExpectPrepare("INSERT INTO EMPLOYEE(name, email, role) VALUES (?,?,?)").ExpectExec().WithArgs("Anish1", "anish1@gmail.com", "SDE1").WillReturnError(errors.New("fail to execute")),
		},
		{
			desc:   "Error Prepare Case",
			emp:    Employee{Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: nil,
			err:    errors.New("fail to prepare"),
			mockQ:  mock.ExpectPrepare("INSERT INTO EMPLOYEE(name, email, role) VALUES (?,?,?)").WillReturnError(errors.New("fail to prepare")),
		},
		{
			desc:   "empty Case",
			emp:    Employee{Name: "", Email: "", Role: ""},
			output: nil,
			err:    errors.New("empty"),
		},
	}
	for i, tc := range tcs {
		fmt.Println("Test Case #", i, tc)
		out, err := tc.emp.Save(d)
		fmt.Println("ERROR is dfsdaf ", err, tc.err)
		if err != nil && !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\nExpected Error \n%v, \nGot Error \n%v", tc.err, err)
		}
		if err == nil {
			gotLastInserted, _ := out.LastInsertId()
			gotRowsAff, _ := out.RowsAffected()
			expLast, _ := tc.output.LastInsertId()
			expAff, _ := tc.output.RowsAffected()
			if !reflect.DeepEqual(gotLastInserted, expLast) || !reflect.DeepEqual(gotRowsAff, expAff) {
				t.Errorf("\nExpected Employee \n%v, \nGot \n%v", tc.output, out)
			}

		}
	}
}

func TestUpdate(t *testing.T) {
	d, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error("Error Mocking the db", err)
	}
	defer d.Close()

	tcs := []struct {
		desc   string
		id     int
		emp    Employee
		output sql.Result
		err    error
		mockQ  interface{}
	}{
		{
			desc:   "Success Case",
			emp:    Employee{id: 1, Name: "Anish1Updated", Email: "anish1@gmail.com", Role: "SDE1"},
			output: sqlmock.NewResult(0, 1),
			mockQ:  mock.ExpectPrepare("UPDATE EMPLOYEE SET name=? , email=?, role=? WHERE id=?").ExpectExec().WithArgs("Anish1Updated", "anish1@gmail.com", "SDE1", 1).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			desc:   "Fail Case",
			emp:    Employee{id: 1, Name: "", Email: "anish1@gmail.com", Role: "SDE1"},
			output: nil,
			err:    errors.New("empty"),
		},
		{
			desc:   "Error Case",
			emp:    Employee{Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: nil,
			err:    sql.ErrNoRows,
			mockQ:  mock.ExpectPrepare("UPDATE EMPLOYEE SET name=? , email=?, role=? WHERE id=?").ExpectExec().WithArgs("Anish1Updated", "anish1@gmail.com", "SDE1", 1).WillReturnError(sql.ErrNoRows),
		},
		{
			desc:   "Error Prepare Case",
			emp:    Employee{Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: nil,
			err:    errors.New("fail to prepare"),
			mockQ:  mock.ExpectPrepare("UPDATE EMPLOYEE SET name=? , email=?, role=? WHERE id=?").WillReturnError(errors.New("fail to prepare")),
		},
	}
	for i, tc := range tcs {
		fmt.Println("Test Case #", i, tc)
		out, err := tc.emp.Update(d)

		if err != nil && !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\nExpected Error \n%v, \nGot Error \n%v", tc.err, err)
		}
		if err == nil {
			gotLastInserted, _ := out.LastInsertId()
			gotRowsAff, _ := out.RowsAffected()
			expLast, _ := tc.output.LastInsertId()
			expAff, _ := tc.output.RowsAffected()
			if !reflect.DeepEqual(gotLastInserted, expLast) || !reflect.DeepEqual(gotRowsAff, expAff) {
				t.Errorf("\nExpected Employee \n%v, \nGot \n%v", tc.output, out)
			}

		}
	}
}

func TestDelete(t *testing.T) {
	d, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error("Error Mocking the db", err)
	}
	defer d.Close()

	tcs := []struct {
		desc   string
		id     int
		emp    Employee
		output sql.Result
		err    error
		mockQ  interface{}
	}{
		{
			desc:   "Success Case",
			id:     1,
			emp:    Employee{id: 1, Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: sqlmock.NewResult(1, 1),
			mockQ:  mock.ExpectPrepare("DELETE FROM EMPLOYEE WHERE id=?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1)),
		},

		{
			desc:   "Error Prepare Case",
			id:     1,
			emp:    Employee{id: 1, Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: nil,
			err:    errors.New("fail to prepare"),
			mockQ:  mock.ExpectPrepare("DELETE FROM EMPLOYEE WHERE id=?").WillReturnError(errors.New("fail to prepare")),
		},
		{
			desc:   "Error Execute Case",
			id:     1,
			emp:    Employee{id: 1, Name: "Anish1", Email: "anish1@gmail.com", Role: "SDE1"},
			output: nil,
			err:    sql.ErrNoRows,
			mockQ:  mock.ExpectPrepare("DELETE FROM EMPLOYEE WHERE id=?").ExpectExec().WithArgs(1).WillReturnError(sql.ErrNoRows),
		},

		// {
		// 	desc:   "Success Case",
		// 	emp:    Employee{id: 0, Name: "", Email: "", Role: ""},
		// 	output: nil,
		// 	err:    sql.ErrNoRows,
		// },
	}
	for i, tc := range tcs {
		fmt.Println("Test Case #", i, tc)
		out, err := tc.emp.Delete(d)

		if err != nil && !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\nExpected Error \n%v, \nGot Error \n%v", tc.err, err)
		}
		if err == nil {
			gotLastInserted, _ := out.LastInsertId()
			gotRowsAff, _ := out.RowsAffected()
			expLast, _ := tc.output.LastInsertId()
			expAff, _ := tc.output.RowsAffected()
			if !reflect.DeepEqual(gotLastInserted, expLast) || !reflect.DeepEqual(gotRowsAff, expAff) {
				t.Errorf("\nExpected Employee \n%v, \nGot \n%v", tc.output, out)
			}

		}
	}
}
