package imgmeta

import (
	"os"
	"fmt"
	"encoding/json"
	"strings"
)
var cacheFolder = ".metaimage_cache"
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

func saveCache(basePath, filename string, data map[string]string) error {
	cache_folder := fmt.Sprintf("%s/%s", basePath, cacheFolder)
	if _, err := os.Stat(cache_folder); err != nil {
		if err = os.Mkdir(cache_folder, os.ModePerm); err != nil {
			return err
		}
	}
	
	file, err := os.Create(fmt.Sprintf("%s/%s", cache_folder, makeFilename(filename)))
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
