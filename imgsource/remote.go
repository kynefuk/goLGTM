package imgsource

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/kynefuk/goLGTM/processor"
)

// OutFile is dst file name
const OutFile = "out_LGTM_.png"

// RemoteImage deal with remote image from url
type RemoteImage struct {
	source string
}

// AddMessage adds message to image
func (remoteImage *RemoteImage) AddMessage(message string) error {
	u, err := url.Parse(remoteImage.source)
	if err != nil {
		return err
	}
	res, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to request to URL. error: %s", err)
	}
	defer res.Body.Close()

	// save http response body to tmp file.
	tmpFile, err := ioutil.TempFile("", "sample.jpeg")
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
	outFile, err := os.Create(newImgName(OutFile))
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

// NewRemoteImage is a constructor of RemoteImage
func NewRemoteImage(source string) *RemoteImage {
	return &RemoteImage{source: source}
}
