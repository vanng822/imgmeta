package imgmeta

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const cacheFolder = ".metaimage_cache"

func makeFilename(filename string) string {
	return strings.Replace(strings.Replace(strings.Trim(filename, "/"), "/", "_", -1), ".", "_", -1)
}

func getCache(basePath, filename string) map[string]string {
	file, err := os.Open(fmt.Sprintf("%s/%s/%s", basePath, cacheFolder, makeFilename(filename)))
	if err != nil {
		return nil
	}
	defer file.Close()
	var res map[string]string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&res)
	if err != nil {
		return nil
	}
	return res
}

func ensureCacheFolder(cFolder string) error {
	if _, err := os.Stat(cFolder); err != nil {
		if err = os.Mkdir(cFolder, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func saveCache(basePath, filename string, data map[string]string) error {
	cFolder := fmt.Sprintf("%s/%s", basePath, cacheFolder)
	if err := ensureCacheFolder(cFolder); err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", cFolder, makeFilename(filename)))
	if err != nil {
		return err
	}
	defer file.Close()
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	_, err = file.Write(d)
	return err
}

func SaveCacheMeta(basePath string, imf *ImgFolder) error {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return err
	}
	cFolder := fmt.Sprintf("%s/%s", p, cacheFolder)
	if err := ensureCacheFolder(cFolder); err != nil {
		return err
	}
	d, err := json.Marshal(imf)
	if err != nil {
		return err
	}
	file, err := os.Create(fmt.Sprintf("%s/.imgmeta", cFolder))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(d)
	return err
}

func LoadCacheMeta(basePath string) (*ImgFolder, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fmt.Sprintf("%s/%s/.imgmeta", p, cacheFolder))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var res *ImgFolder
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
