/*	==========================================================================
	Yolov8 dataset
	Filename: server.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/
	Purpose: Server implementation to render project webpages
			 Includes assets and favicon

	Handlers supported by server:
		1. messageStatusHandler
		2. triggerActionHandler
		3. wsChargerHandler
	=============================================================================
*/

package main

import (
	"github.com/CoderSergiy/yolov8-dataset/pages"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()

	// Assets files handler
	//router.GET("/favicon.ico", FaviconHandler)
	router.ServeFiles("/assets/*filepath", http.Dir("./web/statics/"))

	// Dataset dashboard
	router.GET("/dataset/:datasetname/dashboard", pages.DashBoardHandler)

	// Dataset classes
	//router.GET("/dataset/:name/classes/list", pages.DashBoardHandler)
	//router.POST("/dataset/:name/classes/update", pages.DashBoardHandler)

	// Images page
	router.GET("/dataset/:datasetname/images", pages.ImagesHandler)
	router.POST("/dataset/:datasetname/upload", pages.UploadFilesHandler)

	// Landing page
	router.GET("/", pages.IndexHandler)
	router.POST("/create/dataset", pages.DatasetCreationHandler)

	// Run server
	log.Fatal(http.ListenAndServe(":8080", router))
}
