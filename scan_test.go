package imgmeta

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)

func cleanCache(basePath string) {
	os.RemoveAll(fmt.Sprintf("%s/%s", basePath, cacheFolder))
}

func TestMain(m *testing.M) {
	fmt.Println("Test starting")
	retCode := m.Run()
	cleanCache("./data")
	fmt.Println("Test ending")
	os.Exit(retCode)
}


func TestFetchOK(t *testing.T) {
	f, _ := Scan("./data")
	assert.NotNil(t, f)
	assert.Equal(t, 1, len(f.Folders))
	ff := f.Find("/photos/")
	assert.NotNil(t, ff)
	assert.Equal(t, "photos", ff.Name)
	assert.Equal(t, "photos", ff.RelPath())
	assert.Equal(t, 0, len(ff.Folders))
	assert.Equal(t, 2, len(ff.Images))
}
