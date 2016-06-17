package imgmeta

import (
	//"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

type ByImageName []*ImgMeta

func (b ByImageName) Len() int {
	return len(b)
}
func (b ByImageName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByImageName) Less(i, j int) bool {
	return b[i].Name < b[j].Name
}

type ImgFolder struct {
	BasePath string
	Name     string
	Path     string
	Folders  []*ImgFolder
	Images   []*ImgMeta
}

func (im *ImgFolder) RelPath() string {
	rel, _ := filepath.Rel(im.BasePath, im.Path)
	return rel
}

func _matchFolder(folders []*ImgFolder, folder string) *ImgFolder {
	for _, f := range folders {
		if f.Name == folder {
			return f
		}
	}
	return nil
}

func (im *ImgFolder) Find(path string) *ImgFolder {
	var parent *ImgFolder
	folders := strings.Split(strings.Trim(path, "/"), "/")
	parent = im
	for _, f := range folders {
		p := _matchFolder(parent.Folders, f)
		if p == nil {
			return nil
		}
		parent = p
	}
	return parent
}

type ByFolderName []*ImgFolder

func (b ByFolderName) Len() int {
	return len(b)
}
func (b ByFolderName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByFolderName) Less(i, j int) bool {
	return b[i].Name < b[j].Name
}

func scanFolder(path string, basePath string, name string) (*ImgFolder, error) {
	re := &ImgFolder{
		BasePath: basePath,
		Path:     path,
		Folders:  make([]*ImgFolder, 0),
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
				re.Folders = append(re.Folders, subfolder)
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
	if len(re.Images) > 0 {
		sort.Sort(ByImageName(re.Images))
	}
	if len(re.Folders) > 0 {
		sort.Sort(sort.Reverse(ByFolderName(re.Folders)))
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
