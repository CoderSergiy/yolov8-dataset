/*	==========================================================================
	Yolov8 dataset
	Filename: images.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/pages/
	Purpose: File has handlers to work with images

	In file:
		1. ImagesHandler
		2. UploadFilesHandler
	=============================================================================
*/

package pages

import (
	"fmt"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/timelib"
	"github.com/CoderSergiy/golib/tools"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

// Model to pass data to the html template
type ImagesModel struct {
	Title        string
	Menu         string
	ErrorMessage string
	DatasetName  string
}

/****************************************************************************************
 *
 * Function : ImagesHandler
 *
 * Purpose : Render images page
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func ImagesHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ET := timelib.EventTimerConstructor()
	logging.Info_Log("Render Images page")

	// Set response headers
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Pointed all template files to render current page
	parsedPage := template.Must(template.ParseFiles(
		templatePath+"layouts/index.gohtml", // Must to be first in the list
		templatePath+"layouts/logo.gohtml",
		templatePath+"layouts/header.gohtml",
		templatePath+"layouts/notifications.gohtml",
		templatePath+"images/body.gohtml",  // page body
		templatePath+"layouts/menu.gohtml", // menu is using in body, so it shouls be after body.gohtml
		templatePath+"layouts/footer.gohtml"))

	// Get dataset name from the request parameters
	datasetName := p.ByName("datasetname")

	// Initialise model
	model := ImagesModel{Menu: "images"}  // Set active menu button
	model.Title = datasetName + " Images" // Set title of the webpage
	model.DatasetName = datasetName

	// Render the page
	err := parsedPage.Execute(w, &model)
	if err != nil {
		logging.Error_Log("Error render dashboard : '%v'", err)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	logging.Info_Log("Finish render images for '%v' dataset page in %s", datasetName, ET.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : UploadFilesHandler
 *
 * Purpose : Handler for the request to upload images
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func UploadFilesHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ET := timelib.EventTimerConstructor()
	logging.Info_Log("Upload image")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `dataset_image`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("dataset_image")
	if err != nil {
		logging.Error_Log("Error Retrieving the File: '%s'", err)
		return
	}
	defer file.Close()

	// Get dataset name
	datasetName := p.ByName("datasetname")
	fileNamePath := tools.EnsureSlashInEnd(datsetsPath) + datasetName + "/uploaded/images/" + handler.Filename

	storeFile, err := os.Create(fileNamePath)
	if err != nil {
		logging.Error_Log("Error when creating file: '%s' , error: '%s'", fileNamePath, err)
		return
	}
	defer storeFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logging.Error_Log("Error reading sended file: '%s'", err)
		return
	}

	// write this byte array to our temporary file
	storeFile.Write(fileBytes)

	// Response with file status to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	logging.Info_Log("File '%v' wirh size '%v'", handler.Filename, PrintFileSize(handler.Size))
	logging.Info_Log("Successfully finish uploading file request in %s", ET.PrintTimerString())
}

func PrintFileSize(fileSize int64) string {
	if int(fileSize/(1024*1024*1024*1024)) > 0 {
		return fmt.Sprintf("%vTB", fileSize/(1024*1024*1024*1024))
	} else if int(fileSize/(1024*1024*1024)) > 0 {
		return fmt.Sprintf("%vGB", fileSize/(1024*1024*1024))
	} else if int(fileSize/(1024*1024)) > 0 {
		return fmt.Sprintf("%vMB", fileSize/(1024*1024))
	} else if int(fileSize/1024) > 0 {
		return fmt.Sprintf("%vKB", fileSize/1024)
	}

	return fmt.Sprintf("%fBytes", fileSize)
}
