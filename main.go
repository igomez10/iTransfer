package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	authy "iTransfer/authyInteraction"
	api "iTransfer/driveAPISetup"

	drive "google.golang.org/api/drive/v3"
)

var iTransferFolder *drive.File

type handler struct {
	Srv *drive.Service
}

type confirmationMessage struct {
	Message  string `json:"message"`
	ID       string `json:"id"`
	FileName string `json:"file_name"`
}

func main() {
	srv := api.SetupCloud("client_secret.json")
	h := handler{}
	h.Srv = srv
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
		fmt.Printf("\rSUCCESSFULLY RECEIVED RESPLIST                  ")

	}

	if len(respList.Files) == 0 {
		fmt.Printf("\nCREATING FOLDER iTransfer\n")
		api.CreateFolder("iTransfer", srv)
	} else {
		fmt.Printf("\rSUCESS, FOLDER ALREADY EXISTS                ")
	}

	iTransferFolder = respList.Files[0]

	http.HandleFunc("/health", h.GetHealth)

	http.HandleFunc("/entry", h.PostFile)
	port := os.Getenv("PORT")
	fmt.Printf("\rREADY - LISTENING ON PORT %s      ", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func (h *handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, "OK")
		log.Println(r.Method, r.URL)
	}
}

func (h *handler) PostFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		w.Header().Set("Server", "iTransfer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		body, meta, err := r.FormFile("file")
		if err != nil {
			log.Println("Error reading file from request", err)
		}
		defer body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)

		log.Println(r.Method, r.URL)
		userID := os.Getenv("USERID")
		requestID := authy.CreateApprovalRequest(userID, 120, fmt.Sprintf("Request to upload %s", meta.Filename))
		for authy.CheckApprovalRequest(requestID) == -1 {
		}

		switch authy.CheckApprovalRequest(requestID) {
		case 2:
			w.Write([]byte("{\"error_message\":\"Request expired\"}"))
			return
		case 0:
			w.Write([]byte("{\"error_message\":\"Request denied\"}"))
			return
		}

		uploadedFileID, err := api.CreateFile(
			meta.Filename,
			[]string{iTransferFolder.Id},
			buf.Bytes(),
			h.Srv)

		res, _ := json.Marshal(confirmationMessage{"SUCESS", uploadedFileID, meta.Filename})

		w.Write(res)
	}
}

func init() {
	neededEnvVars := []string{"USERID", "PORT"}

	for i := 0; i < len(neededEnvVars); i++ {
		if os.Getenv(neededEnvVars[i]) == "" {
			fmt.Printf("\nERROR: env variable %s not found \n", neededEnvVars[i])
			os.Exit(1)
		}
	}
}
