package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"strings"
)

const StickerSize uint = 512

func StickersHandler(c *gin.Context) {
	runTemplate(c, "stickers", nil)
}

func UploadStickerHandler(c *gin.Context) {
	// Get Input
	inFile, header, err := c.Request.FormFile("upload")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Decode
	img, _, err := image.Decode(inFile)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Scale
	outImg := scaleImage(StickerSize, img)

	// Create Headers
	name := header.Filename[:strings.LastIndex(header.Filename, ".")]
	c.Header("Content-Disposition", `attachment; filename="` + name + `_converted.png"`)
	c.Header("Content-Type", "image/png")

	// Encode
	encoder := png.Encoder{png.BestCompression}
	encoder.Encode(c.Writer, outImg)
}

func scaleImage(imgSideLen uint, img image.Image) image.Image {
	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	newWidth, newHeight := origWidth, origHeight

	if origWidth > origHeight {
		newWidth = imgSideLen
		newHeight = 0
	} else if origWidth < origHeight {
		newWidth = 0
		newHeight = imgSideLen
	} else {
		newWidth = imgSideLen
		newHeight = imgSideLen
	}

	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}