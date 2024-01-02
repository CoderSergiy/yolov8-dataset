/*	==========================================================================
	Yolov8 dataset
	Filename: dashboard.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/pages
	Purpose: File has handler to render dashboard webpage

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
type DashboardModel struct {
	Title        string
	Menu         string
	ErrorMessage string
	DatasetName  string
}

/****************************************************************************************
 *
 * Function : DashBoardHandler
 *
 * Purpose : Render dashboard page
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func DashBoardHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ET := timelib.EventTimerConstructor()
	logging.Info_Log("Render Dashboard page")

	// Set response headers
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Pointed all template files to render current page
	parsedPage := template.Must(template.ParseFiles(
		templatePath+"layouts/index.gohtml", // Must to be first in the list
		templatePath+"layouts/logo.gohtml",
		templatePath+"layouts/header.gohtml",
		templatePath+"layouts/notifications.gohtml",
		templatePath+"dashboard/body.gohtml", // page body
		templatePath+"layouts/menu.gohtml",   // menu is using in bodu, so it shouls be after body.gohtml
		templatePath+"layouts/footer.gohtml"))

	// Get dataset name from the request parameters
	datasetName := p.ByName("datasetname")

	// Initialise model
	model := DashboardModel{Menu: "dashboard"} // Set active menu button
	model.Title = datasetName + " Dashboard"   // Set title of the webpage
	model.DatasetName = datasetName

	// Render the page
	err := parsedPage.Execute(w, &model)
	if err != nil {
		logging.Error_Log("Error render dashboard : '%v'", err)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	logging.Info_Log("Finish render '%v' dashboard page in %s", datasetName, ET.PrintTimerString())
}
