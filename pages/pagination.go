/*	==========================================================================
	Yolov8 dataset
	Filename: pagination.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/pages
	Purpose: Methods to calculate pagination

	=============================================================================
*/

package pages

import (
	"html/template"
)

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
 * Function : getPaginationModel
 *
 * Purpose : Constructor for the PaginationModel
 *
 *   Input : page int64 - current page,
 *			 totalItems int64 - total items to split in pages
 *			 itemsPerPage int64 - setted items per page
 *			 url string - url to create page button
 *
 *  Return : PaginationModel
 */

func getPaginationModel(page int64, totalItems int64, itemsPerPage int64, url string) PaginationModel {

	model := PaginationModel{Page: page,
		Items:          totalItems,
		ItemsPerPage:   itemsPerPage,
		LastPageNumber: GetPaginationPages(totalItems, itemsPerPage),
		Url:            url}

	return model
}

/****************************************************************************************
 *
 * Function : GetPaginationPages
 *
 * Purpose : Calculate how many pages do we need to split total amount of items
 *
 *   Input : totalItems int64 - total items
 *			 itemsPerPage int64 - setted items per page
 *
 *  Return : int64 - total pages
 */
func GetPaginationPages(totalItems int64, itemsPerPage int64) int64 {

	lastPageNumber := totalItems / itemsPerPage

	if (totalItems % itemsPerPage) != 0 {
		return (lastPageNumber + 1)
	}

	return lastPageNumber
}

/****************************************************************************************
 *
 * Function : funcPaginationMap
 *
 * Purpose : Create map of functions for the pagination
 *
 *   Input : set of methods
 *
 *  Return : template.FuncMap
 */
var funcPaginationMap template.FuncMap = template.FuncMap{
	"minus": func(a, b int64) int64 {
		return a - b
	},
	"add": func(a, b int64) int64 {
		return a + b
	},
	"doPrint": func(a int64, b int64, c int64) bool {
		if (a - b) < c {
			return false
		}
		return true
	},
	"pagesRangeUp": func(page int64, numberPages int64, totalPages int64) []int64 {
		var pages []int64
		for number := page + 1; number < (page + numberPages); number++ {
			if int64(number) < totalPages {
				pages = append(pages, number)
			}
		}
		return pages
	},
	"pagesRangeDown": func(page int64, numberPages int64) []int64 {
		var pages []int64
		for number := page - numberPages; number < page; number++ {
			if int64(number) > 1 {
				pages = append(pages, number)
			}
		}
		return pages
	},
}
