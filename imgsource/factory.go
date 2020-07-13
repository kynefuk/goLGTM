package imgsource

import (
	"os"
	"strings"
)

// ImgSrcFactory is a factory of image source struct
type ImgSrcFactory struct {
	source string
}

// GetImgSrc returns image source struct from source
func (factory *ImgSrcFactory) GetImgSrc() Image {
	if strings.HasPrefix(factory.source, "https://") {
		// return online image
	}
	if fileExists(factory.source) {
		return NewLocalImage(factory.source)
	}
	return NewKeywordImage(factory.source)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// NewImgSrcFactory is a constructor of NewImgSrcFactory
func NewImgSrcFactory(source string) *ImgSrcFactory {
	return &ImgSrcFactory{source: source}
}
