package api

import (
	"musicd/graph/model"
	mediaImage "musicd/internal/media/image"
)

type Transformer struct {
	urlConstructor mediaImage.UrlConstructor
}

func NewTransformer(urlConstructor mediaImage.UrlConstructor) *Transformer {
	return &Transformer{urlConstructor: urlConstructor}
}

func (t *Transformer) TransformImagesToGraphql(images mediaImage.Images) []*model.Image {
	graphqlImages := make([]*model.Image, len(images))
	for i, image := range images {
		graphqlImages[i] = t.TransformImageToGraphql(image)
	}
	return graphqlImages
}

func (t *Transformer) TransformImageToGraphql(image mediaImage.Image) *model.Image {
	return &model.Image{
		ID:  image.Id,
		URL: t.urlConstructor.GetLinkToOriginal(image.Id),
		Thumbnails: &model.ImageThumbnails{
			S295:  t.getThumbnail(image, 295, 0),
			M600:  t.getThumbnail(image, 600, 0),
			B960:  t.getThumbnail(image, 960, 0),
			W1200: t.getThumbnail(image, 1200, 0),
			F1920: t.getThumbnail(image, 1920, 0),
		},
		Sizes: &model.ImageSizes{
			Width:  image.Sizes.Width,
			Height: image.Sizes.Height,
		},
	}
}

func (t *Transformer) getThumbnail(
	image mediaImage.Image,
	width int,
	height int,
) *model.ImageThumbnail {
	sizes := mediaImage.Sizes{Width: width, Height: height}

	if image.Sizes.Height != 0 && sizes.Width == 0 && sizes.Height != 0 {
		sizes.Width = (sizes.Height / image.Sizes.Height) * image.Sizes.Width
	}

	if image.Sizes.Width != 0 && sizes.Height == 0 && sizes.Width != 0 {
		sizes.Height = (sizes.Width / image.Sizes.Width) * image.Sizes.Height
	}

	return &model.ImageThumbnail{
		URL: t.urlConstructor.GetLinkToThumbnail(image.Id, sizes.Width, sizes.Height, image.Sizes.Width),
		Sizes: &model.ImageSizes{
			Width:  sizes.Width,
			Height: sizes.Height,
		},
	}
}
