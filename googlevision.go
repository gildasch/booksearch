package booksearch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
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

type googleVisionResponse struct {
	Responses []response `json:"responses"`
}

type response struct {
	Detection entities `json:"webDetection"`
}

type entities struct {
	Entities []webEntity `json:"webEntities"`
}

type webEntity struct {
	ID          string  `json:"entityId"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}

func GoogleVision(image []byte) (name, short, description string, err error) {
	r := VisionRequest{
		Requests: []ImageRequest{ImageRequest{
			Image: base64.StdEncoding.EncodeToString(image),
		}},
	}

	j, err := json.Marshal(r)
	if err != nil {
		return "", "", "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(j))
	if err != nil {
		return "", "", "", err
	}

	query := req.URL.Query()
	query.Set("key", os.Getenv("API_KEY"))
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

	var g googleVisionResponse
	err = json.Unmarshal(b, &g)
	if err != nil {
		return "", "", "", err
	}

	if len(g.Responses) == 0 {
		return "", "", "", errors.Errorf("no response returned")
	}

	for _, e := range g.Responses[0].Detection.Entities {
		name, short, desc, err := entity(e.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		return name, short, desc, nil
	}

	return "", "", "", nil
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
	Types       []string    `json:"@type"`
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

	if len(e.Items) == 0 {
		return "", "", "", errors.Errorf("entity %q not found", eid)
	}

	isBook := false
	for _, t := range e.Items[0].Result.Types {
		if t == "Book" {
			isBook = true
			break
		}
	}

	if !isBook {
		return "", "", "",
			errors.Errorf("entity %q (types %v) is not a book", eid, e.Items[0].Result.Types)
	}

	return e.Items[0].Result.Name,
		e.Items[0].Result.Short,
		e.Items[0].Result.Description.Body,
		nil
}
