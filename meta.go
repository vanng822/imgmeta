package imgmeta

import (
	"github.com/vanng822/imgscale/imagick"
)

func getMeta(filename string) map[string]string {
	img := imagick.NewMagickWand() 
	defer img.Destroy()
	img.ReadImage(filename)
	return img.GetImagePropertyValues("*")
}