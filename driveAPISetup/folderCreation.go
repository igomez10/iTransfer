package driveAPISetup

import drive "google.golang.org/api/drive/v3"
import "log"

// creates a root folder with the given name
func CreateFolder(name string, srv *drive.Service) (string, error) {
	newFolder := drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}
	id, err := srv.Files.Create(&newFolder).Do()
	if err != nil {
		log.Println("Unable to create folder", err)
		return "", err
	}

	return id.Id, err

}
