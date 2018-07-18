package main

import (
	"fmt"
	"log"

	api "iTransfer/driveAPISetup"

	"google.golang.org/api/drive/v3"
)

func main() {
	srv := api.SetupCloud("client_secret.json")
	respList, err := srv.Files.List().
		Corpora("user").
		IncludeTeamDriveItems(false).
		PageSize(100).
		Q("name = 'iTransfer' and mimeType = 'application/vnd.google-apps.folder'").
		Spaces("drive").
		Do()

	if err != nil {
		fmt.Println("UNEXPECTED ERROR", err)
	} else {
		log.Println("SUCCESSFULLY RECEIVED RESPLIST")
	}

	if len(respList.Files) == 0 {
		// HANDLE FOLDER CREATION
		log.Println("CREATING FOLDER iTransfer")
		api.CreateFolder("iTransfer", srv)
	} else {
		log.Println("SUCESS, FOLDER ALREADY EXISTS")
	}

	iTransferFolder := respList.Files[0]

	newFile := drive.File{}
	newFile.Name = "ALOHA"
	newFile.Parents = []string{iTransferFolder.Id}
	newID, err := srv.Files.Create(&newFile).Do()
	if err != nil {
		fmt.Println("ERROR EXECUTING DO TO CREATE FILE:", err)
	} else {
		fmt.Println("SUCCESS:", newID.Id)
	}
}
