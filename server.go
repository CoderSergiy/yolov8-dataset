/*	==========================================================================
	Yolov8 dataset
	Filename: server.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/
	Purpose: Server implementation to render project webpages
			 Includes assets and favicon

	=============================================================================
*/

package main

import (
	"github.com/CoderSergiy/yolov8-dataset/pages"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func faviconHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./web/statics/img/favicon.ico")
}

func main() {
	router := httprouter.New()

	// Assets files handler
	router.GET("/favicon.ico", faviconHandler)
	router.ServeFiles("/assets/*filepath", http.Dir("./web/statics/"))

	// Dataset dashboard
	router.GET("/dataset/:datasetname/", pages.DashBoardHandler) // Dashboard index page
	router.GET("/dataset/:datasetname/dashboard", pages.DashBoardHandler)

	// Dataset classes
	//router.GET("/dataset/:name/classes/list", pages.DashBoardHandler)
	//router.POST("/dataset/:name/classes/update", pages.DashBoardHandler)

	// Images page
	router.GET("/dataset/:datasetname/images", pages.ImagesHandler)
	router.GET("/dataset/:datasetname/uploaded", pages.UploadedHandler)
	router.GET("/dataset/:datasetname/uploaded/:page", pages.UploadedHandler)
	router.GET("/dataset/:datasetname/images/annotated/:page", pages.UploadedHandler)
	router.POST("/dataset/:datasetname/upload", pages.UploadFilesHandler)              // Handle 'file upload' request
	router.GET("/dataset/:datasetname/download/:filename", pages.DownloadImageHandler) // Handle 'file download' request

	// Landing page
	router.GET("/", pages.IndexHandler)
	router.POST("/create/dataset", pages.DatasetCreationHandler)

	// Run server
	log.Fatal(http.ListenAndServe(":8080", router))
}
