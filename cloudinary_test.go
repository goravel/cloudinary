package cloudinary

import (
	"context"
	"mime"
	"os"
	"testing"
	"time"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"

	configmocks "github.com/goravel/framework/contracts/config/mocks"
	contractsfilesystem "github.com/goravel/framework/contracts/filesystem"
)

func TestStorage(t *testing.T) {
	if os.Getenv("CLOUDINARY_API_KEY") == "" {
		color.Redln("No filesystem tests run, please add Cloudinary configuration: CLOUDINARY_API_KEY= CLOUDINARY_API_SECRET= CLOUDINARY_CLOUD_NAME= CLOUDINARY_URL= go test ./...")
		return
	}

	assert.Nil(t, os.WriteFile("test.txt", []byte("Goravel"), 0644))

	mockConfig := &configmocks.Config{}
	mockConfig.On("GetString", "filesystems.disks.cloudinary.key").Return(os.Getenv("CLOUDINARY_API_KEY"))
	mockConfig.On("GetString", "filesystems.disks.cloudinary.secret").Return(os.Getenv("CLOUDINARY_API_SECRET"))
	mockConfig.On("GetString", "filesystems.disks.cloudinary.cloud").Return(os.Getenv("CLOUDINARY_CLOUD_NAME"))
	mockConfig.On("GetString", "filesystems.disks.cloudinary.url").Return(os.Getenv("CLOUDINARY_URL"))

	var driver contractsfilesystem.Driver
	//url := os.Getenv("CLOUDINARY_URL")

	tests := []struct {
		name  string
		setup func()
	}{
		{
			name: "AllDirectories",
			setup: func() {
				assert.Nil(t, driver.Put("AllDirectories/1.txt", "Goravel"))
				assert.Nil(t, driver.Put("AllDirectories/2.txt", "Goravel"))
				assert.Nil(t, driver.Put("AllDirectories/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put("AllDirectories/3/5/6/6.txt", "Goravel"))
				assert.Nil(t, driver.MakeDirectory("AllDirectories/3/4"))
				assert.True(t, driver.Exists("AllDirectories/1.txt"))
				assert.True(t, driver.Exists("AllDirectories/2.txt"))
				assert.True(t, driver.Exists("AllDirectories/3/3.txt"))
				assert.True(t, driver.Exists("AllDirectories/3/4"))
				assert.True(t, driver.Exists("AllDirectories/3/5/6/6.txt"))
				files, err := driver.AllDirectories("AllDirectories")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllDirectories/3", "AllDirectories/3/4", "AllDirectories/3/5", "AllDirectories/3/5/6"}, files)
				files, err = driver.AllDirectories("./AllDirectories")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllDirectories/3", "AllDirectories/3/4", "AllDirectories/3/5", "AllDirectories/3/5/6"}, files)
				files, err = driver.AllDirectories("/AllDirectories")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllDirectories/3", "AllDirectories/3/4", "AllDirectories/3/5", "AllDirectories/3/5/6"}, files)
				files, err = driver.AllDirectories("./AllDirectories/")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllDirectories/3", "AllDirectories/3/4", "AllDirectories/3/5", "AllDirectories/3/5/6"}, files)
				assert.Nil(t, driver.DeleteDirectory("AllDirectories"))
			},
		},
		{
			name: "AllFiles",
			setup: func() {
				assert.Nil(t, driver.Put("AllFiles/1.txt", "Goravel"))
				assert.Nil(t, driver.Put("AllFiles/2.txt", "Goravel"))
				assert.Nil(t, driver.Put("AllFiles/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put("AllFiles/3/4/4.txt", "Goravel"))
				assert.True(t, driver.Exists("AllFiles/1.txt"))
				assert.True(t, driver.Exists("AllFiles/2.txt"))
				assert.True(t, driver.Exists("AllFiles/3/3.txt"))
				assert.True(t, driver.Exists("AllFiles/3/4/4.txt"))
				files, err := driver.AllFiles("AllFiles")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllFiles/1.txt", "AllFiles/2.txt", "AllFiles/3/3.txt", "AllFiles/3/4/4.txt"}, files)
				files, err = driver.AllFiles("./AllFiles")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllFiles/1.txt", "AllFiles/2.txt", "AllFiles/3/3.txt", "AllFiles/3/4/4.txt"}, files)
				files, err = driver.AllFiles("/AllFiles")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllFiles/1.txt", "AllFiles/2.txt", "AllFiles/3/3.txt", "AllFiles/3/4/4.txt"}, files)
				files, err = driver.AllFiles("./AllFiles/")
				assert.Nil(t, err)
				assert.Equal(t, []string{"AllFiles/1.txt", "AllFiles/2.txt", "AllFiles/3/3.txt", "AllFiles/3/4/4.txt"}, files)
				assert.Nil(t, driver.DeleteDirectory("AllFiles"))
			},
		},
		{
			name: "Copy",
			setup: func() {
				assert.Nil(t, driver.Put("Copy/1.txt", "Goravel"))
				assert.True(t, driver.Exists("Copy/1.txt"))
				assert.Nil(t, driver.Copy("Copy/1.txt", "Copy1/1.txt"))
				assert.True(t, driver.Exists("Copy/1.txt"))
				assert.True(t, driver.Exists("Copy1/1.txt"))
				assert.Nil(t, driver.DeleteDirectory("Copy"))
				assert.Nil(t, driver.DeleteDirectory("Copy1"))
			},
		},
		{
			name: "Delete",
			setup: func() {
				assert.Nil(t, driver.Put("Delete/1.txt", "Goravel"))
				assert.True(t, driver.Exists("Delete/1.txt"))
				assert.Nil(t, driver.Delete("Delete/1.txt"))
				assert.True(t, driver.Missing("Delete/1.txt"))
				assert.Nil(t, driver.DeleteDirectory("Delete"))
			},
		},
		{
			name: "DeleteDirectory",
			setup: func() {
				assert.Nil(t, driver.Put("DeleteDirectory/1.txt", "Goravel"))
				assert.True(t, driver.Exists("DeleteDirectory/1.txt"))
				assert.Nil(t, driver.DeleteDirectory("DeleteDirectory"))
				assert.True(t, driver.Missing("DeleteDirectory/1.txt"))
				assert.Nil(t, driver.DeleteDirectory("DeleteDirectory"))
			},
		},
		{
			name: "Directories",
			setup: func() {
				assert.Nil(t, driver.Put("Directories/1.txt", "Goravel"))
				assert.Nil(t, driver.Put("Directories/2.txt", "Goravel"))
				assert.Nil(t, driver.Put("Directories/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put("Directories/3/5/5.txt", "Goravel"))
				assert.Nil(t, driver.MakeDirectory("Directories/3/4"))
				assert.True(t, driver.Exists("Directories/1.txt"))
				assert.True(t, driver.Exists("Directories/2.txt"))
				assert.True(t, driver.Exists("Directories/3/3.txt"))
				assert.True(t, driver.Exists("Directories/3/4"))
				assert.True(t, driver.Exists("Directories/3/5/5.txt"))
				files, err := driver.Directories("Directories")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Directories/3"}, files)
				files, err = driver.Directories("./Directories")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Directories/3"}, files)
				files, err = driver.Directories("/Directories")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Directories/3"}, files)
				files, err = driver.Directories("./Directories/")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Directories/3"}, files)
				assert.Nil(t, driver.DeleteDirectory("Directories"))
			},
		},
		{
			name: "Files",
			setup: func() {
				assert.Nil(t, driver.Put("Files/1.txt", "Goravel"))
				assert.Nil(t, driver.Put("Files/2.txt", "Goravel"))
				assert.Nil(t, driver.Put("Files/3/3.txt", "Goravel"))
				assert.Nil(t, driver.Put("Files/3/4/4.txt", "Goravel"))
				assert.True(t, driver.Exists("Files/1.txt"))
				assert.True(t, driver.Exists("Files/2.txt"))
				assert.True(t, driver.Exists("Files/3/3.txt"))
				assert.True(t, driver.Exists("Files/3/4/4.txt"))
				files, err := driver.Files("Files")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Files/1.txt", "Files/2.txt"}, files)
				files, err = driver.Files("./Files")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Files/1.txt", "Files/2.txt"}, files)
				files, err = driver.Files("/Files")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Files/1.txt", "Files/2.txt"}, files)
				files, err = driver.Files("./Files/")
				assert.Nil(t, err)
				assert.Equal(t, []string{"Files/1.txt", "Files/2.txt"}, files)
				assert.Nil(t, driver.DeleteDirectory("Files"))
			},
		},
		{
			name: "Get",
			setup: func() {
				assert.Nil(t, driver.Put("Get/1.txt", "Goravel"))
				assert.True(t, driver.Exists("Get/1.txt"))
				data, err := driver.Get("Get/1.txt")
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)
				length, err := driver.Size("Get/1.txt")
				assert.Nil(t, err)
				assert.Equal(t, int64(7), length)
				assert.Nil(t, driver.DeleteDirectory("Get"))
			},
		},
		{
			name: "LastModified",
			setup: func() {
				assert.Nil(t, driver.Put("LastModified/1.txt", "Goravel"))
				assert.True(t, driver.Exists("LastModified/1.txt"))
				assert.Nil(t, driver.Put("LastModified/1.txt", "Goravel-new"))
				date, err := driver.LastModified("LastModified/1.txt")
				assert.Nil(t, err)

				l, err := time.LoadLocation("UTC")
				assert.Nil(t, err)
				assert.Equal(t, time.Now().In(l).Format("2006-01-02 15"), date.Format("2006-01-02 15"))
				assert.Nil(t, driver.DeleteDirectory("LastModified"))
			},
		},
		{
			name: "MakeDirectory",
			setup: func() {
				assert.Nil(t, driver.MakeDirectory("MakeDirectory1/"))
				assert.Nil(t, driver.MakeDirectory("MakeDirectory2"))
				assert.Nil(t, driver.MakeDirectory("MakeDirectory3/MakeDirectory4"))
				assert.Nil(t, driver.DeleteDirectory("MakeDirectory1"))
				assert.Nil(t, driver.DeleteDirectory("MakeDirectory2"))
				assert.Nil(t, driver.DeleteDirectory("MakeDirectory3"))
				assert.Nil(t, driver.DeleteDirectory("MakeDirectory4"))
			},
		},
		{
			name: "MimeType",
			setup: func() {
				assert.Nil(t, driver.Put("MimeType/1.txt", "Goravel"))
				assert.True(t, driver.Exists("MimeType/1.txt"))
				mimeType, err := driver.MimeType("MimeType/1.txt")
				assert.Nil(t, err)
				mediaType, _, err := mime.ParseMediaType(mimeType)
				assert.Nil(t, err)
				assert.Equal(t, "raw", mediaType)

				fileInfo := &File{path: "logo.png"}
				path, err := driver.PutFile("MimeType", fileInfo)
				assert.Nil(t, err)
				assert.True(t, driver.Exists(path))
				mimeType, err = driver.MimeType(path)
				assert.Nil(t, err)
				assert.Equal(t, "image/png", mimeType)
			},
		},
		{
			name: "Move",
			setup: func() {
				assert.Nil(t, driver.Put("Move/1.txt", "Goravel"))
				assert.True(t, driver.Exists("Move/1.txt"))
				assert.Nil(t, driver.Move("Move/1.txt", "Move1/1.txt"))
				assert.True(t, driver.Missing("Move/1.txt"))
				assert.True(t, driver.Exists("Move1/1.txt"))
				assert.Nil(t, driver.DeleteDirectory("Move"))
				assert.Nil(t, driver.DeleteDirectory("Move1"))
			},
		},
		{
			name: "Put",
			setup: func() {
				assert.Nil(t, driver.Put("Put/1.txt", "Goravel"))
				assert.True(t, driver.Exists("Put/1.txt"))
				assert.True(t, driver.Missing("Put/2.txt"))
				assert.Nil(t, driver.DeleteDirectory("Put"))
			},
		},
		{
			name: "PutFile_Image",
			setup: func() {
				fileInfo := &File{path: "logo.png"}
				path, err := driver.PutFile("PutFile1", fileInfo)
				assert.Nil(t, err)
				assert.True(t, driver.Exists(path))
				assert.Nil(t, driver.DeleteDirectory("PutFile1"))
			},
		},
		{
			name: "PutFile_Text",
			setup: func() {
				fileInfo := &File{path: "test.txt"}
				path, err := driver.PutFile("PutFile", fileInfo)
				assert.Nil(t, err)
				assert.True(t, driver.Exists(path))
				data, err := driver.Get(path)
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)
				assert.Nil(t, driver.DeleteDirectory("PutFile"))
			},
		},
		{
			name: "PutFileAs_Text",
			setup: func() {
				fileInfo := &File{path: "test.txt"}
				path, err := driver.PutFileAs("PutFileAs", fileInfo, "text")
				assert.Nil(t, err)
				assert.Equal(t, "PutFileAs/text.txt", path)
				assert.True(t, driver.Exists(path))
				data, err := driver.Get(path)
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)

				path, err = driver.PutFileAs("PutFileAs", fileInfo, "text1.txt")
				assert.Nil(t, err)
				assert.Equal(t, "PutFileAs/text1.txt", path)
				assert.True(t, driver.Exists(path))
				data, err = driver.Get(path)
				assert.Nil(t, err)
				assert.Equal(t, "Goravel", data)

				assert.Nil(t, driver.DeleteDirectory("PutFileAs"))
			},
		},
		{
			name: "PutFileAs_Image",
			setup: func() {
				fileInfo := &File{path: "logo.png"}
				path, err := driver.PutFileAs("PutFileAs1", fileInfo, "image")
				assert.Nil(t, err)
				assert.Equal(t, "PutFileAs1/image", path)
				assert.True(t, driver.Exists(path))

				path, err = driver.PutFileAs("PutFileAs1", fileInfo, "image1.png")
				assert.Nil(t, err)
				assert.Equal(t, "PutFileAs1/image1", path)
				assert.True(t, driver.Exists(path))

				assert.Nil(t, driver.DeleteDirectory("PutFileAs1"))
			},
		},
		{
			name: "Size",
			setup: func() {
				assert.Nil(t, driver.Put("Size/1.txt", "Goravel"))
				assert.True(t, driver.Exists("Size/1.txt"))
				length, err := driver.Size("Size/1.txt")
				assert.Nil(t, err)
				assert.Equal(t, int64(7), length)
				assert.Nil(t, driver.DeleteDirectory("Size"))
			},
		},
		{
			name: "TemporaryUrl",
			setup: func() {
				url, err := driver.TemporaryUrl("TemporaryUrl/1.txt", time.Now().Add(5*time.Second))
				assert.NotNil(t, err)
				assert.Empty(t, url)
			},
		},
		//{
		//	name: "Url",
		//	setup: func() {
		//		assert.Nil(t, driver.Put("Url/1.txt", "Goravel"))
		//		assert.True(t, driver.Exists("Url/1.txt"))
		//		url := url + "/Url/1.txt"
		//		assert.Equal(t, url, driver.Url("Url/1.txt"))
		//		resp, err := http.Get(url)
		//		assert.Nil(t, err)
		//		content, err := io.ReadAll(resp.Body)
		//		assert.Nil(t, resp.Body.Close())
		//		assert.Nil(t, err)
		//		assert.Equal(t, "Goravel", string(content))
		//		assert.Nil(t, driver.DeleteDirectory("Url"))
		//	},
		//},
	}

	var err error
	driver, err = NewCloudinary(context.Background(), mockConfig, "cloudinary")
	assert.NotNil(t, driver)
	assert.Nil(t, err)
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
