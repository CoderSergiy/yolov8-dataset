/*	==========================================================================
	Yolov8 dataset
	Filename: index.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/pages/
	Purpose: File has handlers for index page

	In file:
		1. IndexHandler
		2. UploadFilesHandler
	=============================================================================
*/

package pages

import (
	"errors"
	"fmt"
	"github.com/CoderSergiy/golib/file"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/timelib"
	"github.com/CoderSergiy/golib/tools"
	"github.com/CoderSergiy/yolov8-dataset/core"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"os"
)

// Model to pass data to the html template
type IndexModel struct {
	Title        string
	Directories  []string
	ErrorMessage string
}

/****************************************************************************************
 *
 * Function : IndexHandler
 *
 * Purpose : Handler for index page request
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 _ httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Info_Log("Get dataset list")

	// Get list of the datasets
	model, err := getDatasets()
	if err != nil {
		logging.Error_Log("Error to get datasets list : '%v'", err)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Render Index page with created model
	renderIndexPage(w, model)
}

/****************************************************************************************
 *
 * Function : DatasetCreationHandler
 *
 * Purpose : Render list of curren dataset
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 _ httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func DatasetCreationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get new dataset name from the request
	newFolder := r.FormValue("dataset")

	if newFolder == "" {
		// New dataset name is empty
		renderIndexWithError(w, "New dataset name is empty")
		return
	}

	logging.Info_Log("Create a new dataset '%v'", newFolder)
	fullPathToNewFolder := tools.EnsureSlashInEnd(datsetsPath) + newFolder

	// Check if folder already exists
	if success, _ := file.IsFolderExists(fullPathToNewFolder); success {
		logging.Error_Log("Folder '%v' already exists", datsetsPath)
		renderIndexWithError(w, fmt.Sprintf("Foler '%v' already exists"))
		return
	}

	// Create a new folder
	if err := os.Mkdir(fullPathToNewFolder, os.ModePerm); err != nil {
		logging.Error_Log("Error occur during folder creation: '%v'", err)
		renderIndexWithError(w, fmt.Sprintf("Cannot create folder for the new dataset '%v'"))
		return
	}

	if err := core.CreateNewDataset(fullPathToNewFolder); err != nil {
		logging.Error_Log("Error occur during creation all files/folders for : '%v'", err)
		renderIndexWithError(w, fmt.Sprintf("Cannot create dataset '%v'"))
		return
	}

	logging.Info_Log("Dataset '%v' created", newFolder)

	// Redirect browser to the new dataset dashboard page
	redirectURL := fmt.Sprintf("/dataset/%v/dashboard", newFolder)

	// Redirect to the dashboard of the new dataset
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

/****************************************************************************************
 *
 * Function : renderIndexWithError
 *
 * Purpose : Wrapper to the index page render with an specified error
 *
 *   Input : w http.ResponseWriter - output value
 *			 message string - error message
 *
 *  Return : Nothing
 */
func renderIndexWithError(w http.ResponseWriter, message string) {
	model := IndexModel{ErrorMessage: message}
	renderIndexPage(w, model)
}

/****************************************************************************************
 *
 * Function : renderIndexPage
 *
 * Purpose : Render index page from templates
 *
 *   Input : w http.ResponseWriter - output value
 *			 model IndexModel - model to render index page
 *
 *  Return : Nothing
 */
func renderIndexPage(w http.ResponseWriter, model IndexModel) {
	ET := timelib.EventTimerConstructor()
	logging.Info_Log("Render Index page")
	model.Title = "Yolov8 Vision"

	// Set response headers
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Combined all template files to render current page
	parsedPage := template.Must(template.ParseFiles(
		templatePath+"layouts/index.gohtml", // Must to be first in the list
		templatePath+"layouts/logo.gohtml",
		templatePath+"layouts/header.gohtml",
		templatePath+"layouts/footer.gohtml",
		templatePath+"layouts/notifications.gohtml",
		templatePath+"landingpage/body.gohtml")) // page definition

	// Render the page
	err := parsedPage.Execute(w, &model)
	if err != nil {
		logging.Error_Log("Error render dashboard : '%v'", err)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	logging.Info_Log("Finish render index page in %s", ET.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : getDatasets
 *
 * Purpose : Render index page from templates
 *
 *   Input : Nothing
 *
 *  Return : IndexModel - model to render the current page
 *			 error - error if occur
 */
func getDatasets() (IndexModel, error) {
	model := IndexModel{Title: "Yolov8 Vision"}

	// Check if folder already exists
	if success, _ := file.IsFolderExists(datsetsPath); !success {
		logging.Error_Log("Folder '%v' is not exists", datsetsPath)
		return model, errors.New("Folder is not exists")
	}

	// Get all folders by the path
	var directories []string
	entries, err := os.ReadDir(datsetsPath)
	if err != nil {
		logging.Error_Log("%v", err)
		return model, errors.New("There is no directories")
	}

	// Make array of folders
	for _, e := range entries {
		directories = append(directories, e.Name())
	}
	model.Directories = directories

	logging.Info_Log("Folders [%v]: '%v'", len(directories), tools.Implode(directories))
	return model, nil
}
