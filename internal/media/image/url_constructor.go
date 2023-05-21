package image

type UrlConstructor interface {
	GetLinkToOriginal(originalPath string) string
	GetLinkToThumbnail(originalPath string, width int, height int, originalWidth int) string
}
