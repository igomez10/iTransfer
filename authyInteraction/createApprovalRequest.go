package authyInteraction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type AuthyApprovalResponse struct {
	ApprovalRequest struct {
		UUID string `json:"uuid"`
	} `json:"approval_request"`
	Success bool `json:"success"`
}

func CreateApprovalRequest(userID string, seconds int, message string) string {

	url := fmt.Sprintf("https://api.authy.com/onetouch/json/users/%s/approval_requests", userID)
	content := fmt.Sprintf("{\"message\": \"%s\",\"seconds_to_expire\": %d}", message, seconds)
	payload := strings.NewReader(content)

	req, _ := http.NewRequest("POST", url, payload)
	authyAPI := os.Getenv("Authy_API_KEY")
	req.Header.Add("X-Authy-API-Key", authyAPI)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Unable to make approval request http request")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	uuid := parseResponse(body)
	return uuid
}

func parseResponse(response []byte) string {
	var aar AuthyApprovalResponse
	json.Unmarshal(response, &aar)

	return aar.ApprovalRequest.UUID

}
