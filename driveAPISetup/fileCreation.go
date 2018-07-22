package driveAPISetup

import (
	"bytes"
	"fmt"

	drive "google.golang.org/api/drive/v3"
)

// creates a file with the given name under the give parent folder
func CreateFile(name string, parents []string, media []byte, srv *drive.Service) (string, error) {

	newFile := drive.File{}
	newFile.Name = name
	newFile.Parents = parents
	reader := bytes.NewReader(media)
	newID, err := srv.Files.Create(&newFile).Media(reader).Do()
	if err != nil {
		fmt.Println("ERROR EXECUTING DO TO CREATE FILE:", err)
	} else {
		fmt.Println("SUCCESS:", newID.Id)
	}
	return newID.Id, nil
}
