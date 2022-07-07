package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"context"
)

func main() {
	metadataHost := "metadata.google.internal"
	if os.Getenv("GCE_METADATA_HOST") != "" {
		metadataHost = os.Getenv("GCE_METADATA_HOST")
	}
	url := fmt.Sprintf("http://%s/computeMetadata/v1/instance/service-accounts/default/token", metadataHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("error in NewRequestWithContext: %v", err)
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req.WithContext(context.Background()))
	if err != nil {
		log.Fatalf("error in Do: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	// Attempt to read body as JSON
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Fatalf("error in json.Unmarshal: %v", err)
	}
	if tokenResponse.AccessToken == "" {
		log.Fatalf("got empty accessToken. Full response: %s", body)
	}
}
