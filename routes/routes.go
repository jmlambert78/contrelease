package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/jmlambert78/contrelease/db"
	"github.com/jmlambert78/contrelease/models"

	"github.com/labstack/echo"
)

func AddNewRelease(c echo.Context) error {
	results := getDcs(c)
	return c.Render(http.StatusOK, "newrelease", struct {
		Valid   bool
		Success bool
		DCs     []models.Datacenter
	}{true, true, results})
}

func NewRelease(c echo.Context) error {
	r := &models.Release{}
	if err := c.Bind(r); err != nil {
		return err
	}
	fmt.Println("r:", r)
	r.InsertDate = time.Now()
	r.ReleaseStatus = false // default value for a new inserted release

	// Check if the ZIP file exists and extract Central & LocalImage URLs
	err, imageName, refImageName := CheckZipURLForImage(r.CentralZipURL)
	fmt.Println("Check the zip content : Err:", err, "Found Localimage: ", imageName, "refImage:", refImageName)
	if err == nil {
		r.CentralImage = refImageName
		r.LocalImage = imageName
	}
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("release")
	// Insert
	valid := true
	if err := col.Insert(&r); err != nil {
		valid = false
		fmt.Println("NOK Not Inserted")
	} else {
		valid = true
		//c.JSON(http.StatusOK, "OK")
		fmt.Println("OK Inserted")
	}

	return c.Render(http.StatusOK, "success", struct {
		Valid   bool
		Success bool
	}{true, valid})
}
func NewReleaseSave(c echo.Context) error {
	r := &models.Release{}
	if err := c.Bind(r); err != nil {
		return err
	}
	r.InsertDate = time.Now()
	r.ReleaseStatus = false // default value for a new inserted release
	// Process the URLs for registries
	if strings.Index(r.CentralImage, "://") < 0 {
		r.CentralImage = "https://" + r.CentralImage
	}
	uc, err := url.Parse(r.CentralImage)
	if err != nil {
		panic(err)
	}
	r.CentralImage = uc.Host + uc.Path // update and remove header name.
	if strings.Index(r.LocalImage, "://") < 0 {
		r.LocalImage = "https://" + r.LocalImage
	}
	ul, err := url.Parse(r.LocalImage)
	if err != nil {
		panic(err)
	}
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("dcs")
	if len(ul.Path) <= 0 {
		// need to create an url from the default DC URLs
		var results models.Datacenter
		if err := col.Find(bson.M{"name": r.Destination}).One(&results); err != nil {
			fmt.Println("NOK DC Not found")
		} else {
			fmt.Println("OK DC found", results)
			r.LocalImage = results.RegUrl + uc.Path
		}
	}

	col = Db.C("release")
	// Insert
	valid := true
	if err := col.Insert(&r); err != nil {
		valid = false
		fmt.Println("NOK Not Inserted")
	} else {
		valid = true
		//c.JSON(http.StatusOK, "OK")
		fmt.Println("OK Inserted")
	}

	return c.Render(http.StatusOK, "success", struct {
		Valid   bool
		Success bool
	}{true, valid})
}

