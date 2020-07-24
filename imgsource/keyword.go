package imgsource

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/kynefuk/goLGTM/processor"
)

// URL is a img web site's url
const (
	URL        = "https://loremflickr.com"
	Width  int = 800
	Height int = 600
)

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

	// save http response body to tmp file.
	tmpFile, err := os.Create(keywordImg.source + ".jpeg")
	if err != nil {
		return fmt.Errorf("failed to create out file. error: %s", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	io.Copy(tmpFile, res.Body)

	file, err := os.Open(tmpFile.Name())
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to open file. error: %s", err)
	}

	//read tmp file to get base image.
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image. error: %s", err)
	}

	msgProcessor := processor.NewProcessor(img, message)
	rgba, err := msgProcessor.AddMessage()
	if err != nil {
		return fmt.Errorf("failed to process message. error: %s", err)
	}

	// save the image with message in local new file.
	outFile, err := os.Create(newImgName(file.Name()))
	if err != nil {
		return fmt.Errorf("failed to create out file. error: %s", err)
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		return fmt.Errorf("failed to create out file. error: %s", err)
	}

	err = b.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush buffer. error: %s", err)
	}

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
