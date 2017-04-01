package models

import "testing"

import (
	"fmt"

	"database/sql"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", "root", "sandman123", "localhost:3306"))
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", "sandman_test")
	if err != nil {
		panic(err)
	}
	if !rows.Next() {
		db.Exec(fmt.Sprintf("CREATE DATABASE %s", "sandman_test"))
	}

	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", "sandman123", "localhost:3306", "sandman_test"))

}

func TestCreateDB(t *testing.T) {
	// Database alias.
	name := "default"
	// Drop table and re-create.
	force := true
	// Print log.
	verbose := true
	// Error.
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		panic(err)
	}
}

func TestAddAccount(t *testing.T) {
	model := NewAccountModel()

	acc := &Account{
		Id:   1,
		Name: "teste1",
	}

	_, err := model.Add(acc)
	if err != nil {
		t.Fail()
		panic(err)
	}

}

func TestLoadAccountById(t *testing.T) {
	model := NewAccountModel()

	v, err := model.GetById(1)
	if err != nil {
		t.Fail()
		panic(err)
	}

	acc := v.(*Account)

	if acc.Id != 1 {
		t.Fail()
	}

}

func TestLoadAllAccount(t *testing.T) {
	model := NewAccountModel()

	acc := &Account{
		Name: "teste2",
	}

	_, err := model.Add(acc)
	if err != nil {
		t.Fail()
		panic(err)
	}

	list, err := model.GetAll(make(map[string]string), []string{}, []string{}, []string{}, 0, 0)
	if err != nil {
		t.Fail()
		panic(err)
	}

	if len(list) != 2 {
		t.Fail()
	}

}

func TestUpdateAccountById(t *testing.T) {
	model := NewAccountModel()

	acc := &Account{
		Id:   1,
		Name: "novo nome",
	}

	err := model.UpdateById(acc)
	if err != nil {
		t.Fail()
		panic(err)
	}

	v, err := model.GetById(1)
	if err != nil {
		t.Fail()
		panic(err)
	}

	acc = v.(*Account)

	if acc.Name != "novo nome" {
		t.Fail()
	}

}

func TestDeleteAccountById(t *testing.T) {
	model := NewAccountModel()

	err := model.Delete(1)
	if err != nil {
		t.Fail()
		panic(err)
	}
}

func TestNotFoundAccountById(t *testing.T) {
	model := NewAccountModel()

	_, err := model.GetById(1)
	if err != orm.ErrNoRows {
		t.Fail()
		panic(err)
	}
}
