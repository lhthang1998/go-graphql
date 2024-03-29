package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mgo-gin/app/api"
	"mgo-gin/db"
	my_graphql "mgo-gin/graphql"
	"mgo-gin/middlewares"
	"os"
)

type Routes struct {
}

type reqBody struct {
	Query string `json:"query"`
}

func (app Routes) StartGin() {
	r := gin.Default()
	publicRoute := r.Group("/api/v1")
	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	r.GET("swagger/*any", middlewares.NewSwagger())

	r.Static("/template/css", "./template/css")
	r.Static("/template/images", "./template/images")
	//r.Static("/template", "./template")

	r.NoRoute(func(context *gin.Context) {
		//context.File("./template/route_not_found.html")
		context.File("./template/index.html")
	})
	api.ApplyToDoAPI(publicRoute, resource)
	api.ApplyUserAPI(publicRoute, resource)

	graphqlRoute := r.Group("")
	schema := my_graphql.Init()
	api.ApplyGraphAPI(graphqlRoute, &schema)

	r.Run(":" + os.Getenv("PORT"))
}
