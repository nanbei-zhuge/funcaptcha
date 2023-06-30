package funcaptcha

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func toJSON(data interface{}) string {
	str, _ := json.Marshal(data)
	return string(str)
}

func jsonToForm(data string) string {
	// Unmarshal into map
	var form_data map[string]interface{}
	json.Unmarshal([]byte(data), &form_data)
	// Use reflection to convert to form data
	var form url.Values = url.Values{}
	for k, v := range form_data {
		form.Add(k, fmt.Sprintf("%v", v))
	}
	return form.Encode()
}

func DownloadChallenge(urls []string) error {
	for _, url := range urls {
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header = headers
		resp, err := (*client).Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("status code %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		// Figure out filename from URL
		url_paths := strings.Split(url, "/")
		filename := strings.Split(url_paths[len(url_paths)-1], "?")[0]
		err = os.WriteFile(filename, body, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTimeStamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

func getRequestId(sessionId string) string {
	pwd := fmt.Sprintf("REQUESTED%sID", sessionId)
	return encrypt(`{"sc":[147,307]}`, pwd)
}
