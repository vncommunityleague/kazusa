package internal

import (
	"encoding/json"
	"net/http"
)

func RequestOAuthData(oauthHttp *http.Client, url string, output any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Vietnam Community League")
	req.Header.Add("Content-Type", "application/json")

	resp, err := oauthHttp.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return err
	}

	return nil
}
