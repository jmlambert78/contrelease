package main

import (
	"github.com/jmlambert78/ContinuousRelease/db"
	"github.com/jmlambert78/ContinuousRelease/routes"
	"github.com/jmlambert78/ContinuousRelease/views"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	//proxyUrl, _ := url.Parse("http://localhost:5000")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	e := echo.New()
	e.Renderer = views.Init()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database initialisation
	db.DbMain()
	// Routes
	e.GET("/", routes.Home)
	//e.GET("/login", routes.Login)
	e.GET("/signup", routes.Signup)
	e.GET("/signupfill", routes.SignupFilled)
	e.GET("/checklogin", routes.CheckLogin)
	e.GET("/deleteuser", routes.DeleteUser)

	e.GET("/addnewrelease", routes.AddNewRelease)
	e.POST("/newrelease", routes.NewRelease)
	e.GET("/getallreleases", routes.GetAllReleases)
	e.GET("/getalldcsreleases", routes.GetAllDCReleases)

	e.GET("/getvalidreleases", routes.GetValidReleases)
	e.GET("/validrm", routes.ValidateRelease)
	e.GET("/getscript", routes.GetScript)

	// Mgt of Orgs & URL/Repos default paths
	e.POST("/adddc", routes.AddDc)
	e.GET("/listdcs", routes.ListDcs)
	e.GET("/deletedc", routes.DeleteDc)
	e.GET("/updatedc", routes.GetUpdateDc)
	e.GET("/updatedcput", routes.PutUpdateDc)

	e.GET("/manageroles", routes.ManageUsersRoles)
	e.GET("/editrole", routes.EditUserRoles)
	e.GET("/validateuserrole", routes.UpdateUserRoles)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
