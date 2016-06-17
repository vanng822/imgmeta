package imgmeta

import (
	//"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var supportedExts = regexp.MustCompile(".+\\.(?i)(jpg|png)$")

type ImgMeta struct {
	BasePath string
	Name     string
	Path     string
	Meta     map[string]string
}

func (im *ImgMeta) RelPath() string {
	rel, _ := filepath.Rel(im.BasePath, im.FullPath())
	return rel
}

func (im *ImgMeta) FullPath() string {
	return im.Path + "/" + im.Name
}

type ImgFolder struct {
	BasePath string
	Name     string
	Path     string
	Folders  map[string]*ImgFolder
	Images   []*ImgMeta
}

func (im *ImgFolder) RelPath() string {
	rel, _ := filepath.Rel(im.BasePath, im.Path)
	return rel
}

func (im *ImgFolder) Find(path string) *ImgFolder {
	var parent *ImgFolder
	folders := strings.Split(strings.Trim(path, "/"), "/")
	parent = im
	for _, f := range folders {
		ff := parent.Folders
		p, ok := ff[f]
		if !ok {
			return nil
		}
		parent = p
	}
	return parent
}

func scanFolder(path string, basePath string, name string) (*ImgFolder, error) {
	re := &ImgFolder{
		BasePath: basePath,
		Path:     path,
		Folders:  make(map[string]*ImgFolder),
		Name:     name,
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
		// skip all dot files/folders
		if fi.Name()[0] == '.' {
			continue
		}
		if fi.IsDir() {
			subfolder, err := scanFolder(path+"/"+fi.Name(), basePath, fi.Name())
			if err == nil {
				re.Folders[fi.Name()] = subfolder
			}
		}
		if supportedExts.MatchString(fi.Name()) {
			img := &ImgMeta{Name: fi.Name(), Path: path, BasePath: basePath}
			meta := getCache(basePath, img.RelPath())
			if meta == nil {
				meta = getMeta(img.FullPath())
				err = saveCache(basePath, img.RelPath(), meta)
			}
			img.Meta = meta
			re.Images = append(re.Images, img)
		}
	}

	return re, nil
}

func Scan(path string) (*ImgFolder, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return scanFolder(p, p, "")
}
