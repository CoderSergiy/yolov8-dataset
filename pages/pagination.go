/*	==========================================================================
	Yolov8 dataset
	Filename: pagination.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/pages
	Purpose: Methods to calculate pagination

	=============================================================================
*/

package pages

// Model to collect pagination data for the html template
type PaginationModel struct {
	Page           int64
	ItemsPerPage   int64
	Items          int64
	LastPageNumber int64
	Url            string
}

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

func getPaginationModel(page int64, totalItems int64, itemsPerPage int64, url string) PaginationModel {

	model := PaginationModel{Page: page,
		Items:          totalItems,
		ItemsPerPage:   itemsPerPage,
		LastPageNumber: (totalItems / itemsPerPage),
		Url:            url}

	return model
}
