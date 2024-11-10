package server

import (
	"9DotTechnology/trading_platform/constants/common_constants"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"fmt"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

// StartServer - start Server server
func StartServer() {

	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.New()

	/*
		Define middleware
	*/

	// set logger in gin
	router.Use(logwrapper.GinLogger(), gin.Recovery())

	// initializing all the api endpoints and groups
	initializeRoutes()

	// running the gin server
	logwrapper.Logger.Infoln("Running Server on port : ", common_constants.SERVER_PORT)
	logwrapper.Logger.Fatal(router.Run(fmt.Sprintf(":%d", common_constants.SERVER_PORT)))

}
