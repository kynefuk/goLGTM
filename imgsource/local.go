package imgsource

import (
	"bufio"
	"fmt"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
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
	var text = []string{message}

	file, err := os.Open(localImg.source)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to open file. error: %s", err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image. error: %s", err)
	}

	// message aria
	msgAriaWidth := float64(img.Bounds().Dy()) * MaxRatio

	// parse font
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return fmt.Errorf("failed to parse font file. error: %s", err)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))

	var fontSize float64
	var messageXPoint int
	for i := MaxFontSize; i > MinFontSize; i-- {
		opt := truetype.Options{
			Size:              i,
			DPI:               72,
			Hinting:           font.HintingVertical,
			GlyphCacheEntries: 0,
			SubPixelsX:        0,
			SubPixelsY:        0,
		}
		face := truetype.NewFace(ft, &opt)
		dr := &font.Drawer{Dst: rgba, Src: img, Face: face, Dot: fixed.Point26_6{}}
		messageXPoint = ((fixed.I(img.Bounds().Dx()) - dr.MeasureString(message)) / 2).Floor()

		messageAria := int(msgAriaWidth) - dr.MeasureString(message).Floor()
		if messageAria > 0 {
			fontSize = i
			break
		}
	}

	_, fw := image.Black, image.White
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ft)
	c.SetFontSize(float64(fontSize))
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fw)
	c.SetHinting(font.HintingVertical)

	messageYPoint := img.Bounds().Dy() / 2
	// Draw the text.
	pt := freetype.Pt(messageXPoint, messageYPoint)
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			return fmt.Errorf("failed to create out file. error: %s", err)
		}
		pt.Y += c.PointToFixed(12 * 1.5)
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
	fmt.Println("Wrote out.png OK.")

	return nil
}

// NewLocalImage is a constructor of LocalImage
func NewLocalImage(source string) *LocalImage {
	return &LocalImage{source: source}
}

func newImgName(source string) string {
	filePath := filepath.Dir(source)
	fileName := filepath.Base(source)
	fileExt := filepath.Ext(source)
	fileName = strings.Replace(fileName, fileExt, "", 1)
	return filePath + "/" + fileName + "_LGTM_" + fileExt
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
