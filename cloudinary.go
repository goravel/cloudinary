package cloudinary

import (
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/goravel/framework/support/str"
	"os"
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
	apiSecret := config.GetString(fmt.Sprintf("filesystems.disks.%s.secret", disk))
	apiKey := config.GetString(fmt.Sprintf("filesystems.disks.%s.key", disk))
	cloudName := config.GetString(fmt.Sprintf("filesystems.disks.%s.cloud", disk))
	url := config.GetString(fmt.Sprintf("filesystems.disks.%s.url", disk))
	if apiSecret == "" || apiKey == "" || cloudName == "" {
		return nil, fmt.Errorf("cloudinary config not found for disk %s", disk)
	}
	client, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		color.Redln("[Cloudinary] init disk error: ", err)
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

// AllDirectories returns all the directories within a given directory and all its subdirectories.
func (r *Cloudinary) AllDirectories(path string) ([]string, error) {
	var result []string
	err := r.getAllDirectoriesRecursively(validPath(path), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getAllDirectoriesRecursively is a helper function to recursively get all directories within a given directory and its subdirectories.
func (r *Cloudinary) getAllDirectoriesRecursively(path string, result *[]string) error {
	folders, err := r.instance.Admin.SubFolders(r.ctx, admin.SubFoldersParams{Folder: path})
	if err != nil {
		return err
	}

	for _, folder := range folders.Folders {
		*result = append(*result, folder.Path)
		// Recursively call to get directories in the subdirectory
		err := r.getAllDirectoriesRecursively(folder.Path, result)
		if err != nil {
			return err
		}
	}

	return nil
}

// AllFiles returns all the files from the given directory including all its subdirectories.
func (r *Cloudinary) AllFiles(path string) ([]string, error) {
	var result []string
	assetTypes := []api.AssetType{api.Image, api.Video, api.File}
	for _, assetType := range assetTypes {
		folders, err := r.instance.Admin.Assets(r.ctx, admin.AssetsParams{
			Prefix:       validPath(path),
			DeliveryType: "upload",
			AssetType:    assetType,
		})
		if err != nil {
			return nil, err
		}
		for _, folder := range folders.Assets {
			result = append(result, folder.PublicID)
		}
	}
	return result, nil
}

// Copy copies a file to a new location.
func (r *Cloudinary) Copy(oldFile, newFile string) error {
	GetAssetType(&oldFile)
	resource, err := r.getResource(oldFile)
	if err != nil {
		return err
	}
	result, err := r.instance.Upload.Upload(r.ctx, resource.SecureURL, uploader.UploadParams{
		PublicID: newFile,
	})
	if err != nil {
		return err
	}
	// Check if the public_id matches the newFile
	if result.PublicID != newFile {
		return fmt.Errorf("copy file error: public_id mismatch, expected %s but got %s", newFile, result.PublicID)
	}
	return nil
}

// Delete deletes a file.
func (r *Cloudinary) Delete(file ...string) error {
	for _, f := range file {
		GetAssetType(&f)
		result, err := r.instance.Upload.Destroy(r.ctx, uploader.DestroyParams{
			PublicID:     f,
			ResourceType: string(GetAssetType(&f)),
		})
		if err != nil {
			return err
		}
		if result.Result != "ok" {
			return fmt.Errorf("delete file error: %+v", result.Error.Message)
		}
	}
	return nil
}

// DeleteDirectory deletes a directory.
func (r *Cloudinary) DeleteDirectory(directory string) error {
	assetTypes := []api.AssetType{api.Image, api.Video, api.File}
	for _, assetType := range assetTypes {
		_, err := r.instance.Admin.DeleteAssetsByPrefix(r.ctx, admin.DeleteAssetsByPrefixParams{
			Prefix:    []string{validPath(directory)},
			AssetType: assetType,
		})
		if err != nil {
			return err
		}
	}
	_, err := r.instance.Admin.DeleteFolder(r.ctx, admin.DeleteFolderParams{
		Folder: directory,
	})
	if err != nil {
		return err
	}
	return nil
}

// Directories returns all the directories within a given directory.
func (r *Cloudinary) Directories(path string) ([]string, error) {
	folders, err := r.instance.Admin.SubFolders(r.ctx, admin.SubFoldersParams{
		Folder: validPath(path),
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

// Exists checks if a file exists in the Cloudinary storage.
func (r *Cloudinary) Exists(file string) bool {
	assetType := GetAssetType(&file)
	asset, err := r.instance.Admin.Asset(r.ctx, admin.AssetParams{
		PublicID:  file,
		AssetType: assetType,
	})
	if asset.Error.Message != "" {
		return false
	}
	return err == nil
}

// Files returns all the files from the given directory.
func (r *Cloudinary) Files(path string) ([]string, error) {
	folders, err := r.instance.Admin.Search(r.ctx, search.Query{
		Expression: fmt.Sprintf("folder:%s", validPath(path)),
		SortBy: []search.SortByField{
			{"public_id": search.Ascending},
		},
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
	assetType := GetAssetType(&file)
	explicitResult, err := r.instance.Upload.Explicit(r.ctx, uploader.ExplicitParams{
		PublicID:     file,
		RawConvert:   file,
		ResourceType: string(assetType),
		Type:         "upload",
	})
	if err != nil {
		return "", err
	}
	if explicitResult.Error.Message != "" {
		return "", fmt.Errorf(explicitResult.Error.Message)
	}
	rawContent, err := GetRawContent(explicitResult.SecureURL)
	if err != nil {
		return "", err
	}
	return string(rawContent), nil
}

// LastModified returns the last modified time of a file.
func (r *Cloudinary) LastModified(file string) (time.Time, error) {
	resource, err := r.getResource(file)
	color.Redln(resource.CreatedAt.Format("2006-01-02 15"))
	if err != nil {
		return time.Time{}, err
	}
	return resource.CreatedAt, nil
}

// MakeDirectory creates a directory.
func (r *Cloudinary) MakeDirectory(directory string) error {
	result, err := r.instance.Admin.CreateFolder(r.ctx, admin.CreateFolderParams{
		Folder: directory,
	})
	if err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("make directory error: %+v", result.Error.Message)
	}
	return nil
}

// MimeType returns the mime-type of a file.
func (r *Cloudinary) MimeType(file string) (string, error) {
	resource, err := r.getResource(file)
	if err != nil {
		return "", err
	}
	// Check if the resource format is empty, return only the resource type.
	if resource.Format == "" {
		return resource.ResourceType, nil
	}

	return resource.ResourceType + "/" + resource.Format, nil
}

// Missing checks if a file is missing.
func (r *Cloudinary) Missing(file string) bool {
	return !r.Exists(file)
}

// Move moves a file to a new location.
func (r *Cloudinary) Move(oldFile, newFile string) error {
	fromType := GetAssetType(&oldFile)
	rename, err := r.instance.Upload.Rename(r.ctx, uploader.RenameParams{
		FromPublicID: validPath(oldFile),
		ToPublicID:   validPath(newFile),
		ResourceType: string(fromType),
	})
	if err != nil {
		return err
	}
	if rename.PublicID != newFile {
		return fmt.Errorf("move file error: public_id mismatch, expected %s but got %s", newFile, rename.PublicID)
	}
	return nil
}

// Path returns the full path for a file.
func (r *Cloudinary) Path(file string) string {
	//resource, err := r.getResource(file)
	//if err != nil {
	//	color.Redln("[Cloudinary] get resource error: ", err)
	//	return ""
	//}
	file = validPath(file)
	RemoveFileExtension(&file)
	return file
}

// Put stores a new file on the disk.
func (r *Cloudinary) Put(file, content string) error {
	tempFile, err := r.tempFile(content)
	defer os.Remove(tempFile.Name())
	if err != nil {
		return err
	}
	_, err = r.instance.Upload.Upload(r.ctx, tempFile.Name(), uploader.UploadParams{
		PublicID:    file,
		UseFilename: api.Bool(true),
	})
	if err != nil {
		return err
	}
	return nil
}

// PutFile stores a new file on the disk.
func (r *Cloudinary) PutFile(path string, source filesystem.File) (string, error) {
	return r.PutFileAs(path, source, str.Random(40))
}

// PutFileAs stores a new file on the disk.
func (r *Cloudinary) PutFileAs(path string, source filesystem.File, name string) (string, error) {
	RemoveFileExtension(&name)
	uploadResult, err := r.instance.Upload.Upload(r.ctx, source.File(), uploader.UploadParams{
		Folder:         validPath(path),
		PublicID:       name,
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(false),
	})
	if err != nil {
		return "", err
	}
	return uploadResult.PublicID, nil
}

// Size returns the file size of a given file.
func (r *Cloudinary) Size(file string) (int64, error) {
	resource, err := r.getResource(file)
	if err != nil {
		return 0, err
	}
	return int64(resource.Bytes), nil
}

// TemporaryUrl get the temporary url of a file.
func (r *Cloudinary) TemporaryUrl(file string, time time.Time) (string, error) {
	return "", errors.New("we don't support temporary url for cloudinary")
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
	asset, err := r.getResource(file)
	if err != nil {
		return ""
	}
	return asset.SecureURL
}

func (r *Cloudinary) getResource(path string) (*admin.AssetResult, error) {
	assetType := GetAssetType(&path)
	result, err := r.instance.Admin.Asset(r.ctx, admin.AssetParams{
		PublicID:  path,
		AssetType: assetType, // Set the asset type dynamically based on the file extension.
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Cloudinary) tempFile(content string) (*os.File, error) {
	tempFile, err := os.CreateTemp(os.TempDir(), "goravel-")
	if err != nil {
		return nil, err
	}

	if _, err := tempFile.WriteString(content); err != nil {
		return nil, err
	}

	return tempFile, nil
}
