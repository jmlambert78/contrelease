package main

import (
	"github.com/jmlambert78/contrelease/db"
	"github.com/jmlambert78/contrelease/routes"
	"github.com/jmlambert78/contrelease/views"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()
	e.Renderer = views.Init()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database initialisation
	db.DbMain()
	// Routes for Signup, login, user & role management
	e.GET("/", routes.Home)
	e.GET("/signup", routes.Signup)
	e.GET("/signupfill", routes.SignupFilled)
	e.GET("/checklogin", routes.CheckLogin)
	e.GET("/deleteuser", routes.DeleteUser)
	e.GET("/manageroles", routes.ManageUsersRoles)
	e.GET("/editrole", routes.EditUserRoles)
	e.GET("/validateuserrole", routes.UpdateUserRoles)

	// Routes for Releases management
	e.GET("/addnewrelease", routes.AddNewRelease)
	e.POST("/newrelease", routes.NewRelease)
	e.GET("/getallreleases", routes.GetAllReleases)
	e.GET("/getalldcsreleases", routes.GetAllDCReleases)
	e.GET("/getvalidreleases", routes.GetValidReleases)
	e.GET("/validrm", routes.ValidateRelease)
	e.GET("/getscript", routes.GetScript)

	// Mgt of DCs & URL/Repos default paths
	e.POST("/adddc", routes.AddDc)
	e.GET("/listdcs", routes.ListDcs)
	e.GET("/deletedc", routes.DeleteDc)
	e.GET("/updatedc", routes.GetUpdateDc)
	e.GET("/updatedcput", routes.PutUpdateDc)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
