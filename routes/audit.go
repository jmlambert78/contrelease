package routes

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/jmlambert78/contrelease/db"
	"github.com/jmlambert78/contrelease/models"

	"github.com/labstack/echo"
)

func InsertAudit(who string, action string, rel models.Release, status bool) error {
	audit := &models.AuditTrail{}
	audit.Who = who
	audit.Status = status
	audit.Release = rel
	audit.Action = action
	audit.EventDate = time.Now()

	fmt.Println(audit)
	// TO DO insert in DB
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("audittrails")
	// Insert
	if err := col.Insert(&audit); err != nil {
		fmt.Println("NOK audit Not Inserted")
		return err
	} else {
		fmt.Println("OK audit Inserted")
		return err
	}
}

func GetAuditTrails(c echo.Context) error {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("audittrails")
	//find based on the criteria of the query
	var results []models.AuditTrail
	if err := col.Find(bson.M{}).Sort("-adate").All(&results); err != nil {
		fmt.Println("NOK dc Not found")
	} else {
		fmt.Println("OK dcs found", results)
	}
	return c.Render(http.StatusOK, "audittrails", struct {
		Results []models.AuditTrail
	}{results})
}
