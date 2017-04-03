package tests

import "testing"

import (
	"fmt"

	"database/sql"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	config := map[string]string{
		"user":     "root",
		"password": "123456",
		"url":      "localhost:3306",
		"db":       "goshua_test",
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", config["user"], config["password"], config["url"]))
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", config["db"])
	if err != nil {
		panic(err)
	}
	if !rows.Next() {
		db.Exec(fmt.Sprintf("CREATE DATABASE %s", config["db"]))
	}

	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config["user"], config["password"], config["url"], config["db"]))

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
