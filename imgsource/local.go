package imgsource

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/kynefuk/goLGTM/processor"
)

// LocalImage deal with local image file
type LocalImage struct {
	source string
}

// AddMessage adds message to image
func (localImg *LocalImage) AddMessage(message string) error {

	file, err := os.Open(localImg.source)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to open file. error: %s", err)
	}

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

func newImgName(source string) string {
	filePath := filepath.Dir(source)
	fileName := filepath.Base(source)
	fileExt := filepath.Ext(source)
	fileName = strings.Replace(fileName, fileExt, "", 1)
	return filePath + "/" + fileName + "_LGTM_" + fileExt
}

// NewLocalImage is a constructor of LocalImage
func NewLocalImage(source string) *LocalImage {
	return &LocalImage{source: source}
}
