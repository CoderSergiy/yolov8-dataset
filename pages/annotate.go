/*	==========================================================================
	Yolov8 dataset
	Filename: images.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/pages/
	Purpose: File has handlers to work with images

	In file:
		1. ImagesHandler
		2. UploadFilesHandler
		3. UploadedHandler

	Links:
		1. /dataset/:datasetname/images
		2. /dataset/:datasetname/uploaded
		3. /dataset/:datasetname/uploaded/:page
		4. /dataset/:datasetname/images/annotated
		5. /dataset/:datasetname/images/annotated/:page
	=============================================================================
*/

package pages

import (
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/timelib"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

// Model to pass data to the html template
type AnnotateModel struct {
	Title        string
	Menu         string
	Tag          string
	ErrorMessage string
	DatasetName  string

	UploadedPage int64
	UploadedImgs []string

	Pagination PaginationModel
}

/****************************************************************************************
 *
 * Function : AnnotateHandler
 *
 * Purpose : Handle request to render annotate page
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func AnnotateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ET := timelib.EventTimerConstructor()
	logging.Info_Log("Render Annotate page")

	// Check if dataset folder existing
	if !isDatasetExist(w, r, p) {
		return
	}

	// Set model
	model := AnnotateModel{Tag: "upload"}
	model.UploadedPage = getRequestedPage(p)

	// Render the images page
	RenderAnnotatePage(w, r, p, model)

	logging.Info_Log("Finish render images for '%v' dataset page in %s", p.ByName("datasetname"), ET.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : RenderAnnotatePage
 *
 * Purpose : Render annotate page
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 * 			 model ImagesModel - model to render template
 *
 *  Return : Nothing
 */
func RenderAnnotatePage(w http.ResponseWriter, r *http.Request, p httprouter.Params, model AnnotateModel) {
	// Get dataset name from the request parameters
	datasetName := p.ByName("datasetname")

	// Set response headers
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Pointed all template files to render current page
	parsedPage, errTemplate :=
		template.New("index.gohtml").Funcs(funcPaginationMap).ParseFiles(
			templatePath+"layouts/index.gohtml", // Must to be first in the list
			templatePath+"layouts/logo.gohtml",
			templatePath+"layouts/header.gohtml",
			templatePath+"layouts/notifications.gohtml",
			templatePath+"annotate/body.gohtml", // page body
			templatePath+"layouts/pagination.gohtml",
			templatePath+"layouts/menu.gohtml", // menu is using in body, so it shouls be after body.gohtml
			templatePath+"layouts/footer.gohtml")

	if errTemplate != nil {
		logging.Error_Log("Error parse the files : '%v'", errTemplate)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Initialise model
	model.Menu = "annotate"                 // Set active menu button
	model.Title = datasetName + " Annotate" // Set title of the webpage
	model.DatasetName = datasetName

	// Render the page
	err := parsedPage.Execute(w, &model) //ExecuteTemplate(w, templatePath+"layouts/index.gohtml", &model)//
	if err != nil {
		logging.Error_Log("Error render annotate page : '%v'", err)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}
