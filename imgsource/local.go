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

// MaxRatio is a ratio of messaging aria
const MaxRatio float64 = 0.8

// MaxFontSize is max size of font size
const MaxFontSize float64 = 100

// MinFontSize is minimum size of font size
const MinFontSize float64 = 1

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

	// Save that RGBA image to disk.
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
		return fmt.Errorf("failed to create out file. error: %s", err)
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

// ft, err := truetype.Parse(gobold.TTF)
// if err != nil {
// 	return fmt.Errorf("failed to parse font. error: %s", err)
// }

// opt := truetype.Options{
// 	Size:              100,
// 	DPI:               0,
// 	Hinting:           0,
// 	GlyphCacheEntries: 0,
// 	SubPixelsX:        0,
// 	SubPixelsY:        0,
// }

// face := truetype.NewFace(ft, &opt)

// imageWidth := 100
// imageHeight := 100
// textTopMargin := 90
// text := "L"

// dst := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

// dr := &font.Drawer{
// 	Dst:  dst,
// 	Src:  image.White,
// 	Face: face,
// 	Dot:  fixed.Point26_6{},
// }

// dr.Dot.X = (fixed.I(imageWidth) - dr.MeasureString(text)) / 2
// dr.Dot.Y = fixed.I(textTopMargin)

// dr.DrawString(text)

// buf := &bytes.Buffer{}
// err = png.Encode(buf, dst)
// if err != nil {
// 	return fmt.Errorf("failed to encode. error: %s", err)
// }

// out, err := os.Create(newImgName(file.Name()))
// out.Close()
// if err != nil {
// 	return fmt.Errorf("failed to create out file. error: %s", err)
// }
// fmt.Println(buf.Bytes())
// out.Write(buf.Bytes())
