package driveAPISetup

import drive "google.golang.org/api/drive/v3"

func UpdateFile(fileID string, newContent []byte, srv *drive.Service) {

	newFile := drive.File{}
	newFile.Id = fileID
	// srv.Files.Cre
	// srv.Files.Update(fileID, newFile)
}
