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
	"fmt"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/timelib"
	"github.com/CoderSergiy/golib/tools"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// Model to pass data to the html template
type ImagesModel struct {
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

	// Check if dataset folder existing
	if !isDatasetExist(w, r, p) {
		return
	}

	// Set model
	model := ImagesModel{Tag: "upload"}
	model.UploadedPage = getRequestedPage(p)

	// Render the images page
	RenderImagesPage(w, r, p, model)

	logging.Info_Log("Finish render images for '%v' dataset page in %s", p.ByName("datasetname"), ET.PrintTimerString())
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

	// Check if dataset folder existing
	if !isDatasetExist(w, r, p) {
		return
	}

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
	fileNamePath := tools.EnsureSlashInEnd(datsetsPath) + p.ByName("datasetname") + "/uploaded/images/" + handler.Filename

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

/****************************************************************************************
 *
 * Function : UploadedHandler
 *
 * Purpose : Handler for the request to upload images
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func UploadedHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ET := timelib.EventTimerConstructor()
	logging.Info_Log("Uploaded image")

	// Check if dataset folder existing
	if !isDatasetExist(w, r, p) {
		return
	}

	// Check if page number makes sense
	if getRequestedPage(p) == -1 {
		// Redirect to the index again
		http.Redirect(w, r, "/dataset/"+p.ByName("datasetname")+"/uploaded/1", http.StatusSeeOther)
		return
	}

	// Get uploaded Images from dataset
	files, totalFiles, err := getFilesByPath(
		tools.EnsureSlashInEnd(datsetsPath)+p.ByName("datasetname")+"/uploaded/images/",
		p.ByName("datasetname"),
		maxImagesInGallery,
		getRequestedPage(p))

	if err != nil {
		// Redirect to the index again
		http.Redirect(w, r, "/?errorMessage=Cannot%20get%20files%20for%20'"+p.ByName("datasetname")+"'", http.StatusSeeOther)
		return
	}

	if getRequestedPage(p) > GetPaginationPages(totalFiles, maxImagesInGallery) {
		// Redirect to the index again
		http.Redirect(w, r, "/dataset/"+p.ByName("datasetname")+"/uploaded/1", http.StatusSeeOther)
		return
	}

	// Set model
	model := ImagesModel{Tag: "uploaded"}
	model.UploadedPage = getRequestedPage(p)
	model.UploadedImgs = files
	model.Pagination = getPaginationModel(getRequestedPage(p), totalFiles, maxImagesInGallery, "/dataset/"+p.ByName("datasetname")+"/uploaded/")

	// Render the images page
	RenderImagesPage(w, r, p, model)

	logging.Info_Log("Uploaded page render for datatset '%v' and page %v", p.ByName("datasetname"), model.UploadedPage)
	logging.Info_Log("Successfully finish uploading file request in %s", ET.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : RenderImagesPage
 *
 * Purpose : Render images page
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 * 			 model ImagesModel - model to render template
 *
 *  Return : Nothing
 */
func RenderImagesPage(w http.ResponseWriter, r *http.Request, p httprouter.Params, model ImagesModel) {
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
			templatePath+"images/body.gohtml", // page body
			templatePath+"images/upload.gohtml",
			templatePath+"images/uploaded.gohtml",
			templatePath+"layouts/pagination.gohtml",
			templatePath+"layouts/menu.gohtml", // menu is using in body, so it shouls be after body.gohtml
			templatePath+"layouts/footer.gohtml")

	if errTemplate != nil {
		logging.Error_Log("Error parse the files : '%v'", errTemplate)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Initialise model
	model.Menu = "images"                 // Set active menu button
	model.Title = datasetName + " Images" // Set title of the webpage
	model.DatasetName = datasetName

	// Render the page
	err := parsedPage.Execute(w, &model) //ExecuteTemplate(w, templatePath+"layouts/index.gohtml", &model)//
	if err != nil {
		logging.Error_Log("Error render dashboard : '%v'", err)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}

/****************************************************************************************
 *
 * Function : DownloadImageHandler
 *
 * Purpose : Handler for the 'Image download' request
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 *
 *  Return : Nothing
 */
func DownloadImageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.ServeFile(w, r, tools.EnsureSlashInEnd(datsetsPath)+p.ByName("datasetname")+"/uploaded/images/"+p.ByName("filename"))
}

/****************************************************************************************
 *
 * Function : getRequestedPage
 *
 * Purpose : Get page number from the request
 *
 *   Input : p httprouter.Params - parameter request
 *
 *  Return : int64 - page number
 */
func getRequestedPage(p httprouter.Params) int64 {
	pageInt, err := strconv.ParseInt(p.ByName("page"), 10, 0)
	if err != nil {
		return -1
	}

	return pageInt
}

/****************************************************************************************
 *
 * Function : PrintFileSize
 *
 * Purpose : Print file size in string format
 *
 *   Input : fileSize int64 - file size in int format
 *
 *  Return : String
 */
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
