package booksearch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const url = "https://vision.googleapis.com/v1/images:annotate"

type VisionRequest struct {
	Requests []ImageRequest `json:"requests"`
}

type ImageRequest struct {
	Image    string // base64 content of image
	ImageURI string // url to image
}

func (ir ImageRequest) MarshalJSON() ([]byte, error) {
	if ir.Image != "" {
		return json.Marshal(map[string]interface{}{
			"image": map[string]interface{}{
				"content": ir.Image,
			},
			"features": map[string]interface{}{
				"type": "WEB_DETECTION",
			},
		})
	}

	return json.Marshal(map[string]interface{}{
		"image": map[string]interface{}{
			"source": map[string]interface{}{
				"imageUri": ir.ImageURI,
			},
		},
		"features": map[string]interface{}{
			"type": "WEB_DETECTION",
		},
	})
}

func GoogleVision(image []byte) ([]string, error) {
	r := VisionRequest{
		Requests: []ImageRequest{ImageRequest{
			Image: base64.StdEncoding.EncodeToString(image),
		}},
	}

	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("key", os.Getenv("API_KEY"))
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	return nil, nil
}

const entityURL = "https://kgsearch.googleapis.com/v1/entities:search"

type entityResponse struct {
	Items []itemListElement `json:"itemListElement"`
}

type itemListElement struct {
	Result result `json:"Result"`
}

type result struct {
	Name        string      `json:"name"`
	Short       string      `json:"description"`
	Description description `json:"detailedDescription"`
}

type description struct {
	Body string `json:"articleBody"`
	URL  string `json:"url"`
}

func entity(eid string) (name, short, description string, err error) {
	req, err := http.NewRequest("GET", entityURL, nil)
	if err != nil {
		return "", "", "", err
	}

	query := req.URL.Query()
	query.Set("key", os.Getenv("API_KEY"))
	query.Set("ids", eid)
	query.Set("languages", "fr")
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	fmt.Println(string(b))

	var e entityResponse
	err = json.Unmarshal(b, &e)
	if err != nil {
		return "", "", "", err
	}

	return e.Items[0].Result.Name,
		e.Items[0].Result.Short,
		e.Items[0].Result.Description.Body,
		nil
}
