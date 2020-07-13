package imgsource

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

// URL is a img web site's url
const URL = "https://loremflickr.com"

// Width is a width of img
const Width int = 800

// Height is a height of img
const Height int = 600

// KeywordImage deal with remote image coresspond with keyword
type KeywordImage struct {
	source string
}

// AddMessage adds message to image
func (keywordImg *KeywordImage) AddMessage(message string) error {
	u, err := keywordImg.generateURL()
	if err != nil {
		return fmt.Errorf("failed to parse URL. error: %s", err)
	}

	res, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to request to URL. error: %s", err)
	}

	defer res.Body.Close()

	file, err := os.Create(keywordImg.source + ".jpg")
	if err != nil {
		return fmt.Errorf("failed to create file. error: %s", err)
	}
	defer file.Close()

	io.Copy(file, res.Body)

	return nil
}

func (keywordImg *KeywordImage) generateURL() (*url.URL, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	width := strconv.Itoa(Width)
	height := strconv.Itoa(Height)
	u.Path = path.Join(u.Path, width, height, keywordImg.source)

	return u, nil
}

// NewKeywordImage is a constructor of KeywordImage
func NewKeywordImage(source string) *KeywordImage {
	return &KeywordImage{source: source}
}
