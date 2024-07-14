package internal

import (
	"encoding/json"
	"net/http"
)

func RequestOAuthData(oauthHttp *http.Client, req *http.Request, output any) error {
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
