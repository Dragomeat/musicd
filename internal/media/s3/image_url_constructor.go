package s3

import (
	"fmt"
	"hash/crc32"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultWidth   = 0
	defaultHeight  = 0
	defaultFormat  = "jpg"
	defaultQuality = 85
)

type ImageUrlConstructor struct {
	cdnHost          string
	protectionSecret string
	regHttp          *regexp.Regexp
}

func NewImageUrlConstructor(cdnHost string, protectionSecret string) *ImageUrlConstructor {
	return &ImageUrlConstructor{
		cdnHost:          cdnHost,
		protectionSecret: protectionSecret,
		regHttp:          regexp.MustCompile(`^(http://|https://)`),
	}
}

func (c *ImageUrlConstructor) GetLinkToOriginal(originalPath string) string {
	return c.formUrl(originalPath, defaultWidth, defaultHeight, formatFromPath(originalPath), 100)
}

func (c *ImageUrlConstructor) GetLinkToThumbnail(originalPath string, width int, height int, originalWidth int) string {
	if originalWidth < width {
		return c.GetLinkToOriginal(originalPath)
	}
	return c.formUrl(originalPath, width, height, defaultFormat, defaultQuality)
}

func (c *ImageUrlConstructor) formUrl(
	originalPath string,
	width int,
	height int,
	format string,
	quality int,
) string {
	if c.regHttp.Match([]byte(originalPath)) {
		return originalPath
	}

	params := url.Values{}

	w := strconv.Itoa(width)
	h := strconv.Itoa(height)
	q := strconv.Itoa(quality)
	format = correctFormat(format)

	if w == "0" {
		w = "AUTO"
	}

	if h == "0" {
		h = "AUTO"
	}

	if q == "0" {
		q = "100"
	}

	originalFormat := formatFromPath(originalPath)

	if w == "AUTO" && h == "AUTO" && q == "100" && format == originalFormat {
		return fmt.Sprintf("%s/%s", c.cdnHost, originalPath)
	}

	params.Add("s", c.signUrlV2(originalPath, w, h, q))

	if w != "AUTO" {
		params.Add("w", w)
	}
	if h != "AUTO" {
		params.Add("h", h)
	}
	if format != originalFormat {
		params.Add("f", format)
	}
	if q != "100" {
		params.Add("q", q)
	}

	return fmt.Sprintf(
		"%s/%s?%s",
		c.cdnHost,
		originalPath,
		params.Encode(),
	)
}

func (c *ImageUrlConstructor) signUrlV2(
	originalPath string,
	width string,
	height string,
	quality string,
) string {
	data := fmt.Sprintf(
		"%s%s%s%s%s",
		c.protectionSecret,
		originalPath,
		width,
		height,
		quality,
	)

	checksum := crc32.ChecksumIEEE([]byte(data))

	return "v2" + strconv.FormatInt(int64(checksum), 16)
}

func formatFromPath(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		return correctFormat(filepath.Ext(path)[1:])
	} else {
		return ext
	}
}

func correctFormat(format string) string {
	format = strings.ToLower(format)
	if format == "jpeg" {
		return "jpg"
	}
	return format
}
