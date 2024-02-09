/*	==========================================================================
	Yolov8 dataset
	Filename: newdataset.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/yolov8-dataset/core
	Purpose: Methods to create a new dataset

	=============================================================================
*/

package core

import (
	"github.com/CoderSergiy/golib/tools"
	"os"
)

/****************************************************************************************
 *
 * Function : CreateNewDataset
 *
 * Purpose : Create required files for the new dataset
 *
 *   Input : path string - path on drive where to create the dataset folder with required files
 *
 *  Return : error - error if occur
 */
func CreateNewDataset(path string) error {

	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"dataset/test/images", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"dataset/test/labels", os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"dataset/train/images", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"dataset/train/labels", os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"dataset/valid/images", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"dataset/valid/labels", os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(tools.EnsureSlashInEnd(path)+"dataset/data.yaml", generateDataFileContent(), 0644); err != nil {
		return err
	}

	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"uploaded/images", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"uploaded/labels", os.ModePerm); err != nil {
		return err
	}

	// Create folder for versions
	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"versions", os.ModePerm); err != nil {
		return err
	}

	// Create folder for versions
	if err := os.MkdirAll(tools.EnsureSlashInEnd(path)+"models", os.ModePerm); err != nil {
		return err
	}

	return nil
}

/****************************************************************************************
 *
 * Function : generateDataFileContent
 *
 * Purpose : Generate data.yaml file content
 *
 *   Input : Nothing
 *
 *  Return : []byte - file content
 */
func generateDataFileContent() []byte {
	return []byte("train: ../train/images\n" +
		"val: ../valid/images\n" +
		"test: ../test/images\n" +
		"\n" +
		"nc: 0\n" +
		"names: []\n" +
		"\n")
}
