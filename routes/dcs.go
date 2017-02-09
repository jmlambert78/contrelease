package routes

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/jmlambert78/contrelease/db"
	"github.com/jmlambert78/contrelease/models"

	"github.com/labstack/echo"
)

func AddDc(c echo.Context) error {
	dc := &models.Datacenter{}
	if err := c.Bind(dc); err != nil {
		return err
	}
	fmt.Println(dc)
	// TO DO insert in DB
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")
	// Insert
	if err := col.Insert(&dc); err != nil {
		fmt.Println("NOK Not Inserted")
	} else {
		fmt.Println("OK Inserted")
	}
	return ListDcs(c)
}

func getDcs(c echo.Context) []models.Datacenter {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")
	//find based on the criteria of the query
	var results []models.Datacenter
	if err := col.Find(bson.M{}).All(&results); err != nil {
		fmt.Println("NOK dc Not found")
	} else {
		fmt.Println("OK dcs found", results)
	}
	return results
}
func getOneDc(c echo.Context, dc string) models.Datacenter {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")
	//find based on the criteria of the query
	var results models.Datacenter
	if err := col.Find(bson.M{"name": dc}).One(&results); err != nil {
		fmt.Println("NOK dc Not found")
	} else {
		fmt.Println("OK dcs found", results)
	}
	return results
}
func ListDcs(c echo.Context) error {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")
	//find based on the criteria of the query
	fmt.Println("listdcs")
	var results []models.Datacenter
	valid := http.StatusOK
	if err := col.Find(bson.M{}).All(&results); err != nil {
		valid = http.StatusNoContent
		fmt.Println("NOK dc Not found")
	} else {
		valid = http.StatusOK
		fmt.Println("OK dcs found", results)
	}
	return c.Render(valid, "listdcs", struct {
		Results []models.Datacenter
	}{results})
}

func DeleteDc(c echo.Context) error {
	Datacenter := c.FormValue("name")
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")
	// Delete old record (upsert later...)
	if err := col.Remove(bson.M{"name": Datacenter}); err != nil {
		fmt.Println("NOK Not Deleted")
	} else {
		fmt.Println("OK Deleted")
	}
	return ListDcs(c)
}

func PutUpdateDc(c echo.Context) error {
	Datacenter := c.FormValue("name")
	fmt.Println("entry putupdatedc")
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")

	fmt.Println("regurl", c.FormValue("regurl"), "ziprepo", c.FormValue("ziprepo"))
	// // Update
	colQuerier := bson.M{"name": Datacenter}
	change := bson.M{"$set": bson.M{"regurl": c.FormValue("regurl"), "ziprepo": c.FormValue("ziprepo")}}
	err := col.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	return ListDcs(c)
}
func GetUpdateDc(c echo.Context) error {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")
	//find based on the criteria of the query
	var result models.Datacenter
	valid := http.StatusOK
	if err := col.Find(bson.M{"name": c.FormValue("name")}).One(&result); err != nil {
		valid = http.StatusNoContent
		fmt.Println("NOK dc Not found")
	} else {
		valid = http.StatusOK
		fmt.Println("OK dcs found", result)
	}

	return c.Render(valid, "updatedc", struct {
		Result models.Datacenter
	}{result})
}
