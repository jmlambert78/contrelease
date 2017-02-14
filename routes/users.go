package routes

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/jmlambert78/contrelease/db"
	"github.com/jmlambert78/contrelease/libs"
	"github.com/jmlambert78/contrelease/models"

	"github.com/labstack/echo"
)

func Home(c echo.Context) error {
	return c.Render(http.StatusOK, "index", struct {
		Valid   bool
		Success bool
	}{true, false})
}
func Signup(c echo.Context) error {
	results := getDcs(c)
	return c.Render(http.StatusOK, "signup", struct {
		Valid   bool
		Success bool
		DCs     []models.Datacenter
	}{true, true, results})
}
func FillPerso(c echo.Context) error {
	return c.Render(http.StatusOK, "perso", struct {
		Valid   bool
		Success bool
	}{false, false})
}

func SignupFilled(c echo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return err
	}
	fmt.Println(u)
	pass := libs.Password{}
	u.Pass = pass.Gen(u.Pass)
	fmt.Println("after crypt", u)
	// TO DO insert in DB

	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("users")
	// Insert
	valid := true
	if err := col.Insert(&u); err != nil {
		valid = false
		fmt.Println("NOK Not Inserted")
	} else {
		valid = true
		//c.JSON(http.StatusOK, "OK")
		fmt.Println("OK Inserted")
	}

	return c.Render(http.StatusOK, "index", struct {
		Valid   bool
		Success bool
	}{true, valid})
}

func CheckLogin(c echo.Context) error {
	//u := &models.User{}
	fmt.Println(c.FormValue("login"))
	fmt.Println(c.FormValue("password"))

	_pass := string(c.FormValue("password"))
	// Find by login
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("users")
	result := models.User{}
	valid := true
	if err := col.Find(bson.M{"name": c.FormValue("login")}).One(&result); err != nil {
		valid = false
		fmt.Println("NOK User Not found")
	} else {
		valid = true
		fmt.Println("OK User found")
		pass := libs.Password{}
		valid = pass.Compare(result.Pass, _pass)
		fmt.Println("compare pass", valid)
	}

	if valid {
		cookie := new(http.Cookie)
		cookie.Name = "username"
		cookie.Value = result.Name
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)
		fmt.Println("cookie", cookie.Name, cookie.Value)

		roles := getUserRole(c, result.Name)
		return c.Render(http.StatusOK, "cractions", struct {
			Valid    bool
			Success  bool
			Username string
			Roles    models.Roles
		}{valid, true, result.Name, roles})
	} else {
		return Home(c)
	}
}
func getUsers(c echo.Context) []models.User {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("users")
	//find based on the criteria of the query
	var results []models.User
	if err := col.Find(bson.M{}).All(&results); err != nil {
		fmt.Println("NOK dc Not found")
	} else {
		fmt.Println("OK dcs found", results)
	}
	return results
}
func DeleteUser(c echo.Context) error {
	User := c.FormValue("name")
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("users")
	// Delete old record (upsert later...)
	if err := col.Remove(bson.M{"name": User}); err != nil {
		fmt.Println("NOK User Not Deleted")
	} else {
		fmt.Println("OK User Deleted")
	}
	DeleteRole(c, User)
	return ManageUsersRoles(c)
}
func ManageUsersRoles(c echo.Context) error {
	users := getUsers(c)
	roles := getRoles(c)
	return c.Render(http.StatusOK, "selectuser", struct {
		Valid   bool
		Success bool
		Users   []models.User
		Roles   []models.Roles
	}{true, true, users, roles})
}
func EditUserRoles(c echo.Context) error {
	DCs := getDcs(c)
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("roles")
	//find based on the criteria of the query
	var result models.Roles
	valid := http.StatusOK
	if err := col.Find(bson.M{"name": c.FormValue("name")}).One(&result); err != nil {
		valid = http.StatusOK
		fmt.Println("NOK No Roles yet Not found")
		//Initialise default values
		result.UserName = c.FormValue("name")
		result.DC = DCs[0].Name
	} else {
		valid = http.StatusOK
		fmt.Println("OK Roles found found", result)
	}

	return c.Render(valid, "updateuserrole", struct {
		Valid   bool
		Success bool
		Roles   models.Roles
		DCs     []models.Datacenter
	}{true, true, result, DCs})
}
func UpdateUserRoles(c echo.Context) error {
	// Get data from Bind

	role := &models.Roles{}
	role.UserName = c.FormValue("name")
	role.DC = c.FormValue("dc")
	role.DevRole = c.FormValue("devrole") == "yes"
	role.RmRole = c.FormValue("rmrole") == "yes"
	role.AdmRole = c.FormValue("admrole") == "yes"
	role.OpsRole = c.FormValue("opsrole") == "yes"

	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("roles")
	//find based on the criteria of the query
	var result models.Roles
	if err := col.Find(bson.M{"name": c.FormValue("name")}).One(&result); err != nil {
		fmt.Println("NOK No Roles yet Not found")
		// Insert the new roles doc
		if err := col.Insert(&role); err != nil {
			fmt.Println("NOK Role Not Inserted")
		} else {
			fmt.Println("OK Role Inserted")
		}
	} else {
		fmt.Println("OK Roles found found", result)
		// Update
		colQuerier := bson.M{"name": c.FormValue("name")}
		change := bson.M{"$set": bson.M{"name": role.UserName, "dc": role.DC, "devrole": role.DevRole, "rmrole": role.RmRole, "opsrole": role.OpsRole, "admrole": role.AdmRole}}
		err := col.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}
	}

	return c.Render(http.StatusOK, "index", struct {
		Valid   bool
		Success bool
	}{true, true})
}
func getUserRole(c echo.Context, user string) models.Roles {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("roles")
	//find based on the criteria of the query
	var results models.Roles
	if err := col.Find(bson.M{"name": user}).One(&results); err != nil {
		fmt.Println("NOK User Not found")
	} else {
		fmt.Println("OK User found", results)
	}
	return results
}
func getRoles(c echo.Context) []models.Roles {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("roles")
	//find based on the criteria of the query
	var results []models.Roles
	if err := col.Find(bson.M{}).All(&results); err != nil {
		fmt.Println("NOK User Not found")
	} else {
		fmt.Println("OK User found", results)
	}
	return results
}
func DeleteRole(c echo.Context, User string) error {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()
	col := Db.C("roles")
	// Delete old record (upsert later...)
	if err := col.Remove(bson.M{"name": User}); err != nil {
		fmt.Println("NOK User Not Deleted")
		return err
	} else {
		fmt.Println("OK User Deleted")
		return err
	}
}
