package cloudinary

import (
	"fmt"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/filesystem"
	"golang.org/x/net/context"
)

type Cloudinary struct {
	ctx      context.Context
	config   config.Config
	instance *cloudinary.Cloudinary
	disk     string
	url      string
}

func NewCloudinary(ctx context.Context, config config.Config, disk string) (*Cloudinary, error) {
	apiSecret := config.GetString(fmt.Sprintf("cloudinary.%s.secret", disk))
	apiKey := config.GetString(fmt.Sprintf("cloudinary.%s.key", disk))
	cloudName := config.GetString(fmt.Sprintf("cloudinary.%s.cloud", disk))
	url := config.GetString(fmt.Sprintf("cloudinary.%s.url", disk))
	if apiSecret == "" || apiKey == "" || cloudName == "" {
		return nil, fmt.Errorf("cloudinary config not found for disk %s", disk)
	}
	client, err := cloudinary.NewFromParams(apiKey, apiSecret, cloudName)
	if err != nil {
		return nil, err
	}
	return &Cloudinary{
		ctx:      ctx,
		config:   config,
		instance: client,
		disk:     disk,
		url:      url,
	}, nil
}

// AllDirectories returns all the directories within a given directory.(it is not typically used for cloud storage systems)
func (r *Cloudinary) AllDirectories(path string) ([]string, error) {
	folders, err := r.instance.Admin.RootFolders(r.ctx, admin.RootFoldersParams{})
	if err != nil {
		return nil, err
	}
	var result []string
	for _, folder := range folders.Folders {
		result = append(result, folder.Path)
	}
	return result, nil
}

// AllFiles returns all the files from the given directory.(it is not typically used for cloud storage systems)
func (r *Cloudinary) AllFiles(path string) ([]string, error) {
	folders, err := r.instance.Admin.Assets(r.ctx, admin.AssetsParams{
		Prefix: path,
	})
	if err != nil {
		return nil, err
	}
	var result []string
	for _, folder := range folders.Assets {
		result = append(result, folder.PublicID)
	}
	return result, nil
}

// Copy copies a file to a new location.
func (r *Cloudinary) Copy(oldFile, newFile string) error {
	_, err := r.instance.Admin.UpdateAsset(r.ctx, admin.UpdateAssetParams{})
	//	Copy(r.ctx, uploader.CopyParams{
	//	FromPublicID: oldFile,
	//	ToPublicID:   newFile,
	//})
	return err
}

// Delete deletes a file.
func (r *Cloudinary) Delete(file ...string) error {
	_, err := r.instance.Admin.DeleteAssets(r.ctx, admin.DeleteAssetsParams{
		PublicIDs: file,
	})
	return err
}

// DeleteDirectory deletes a directory.
func (r *Cloudinary) DeleteDirectory(directory string) error {
	_, err := r.instance.Admin.DeleteFolder(r.ctx, admin.DeleteFolderParams{
		Folder: directory,
	})
	return err
}

// Directories returns all the directories within a given directory.
func (r *Cloudinary) Directories(path string) ([]string, error) {
	folders, err := r.instance.Admin.SubFolders(r.ctx, admin.SubFoldersParams{
		Folder: path,
	})
	if err != nil {
		return nil, err
	}
	var result []string
	for _, folder := range folders.Folders {
		result = append(result, folder.Path)
	}
	return result, nil
}

// Exists checks if a file exists.
func (r *Cloudinary) Exists(file string) bool {
	_, err := r.instance.Upload.Explicit(r.ctx, uploader.ExplicitParams{
		PublicID: file,
	})
	return err == nil
}

// Files returns all the files from the given directory.
func (r *Cloudinary) Files(path string) ([]string, error) {
	folders, err := r.instance.Admin.Assets(r.ctx, admin.AssetsParams{
		Prefix: path,
	})
	if err != nil {
		return nil, err
	}
	var result []string
	for _, folder := range folders.Assets {
		result = append(result, folder.PublicID)
	}
	return result, nil
}

// Get returns the contents of a file.
func (r *Cloudinary) Get(file string) (string, error) {
	return "", nil
}

// LastModified returns the last modified time of a file.
func (r *Cloudinary) LastModified(file string) (time.Time, error) {
	return time.Time{}, nil
}

// MakeDirectory creates a directory.
func (r *Cloudinary) MakeDirectory(directory string) error {
	_, err := r.instance.Admin.CreateFolder(r.ctx, admin.CreateFolderParams{
		Folder: directory,
	})
	return err
}

// MimeType returns the mime-type of a file.
func (r *Cloudinary) MimeType(file string) (string, error) {
	return "", nil
}

// Missing checks if a file is missing.
func (r *Cloudinary) Missing(file string) bool {
	return false
}

// Move moves a file to a new location.
func (r *Cloudinary) Move(oldFile, newFile string) error {
	return nil
}

// Path returns the full path for a file.
func (r *Cloudinary) Path(file string) string {
	return ""
}

// Put stores a new file on the disk.
func (r *Cloudinary) Put(file, content string) error {
	return nil
}

// PutFile stores a new file on the disk.
func (r *Cloudinary) PutFile(path string, source filesystem.File) (string, error) {
	return "", nil
}

// PutFileAs stores a new file on the disk.
func (r *Cloudinary) PutFileAs(path string, source filesystem.File, name string) (string, error) {
	return "", nil
}

// Size returns the file size of a given file.
func (r *Cloudinary) Size(file string) (int64, error) {
	return 0, nil
}

// TemporaryUrl get the temporary url of a file.
func (r *Cloudinary) TemporaryUrl(file string, time time.Time) (string, error) {
	return "", nil
}

// WithContext sets the context for the driver.
func (r *Cloudinary) WithContext(ctx context.Context) filesystem.Driver {
	driver, err := NewCloudinary(ctx, r.config, r.disk)
	if err != nil {
		color.Redf("[Cloudinary] init disk error: %+v\n", err)
		return nil
	}
	return driver
}

// Url returns the url for a file.
func (r *Cloudinary) Url(file string) string {
	return ""
}