func GetAllReleases(c echo.Context) error {
	Destination := c.FormValue("dest")
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("release")
	//find based on the criteria of the query

	var results []models.Release
	valid := http.StatusOK
	if err := col.Find(bson.M{"dest": Destination}).All(&results); err != nil {
		valid = http.StatusNoContent
		fmt.Println("NOK releases Not found")
	} else {
		valid = http.StatusOK
		fmt.Println("OK Releases found", results)
	}
	cookie, _ := c.Cookie("username")
	fmt.Println("cookie", cookie.Value, getUserRole(c, cookie.Value))
	return c.Render(valid, "releases", struct {
		Destination string
		Results     []models.Release
		Date        time.Time
		AllReleases string
		Roles       models.Roles
	}{Destination, results, time.Now(), "ALL", getUserRole(c, cookie.Value)})
}
func GetAllDCReleases(c echo.Context) error {
	Destination := c.FormValue("dest")
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("release")
	//find based on the criteria of the query
	var results []models.Release
	valid := http.StatusOK
	if err := col.Find(bson.M{}).All(&results); err != nil {
		valid = http.StatusNoContent
		fmt.Println("NOK releases Not found")
	} else {
		valid = http.StatusOK
		fmt.Println("OK Releases found", results)
	}
	checkAllZipImages(results)
	cookie, _ := c.Cookie("username")
	fmt.Println("cookie", cookie.Value, getUserRole(c, cookie.Value))
	return c.Render(valid, "releases", struct {
		Destination string
		Results     []models.Release
		Date        time.Time
		AllReleases string
		Roles       models.Roles
	}{Destination, results, time.Now(), "ALL", getUserRole(c, cookie.Value)})
}
func GetValidReleases(c echo.Context) error {
	Destination := c.FormValue("dest")
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("release")
	//find based on the criteria of the query

	var results []models.Release
	valid := http.StatusOK
	if err := col.Find(bson.M{"dest": Destination, "rmok": true}).All(&results); err != nil {
		valid = http.StatusNoContent
		fmt.Println("NOK releases Not found")
	} else {
		valid = http.StatusOK
		fmt.Println("OK Releases found", results)

	}
	cookie, _ := c.Cookie("username")
	return c.Render(valid, "releases", struct {
		Destination string
		Results     []models.Release
		Date        time.Time
		AllReleases string
		Roles       models.Roles
	}{Destination, results, time.Now(), "RELEASED", getUserRole(c, cookie.Value)})
}

const ScriptModel string = `#!/bin/bash
# This script is generated for the Datacenter : {{.Destination}} to pull {{len .Results}} releases material by {{.Date}}
#
{{$toto :=.LocalZipURL}}
{{range $key, $release := .Results}}
# Service : {{$release.ServiceName}} Version : {{$release.ServiceVersion}} Date: {{$release.InsertDate}}
# Wget the Zip
wget -P {{$toto}} {{$release.CentralZipURL}}
# Pull central image
docker pull {{$release.CentralImage}}
# Tag & Push image to local repo
docker tag $(docker image -q {{$release.CentralImage}}) {{$release.LocalImage}}
docker push {{$release.LocalImage}}
{{end}}
#End of script`

// Function will generate the sh script to get the zip & images and run in a local dc vm to populate the local repos.
// It get only the validated releases (rmok=true)
//
func GetScript(c echo.Context) error {
	Destination := c.FormValue("dest")
	ViewOption := c.FormValue("option")
	dc := getOneDc(c, Destination)
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("release")
	//find based on the criteria of the query
	var results []models.Release
	valid := http.StatusOK
	if err := col.Find(bson.M{"dest": Destination, "rmok": true}).All(&results); err != nil {
		valid = http.StatusNoContent
		fmt.Println("NOK releases Not found")
	} else {
		valid = http.StatusOK
		fmt.Println("OK Releases found", results)
	}

	tmpl, err := template.New("genscriptfile").Parse(ScriptModel)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, struct {
		Destination string
		Results     []models.Release
		Date        time.Time
		LocalZipURL string
	}{Destination, results, time.Now(), dc.ZipRepoUrl})
	if err != nil {
		panic(err)
	}
	if ViewOption == "viewonly" {
		return c.Blob(valid, "text/*", doc.Bytes())
	} else {
		FileName := "attachment;filename=ContRel-" + Destination + "-" + time.Now().Format("2006 Jan 02 15:04") + ".sh"
		c.Response().Header().Set("Content-Disposition", FileName)
		return c.Blob(valid, "text/csv", doc.Bytes())
	}
}
func ValidateRelease(c echo.Context) error {
	Destination := c.FormValue("dest")
	ReleaseId := c.FormValue("id")
	ServiceName := c.FormValue("sname")
	ServiceVersion := c.FormValue("svers")
	fmt.Println(ReleaseId)
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("release")
	// // Update
	colQuerier := bson.M{"dest": Destination, "sname": ServiceName, "svers": ServiceVersion}
	change := bson.M{"$set": bson.M{"rmok": true}}
	err := col.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}
	return GetAllReleases(c)
}
