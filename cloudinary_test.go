package cloudinary

import (
	"context"
	"crypto/rand"
	"io"
	"math/big"
	"mime"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/gookit/color"
	configmocks "github.com/goravel/framework/contracts/config/mocks"
	contractsfilesystem "github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	if os.Getenv("CLOUDINARY_ACCESS_KEY_ID") == "" {
		color.Redln("No filesystem tests run, please add Cloudinary configuration: CLOUDINARY_ACCESS_KEY_ID= CLOUDINARY_ACCESS_KEY_SECRET= CLOUDINARY_CLOUD= go test ./...")
		return
	}

	assert.Nil(t, os.WriteFile("test.txt", []byte("Goravel"), 0644))

	mockConfig := &configmocks.Config{}
	mockConfig.On("GetString", "filesystems.disks.cloudinary.key").Return(os.Getenv("CLOUDINARY_ACCESS_KEY_ID"))
	mockConfig.On("GetString", "filesystems.disks.cloudinary.secret").Return(os.Getenv("CLOUDINARY_ACCESS_KEY_SECRET"))
	mockConfig.On("GetString", "filesystems.disks.cloudinary.cloud").Return(os.Getenv("CLOUDINARY_CLOUD"))
	mockConfig.On("Get", "filesystems.disks.cloudinary.resource_types", map[string][]string{
		"image": {"png"},
		"video": {},
		"raw":   {"txt", "pdf"},
	}).Return(map[string][]string{
		"image": {"png"},
		"video": {},
		"raw":   {"txt", "pdf"},
	})

	var driver contractsfilesystem.Driver
	randNum, err := rand.Int(rand.Reader, big.NewInt(1000))
	rootFolder := randNum.String() + "/"
	assert.Nil(t, err)
	driver, err = NewCloudinary(context.Background(), mockConfig, "cloudinary")
	assert.NotNil(t, driver)
	assert.Nil(t, err)

	tests := []struct {
		name  string
		setup func()
	}{
		{
			name: "AllDirectories",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"AllDirectories/1.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"AllDirectories/2.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"AllDirectories/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"AllDirectories/3/5/6/6.txt", "Goravel"))
				assert.Nil(t, driver.MakeDirectory(rootFolder+"AllDirectories/3/4"))
				assert.True(t, driver.Exists(rootFolder+"AllDirectories/1.txt"))
				assert.True(t, driver.Exists(rootFolder+"AllDirectories/2.txt"))
				assert.True(t, driver.Exists(rootFolder+"AllDirectories/3/3.txt"))
				assert.True(t, driver.Exists(rootFolder+"AllDirectories/3/5/6/6.txt"))
				files, err := driver.AllDirectories(rootFolder + "AllDirectories")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllDirectories/3", rootFolder + "AllDirectories/3/4", rootFolder + "AllDirectories/3/5", rootFolder + "AllDirectories/3/5/6"}, files)
				files, err = driver.AllDirectories("./" + rootFolder + "AllDirectories")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllDirectories/3", rootFolder + "AllDirectories/3/4", rootFolder + "AllDirectories/3/5", rootFolder + "AllDirectories/3/5/6"}, files)
				files, err = driver.AllDirectories("/" + rootFolder + "AllDirectories")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllDirectories/3", rootFolder + "AllDirectories/3/4", rootFolder + "AllDirectories/3/5", rootFolder + "AllDirectories/3/5/6"}, files)
				files, err = driver.AllDirectories("./" + rootFolder + "AllDirectories/")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllDirectories/3", rootFolder + "AllDirectories/3/4", rootFolder + "AllDirectories/3/5", rootFolder + "AllDirectories/3/5/6"}, files)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"AllDirectories"))
			},
		},
		{
			name: "AllFiles",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"AllFiles/1.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"AllFiles/2.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"AllFiles/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"AllFiles/3/4/4.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"AllFiles/1.txt"))
				assert.True(t, driver.Exists(rootFolder+"AllFiles/2.txt"))
				assert.True(t, driver.Exists(rootFolder+"AllFiles/3/3.txt"))
				assert.True(t, driver.Exists(rootFolder+"AllFiles/3/4/4.txt"))
				files, err := driver.AllFiles(rootFolder + "AllFiles")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllFiles/1.txt", rootFolder + "AllFiles/2.txt", rootFolder + "AllFiles/3/3.txt", rootFolder + "AllFiles/3/4/4.txt"}, files)
				files, err = driver.AllFiles("./" + rootFolder + "AllFiles")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllFiles/1.txt", rootFolder + "AllFiles/2.txt", rootFolder + "AllFiles/3/3.txt", rootFolder + "AllFiles/3/4/4.txt"}, files)
				files, err = driver.AllFiles("/" + rootFolder + "AllFiles")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllFiles/1.txt", rootFolder + "AllFiles/2.txt", rootFolder + "AllFiles/3/3.txt", rootFolder + "AllFiles/3/4/4.txt"}, files)
				files, err = driver.AllFiles("./" + rootFolder + "AllFiles/")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "AllFiles/1.txt", rootFolder + "AllFiles/2.txt", rootFolder + "AllFiles/3/3.txt", rootFolder + "AllFiles/3/4/4.txt"}, files)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"AllFiles"))
			},
		},
		{
			name: "Copy",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Copy/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Copy/1.txt"))
				assert.Nil(t, driver.Copy(rootFolder+"Copy/1.txt", rootFolder+"Copy1/1.txt"))
				assert.True(t, driver.Exists(rootFolder+"Copy/1.txt"))
				assert.True(t, driver.Exists(rootFolder+"Copy1/1.txt"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Copy"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Copy1"))
			},
		},
		{
			name: "Delete",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Delete/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Delete/1.txt"))
				assert.Nil(t, driver.Delete(rootFolder+"Delete/1.txt"))
				assert.True(t, driver.Missing(rootFolder+"Delete/1.txt"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Delete"))
			},
		},
		{
			name: "DeleteDirectory",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"DeleteDirectory/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"DeleteDirectory/1.txt"))
				fileInfo := &File{path: "logo.png"}
				path, err := driver.PutFile(rootFolder+"DeleteDirectory", fileInfo)
				assert.Nil(t, err)
				assert.True(t, driver.Exists(path))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"DeleteDirectory"))
				assert.True(t, driver.Missing(rootFolder+"DeleteDirectory/1.txt"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"DeleteDirectory"))
			},
		},
		{
			name: "Directories",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Directories/1.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"Directories/2.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"Directories/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"Directories/3/5/5.txt", "Goravel"))
				assert.Nil(t, driver.MakeDirectory(rootFolder+"Directories/3/4"))
				assert.True(t, driver.Exists(rootFolder+"Directories/1.txt"))
				assert.True(t, driver.Exists(rootFolder+"Directories/2.txt"))
				assert.True(t, driver.Exists(rootFolder+"Directories/3/3.txt"))
				assert.True(t, driver.Exists(rootFolder+"Directories/3/5/5.txt"))
				files, err := driver.Directories(rootFolder + "Directories")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Directories/3"}, files)
				files, err = driver.Directories("./" + rootFolder + "Directories")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Directories/3"}, files)
				files, err = driver.Directories("/" + rootFolder + "Directories")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Directories/3"}, files)
				files, err = driver.Directories("./" + rootFolder + "Directories/")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Directories/3"}, files)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Directories"))
			},
		},
		{
			name: "Files",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Files/1.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"Files/2.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"Files/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put(rootFolder+"Files/3/4/4.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Files/1.txt"))
				assert.True(t, driver.Exists(rootFolder+"Files/2.txt"))
				assert.True(t, driver.Exists(rootFolder+"Files/3/3.txt"))
				assert.True(t, driver.Exists(rootFolder+"Files/3/4/4.txt"))
				files, err := driver.Files(rootFolder + "Files")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Files/1.txt", rootFolder + "Files/2.txt"}, files)
				files, err = driver.Files("./" + rootFolder + "Files")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Files/1.txt", rootFolder + "Files/2.txt"}, files)
				files, err = driver.Files("/" + rootFolder + "Files")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Files/1.txt", rootFolder + "Files/2.txt"}, files)
				files, err = driver.Files("./" + rootFolder + "Files/")
				assert.Nil(t, err)
				assert.Equal(t, []string{rootFolder + "Files/1.txt", rootFolder + "Files/2.txt"}, files)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Files"))
			},
		},
		{
			name: "Get",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Get/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Get/1.txt"))
				data, err := driver.Get(rootFolder + "Get/1.txt")
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)
				length, err := driver.Size(rootFolder + "Get/1.txt")
				assert.Nil(t, err)
				assert.Equal(t, int64(7), length)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Get"))
			},
		},
		{
			name: "LastModified",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"LastModified/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"LastModified/1.txt"))
				assert.Nil(t, driver.Put(rootFolder+"LastModified/1.txt", "Goravel-new"))
				date, err := driver.LastModified(rootFolder + "LastModified/1.txt")
				assert.Nil(t, err)

				l, err := time.LoadLocation("UTC")
				assert.Nil(t, err)
				assert.Equal(t, carbon.Now().ToStdTime().In(l).Format("2006-01-02 15"), date.Format("2006-01-02 15"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"LastModified"))
			},
		},
		{
			name: "MakeDirectory",
			setup: func() {
				assert.Nil(t, driver.MakeDirectory(rootFolder+"MakeDirectory1/"))
				assert.Nil(t, driver.MakeDirectory(rootFolder+"MakeDirectory2"))
				assert.Nil(t, driver.MakeDirectory(rootFolder+"MakeDirectory3/MakeDirectory4"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"MakeDirectory1"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"MakeDirectory2"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"MakeDirectory3"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"MakeDirectory4"))
			},
		},
		{
			name: "MimeType",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"MimeType/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"MimeType/1.txt"))
				mimeType, err := driver.MimeType(rootFolder + "MimeType/1.txt")
				assert.Nil(t, err)
				mediaType, _, err := mime.ParseMediaType(mimeType)
				assert.Nil(t, err)
				assert.Equal(t, "raw", mediaType)

				fileInfo := &File{path: "logo.png"}
				path, err := driver.PutFile(rootFolder+"MimeType", fileInfo)
				assert.Nil(t, err)
				assert.True(t, driver.Exists(path))
				mimeType, err = driver.MimeType(path)
				assert.Nil(t, err)
				assert.Equal(t, "image/png", mimeType)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"MimeType"))
			},
		},
		{
			name: "Move",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Move/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Move/1.txt"))
				assert.Nil(t, driver.Move(rootFolder+"Move/1.txt", rootFolder+"Move1/1.txt"))
				assert.True(t, driver.Missing(rootFolder+"Move/1.txt"))
				assert.True(t, driver.Exists(rootFolder+"Move1/1.txt"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Move"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Move1"))
			},
		},
		{
			name: "Put",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Put/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Put/1.txt"))
				assert.True(t, driver.Missing(rootFolder+"Put/2.txt"))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Put"))
			},
		},
		{
			name: "PutFile_Image",
			setup: func() {
				fileInfo := &File{path: "logo.png"}
				path, err := driver.PutFile(rootFolder+"PutFile1", fileInfo)
				assert.Nil(t, err)
				assert.True(t, driver.Exists(path))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"PutFile1"))
			},
		},
		{
			name: "PutFile_Text",
			setup: func() {
				fileInfo := &File{path: "test.txt"}
				path, err := driver.PutFile(rootFolder+"PutFile", fileInfo)
				assert.Nil(t, err)
				assert.True(t, driver.Exists(path))
				data, err := driver.Get(path)
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"PutFile"))
			},
		},
		{
			name: "PutFileAs_Text",
			setup: func() {
				fileInfo := &File{path: "test.txt"}
				path, err := driver.PutFileAs(rootFolder+"PutFileAs", fileInfo, "text")
				assert.Nil(t, err)
				assert.Equal(t, rootFolder+"PutFileAs/text.txt", path)
				assert.True(t, driver.Exists(path))
				data, err := driver.Get(path)
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)

				path, err = driver.PutFileAs(rootFolder+"PutFileAs", fileInfo, "text1.txt")
				assert.Nil(t, err)
				assert.Equal(t, rootFolder+"PutFileAs/text1.txt", path)
				assert.True(t, driver.Exists(path))
				data, err = driver.Get(path)
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)

				assert.Nil(t, driver.DeleteDirectory(rootFolder+"PutFileAs"))
			},
		},
		{
			name: "PutFileAs_Image",
			setup: func() {
				fileInfo := &File{path: "logo.png"}
				path, err := driver.PutFileAs(rootFolder+"PutFileAs1", fileInfo, "image")
				assert.Nil(t, err)
				assert.Equal(t, rootFolder+"PutFileAs1/image", path)
				assert.True(t, driver.Exists(path))

				path, err = driver.PutFileAs(rootFolder+"PutFileAs1", fileInfo, "image1.png")
				assert.Nil(t, err)
				assert.Equal(t, rootFolder+"PutFileAs1/image1", path)
				assert.True(t, driver.Exists(path))

				assert.Nil(t, driver.DeleteDirectory(rootFolder+"PutFileAs1"))
			},
		},
		{
			name: "Size",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Size/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Size/1.txt"))
				length, err := driver.Size(rootFolder + "Size/1.txt")
				assert.Nil(t, err)
				assert.Equal(t, int64(7), length)
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Size"))
			},
		},
		{
			name: "TemporaryUrl",
			setup: func() {
				url, err := driver.TemporaryUrl(rootFolder+"TemporaryUrl/1.txt", time.Now().Add(5*time.Second))
				assert.NotNil(t, err)
				assert.Empty(t, url)
			},
		},
		{
			name: "Url",
			setup: func() {
				assert.Nil(t, driver.Put(rootFolder+"Url/1.txt", "Goravel"))
				assert.True(t, driver.Exists(rootFolder+"Url/1.txt"))
				url := driver.Url(rootFolder + "Url/1.txt")
				matches := regexp.MustCompile(`/v\d+/(.+)`).FindStringSubmatch(url)
				assert.Equal(t, rootFolder+"Url/1.txt", matches[1])
				resp, err := http.Get(url)
				assert.Nil(t, err)
				content, err := io.ReadAll(resp.Body)
				assert.Nil(t, resp.Body.Close())
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", string(content))
				assert.Nil(t, driver.DeleteDirectory(rootFolder+"Url"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
		})
	}

	assert.Nil(t, os.Remove("test.txt"))
}

type File struct {
	path string
}

func (f *File) Disk(disk string) contractsfilesystem.File {
	return &File{}
}

func (f *File) Extension() (string, error) {
	return "", nil
}

func (f *File) File() string {
	return f.path
}

func (f *File) GetClientOriginalName() string {
	return ""
}

func (f *File) GetClientOriginalExtension() string {
	return ""
}

func (f *File) HashName(path ...string) string {
	return ""
}

func (f *File) LastModified() (time.Time, error) {
	return time.Now(), nil
}

func (f *File) MimeType() (string, error) {
	return "", nil
}

func (f *File) Size() (int64, error) {
	return 0, nil
}

func (f *File) Store(path string) (string, error) {
	return "", nil
}

func (f *File) StoreAs(path string, name string) (string, error) {
	return "", nil
}
