package main

import (
	"database/sql"
	"log"
	"strconv"
	"testing"
)

func TestGetInfo(t *testing.T) {
	var company string
	db, err:= sql.Open("mysql", "root:211user504UI@/productdb")
	if err != nil {
		log.Println(err)
	}
	company="Apple"
	str := "SELECT * FROM productdb.Products WHERE company = " + "'" + company + "'" +";"
	res := db.QueryRow(str)
	prod := Product{}
	if res != nil {
		err := res.Scan(&prod.Id, &prod.Model, &prod.Company, &prod.Price)
		s := prod.Model + ", Price==" + strconv.Itoa(prod.Price)
		if s=="" {t.Error("Wrong data. Empty list")}
		if err==nil{}
		if prod.Price==0{t.Error("Wrong price")}
		if prod.Model==""{t.Error("Wrong model")}
	}
	if res==nil{t.Error("Nil sql")}
}
