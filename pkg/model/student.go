package model

import (
	"fmt"

	"github.com/WANNA959/postgres-writer-operator/pkg/postgre"
)

type Student struct {
	Id         int64
	Name       string
	Age        int32
	Department string
}

func (st *Student) TableName() string {
	return "Student"
}

// Insert inserts a row into the DB to which the receiver PostgresDBClient points
func (student *Student) Insert(st *Student) error {
	pc := postgre.GetDbConnection()
	insertQuery := fmt.Sprintf("insert into \"%s\"(id, name, age, department) VALUES(%d, '%s', %d, '%s');", student.TableName(), st.Id, st.Name, st.Age, st.Department)
	if _, err := pc.Exec(insertQuery); err != nil {
		return err
	}
	return nil
}

// Delete deletes row from the DB to which the receiver PostgresDBClient points
func (student *Student) Delete(id int64) error {
	pc := postgre.GetDbConnection()
	deleteQuery := fmt.Sprintf("DELETE FROM '%s' WHERE id=%d;", student.TableName(), id)
	if _, err := pc.Exec(deleteQuery); err != nil {
		return err
	}
	return nil
}
