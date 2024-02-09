/*	==========================================================================
	Yolov8 dataset
	Filename: utils.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/pages
	Purpose: File has handler to render dashboard webpage

	=============================================================================
*/

package pages

import (
	"errors"
	"github.com/CoderSergiy/golib/file"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/tools"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

/****************************************************************************************
 *
 * Function : isDatasetExist
 *
 * Purpose : Render images page
 *
 *   Input : w http.ResponseWriter - output value
 *			 r *http.Request - request detials
 *			 p httprouter.Params - parameter request
 *
 *  Return : Nothing
 */

func isDatasetExist(w http.ResponseWriter, r *http.Request, p httprouter.Params) bool {
	// Get dataset name from the request parameters
	datasetName := p.ByName("datasetname")
	fullPathToNewFolder := tools.EnsureSlashInEnd(datsetsPath) + datasetName

	// Firstly, check if folder exists
	if success, _ := file.IsFolderExists(fullPathToNewFolder); !success {
		logging.Error_Log("Folder '%v' not exists", fullPathToNewFolder)
		// Redirect to the index again
		http.Redirect(w, r, "/?errorMessage=Dataset%20'"+datasetName+"'%20folder%20not%20exist", http.StatusSeeOther)
		return false
	}

	return true
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
func getFilesByPath(path string, dataset string, showPerPage int64, page int64) ([]string, int64, error) {
	logging.Info_Log("getFilesByPath")
	var filesToPrint []string

	// Check if folder exists
	if success, _ := file.IsFolderExists(path); !success {
		logging.Error_Log("Folder '%v' is not exists", path)
		return filesToPrint, -1, errors.New("Folder is not exists")
	}

	logging.Info_Log("Path '%v'. Pagination page: %v", path, page)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return filesToPrint, -1, err
	}

	startImgIndex := (page - 1) * showPerPage
	lastImgIndex := startImgIndex + showPerPage
	for index, f := range files {
		if int64(index) >= startImgIndex && int64(index) < lastImgIndex {
			//filesToPrint = append(filesToPrint, "/dataset/"+dataset+"/download/"+f.Name())
			filesToPrint = append(filesToPrint, f.Name())
		}
	}

	logging.Info_Log("Find [%v] files between index [%v] and [%v] from total [%v]", len(filesToPrint), startImgIndex, lastImgIndex, len(files))
	return filesToPrint, int64(len(files)), nil
}

func RedirectToPage(w http.ResponseWriter, r *http.Request, p httprouter.Params, path string, errorMessage string) {
	pageToRedirect := "/dataset/"+p.ByName("datasetname")+"/uploaded/1"
	logging.Info_Log("Redirect to '%v' as result of '%v'", pageToRedirect, errorMessage)
	// Redirect to the index again
	http.Redirect(w, r, pageToRedirect, http.StatusSeeOther)
}