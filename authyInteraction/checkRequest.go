package authyInteraction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type authyStatusResponse struct {
	ApprovalRequest struct {
		AppName         string      `json:"_app_name"`
		AppSerialID     int         `json:"_app_serial_id"`
		AuthyID         int         `json:"_authy_id"`
		ID              string      `json:"_id"`
		UserEmail       string      `json:"_user_email"`
		AppID           string      `json:"app_id"`
		CreatedAt       time.Time   `json:"created_at"`
		Notified        bool        `json:"notified"`
		ProcessedAt     time.Time   `json:"processed_at"`
		SecondsToExpire int         `json:"seconds_to_expire"`
		Status          string      `json:"status"`
		UpdatedAt       time.Time   `json:"updated_at"`
		UserID          string      `json:"user_id"`
		UUID            string      `json:"uuid"`
		HiddenDetails   interface{} `json:"hidden_details"`
		Device          struct {
			City                  interface{} `json:"city"`
			Country               string      `json:"country"`
			IP                    string      `json:"ip"`
			Region                interface{} `json:"region"`
			RegistrationCity      interface{} `json:"registration_city"`
			RegistrationCountry   string      `json:"registration_country"`
			RegistrationIP        string      `json:"registration_ip"`
			RegistrationMethod    string      `json:"registration_method"`
			RegistrationRegion    interface{} `json:"registration_region"`
			OsType                string      `json:"os_type"`
			LastAccountRecoveryAt interface{} `json:"last_account_recovery_at"`
			ID                    int         `json:"id"`
			RegistrationDate      int         `json:"registration_date"`
			LastSyncDate          int         `json:"last_sync_date"`
		} `json:"device"`
	} `json:"approval_request"`
	Success bool `json:"success"`
}

func CheckApprovalRequest(requestID string) int {

	url := fmt.Sprintf("https://api.authy.com/onetouch/json/approval_requests/%s", requestID)

	req, _ := http.NewRequest("GET", url, nil)
	authyAPI := os.Getenv("Authy_API_KEY")
	req.Header.Add("X-Authy-API-Key", authyAPI)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Unable to check approval request status")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	status := parseStatusResponse(body)
	return status
}

func parseStatusResponse(response []byte) int {
	var asr authyStatusResponse
	json.Unmarshal(response, &asr)
	var numToReturn int

	switch status := asr.ApprovalRequest.Status; status {
	case "denied":
		numToReturn = 0
	case "approved":
		numToReturn = 1
	case "expired":
		numToReturn = 2
	default:
		numToReturn = -1
	}
	return numToReturn
}
