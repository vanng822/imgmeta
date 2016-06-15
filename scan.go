package imgmeta

import (
	//"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var supportedExts = regexp.MustCompile(".+(?i)(jpg|png)$")

type ImgMeta interface {
	Name() string
	Path() string
	FullPath() string
	RelPath() string
	Meta() map[string]string
}

type imgMeta struct {
	basepath string
	name     string
	path     string
	meta     map[string]string
}

func (im *imgMeta) Name() string {
	return im.name
}

func (im *imgMeta) Path() string {
	return im.path
}

func (im *imgMeta) RelPath() string {
	rel, _ := filepath.Rel(im.basepath, im.FullPath())
	return rel
}

func (im *imgMeta) FullPath() string {
	return im.path + "/" + im.Name()
}

func (im *imgMeta) Meta() map[string]string {
	return im.meta
}

type ImgFolder interface {
	Name() string
	Path() string
	Folders() map[string]ImgFolder
	Images() []ImgMeta
	Find(path string) ImgFolder
	RelPath() string
}

type imgFolder struct {
	basepath string
	name     string
	path     string
	folders  map[string]ImgFolder
	images   []ImgMeta
}

func (im *imgFolder) RelPath() string {
	rel, _ := filepath.Rel(im.basepath, im.path)
	return rel
}

func (im *imgFolder) Name() string {
	return im.name
}

func (im *imgFolder) Path() string {
	return im.path
}

func (im *imgFolder) Folders() map[string]ImgFolder {
	return im.folders
}

func (im *imgFolder) Images() []ImgMeta {
	return im.images
}

func (im *imgFolder) Find(path string) ImgFolder {
	var parent ImgFolder
	folders := strings.Split(strings.Trim(path, "/"), "/")
	parent = im
	for _, f := range folders {
		ff := parent.Folders()
		p, ok := ff[f]
		if !ok {
			return nil
		}
		parent = p
	}
	return parent
}

func scanFolder(path string, basepath string, name string) (ImgFolder, error) {
	re := &imgFolder{
		basepath: basepath,
		path:     path,
		folders:  make(map[string]ImgFolder),
		name:     name,
	}
	folder, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := folder.Readdir(-1)

	if err != nil {
		return nil, err
	}
	for _, fi := range files {
		if fi.IsDir() {
			subfolder, err := scanFolder(path+"/"+fi.Name(), basepath, fi.Name())
			if err == nil {
				re.folders[fi.Name()] = subfolder
			}
		}
		if supportedExts.MatchString(fi.Name()) {
			img := &imgMeta{name: fi.Name(), path: path, basepath: basepath}
			img.meta = getMeta(img.FullPath())
			re.images = append(re.images, img)
		}
	}

	return re, nil
}

func Scan(path string) (ImgFolder, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return scanFolder(p, p, "")
}
