package processor

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
)

// MaxRatio is a ratio of messaging aria
const MaxRatio float64 = 0.8

// MaxFontSize is max size of font size
const MaxFontSize float64 = 100

// MinFontSize is minimum size of font size
const MinFontSize float64 = 1

// Processor put message on image
type Processor struct {
	img     image.Image
	message string
}

// AddMessage put message on image
func (p *Processor) AddMessage() (*image.RGBA, error) {
	var text = []string{p.message}
	// message aria
	msgAriaWidth := float64(p.img.Bounds().Dy()) * MaxRatio

	// parse font
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font file. error: %s", err)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, p.img.Bounds().Dx(), p.img.Bounds().Dy()))

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
		dr := &font.Drawer{Dst: rgba, Src: p.img, Face: face, Dot: fixed.Point26_6{}}
		messageXPoint = ((fixed.I(p.img.Bounds().Dx()) - dr.MeasureString(p.message)) / 2).Floor()

		messageAria := int(msgAriaWidth) - dr.MeasureString(p.message).Floor()
		if messageAria > 0 {
			fontSize = i
			break
		}
	}

	fontColor := image.Black
	draw.Draw(rgba, rgba.Bounds(), p.img, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ft)
	c.SetFontSize(float64(fontSize))
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fontColor)
	c.SetHinting(font.HintingVertical)

	messageYPoint := p.img.Bounds().Dy() / 2
	// Draw the text.
	pt := freetype.Pt(messageXPoint, messageYPoint)
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			return nil, fmt.Errorf("failed to create out file. error: %s", err)
		}
		pt.Y += c.PointToFixed(12 * 1.5)
	}

	return rgba, nil
}

// NewProcessor is a constructor of Processor
func NewProcessor(img image.Image, message string) *Processor {
	return &Processor{img: img, message: message}
}
