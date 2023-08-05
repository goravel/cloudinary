package cloudinary

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/filesystem"
)

type Cloudinary struct {
	ctx      context.Context
	config   config.Config
	instance *cloudinary.Cloudinary
	disk     string
}

func NewCloudinary(ctx context.Context, config config.Config, disk string) (*Cloudinary, error) {
	cloudName := config.GetString(fmt.Sprintf("filesystems.disks.%s.cloud", disk))
	apiKey := config.GetString(fmt.Sprintf("filesystems.disks.%s.key", disk))
	apiSecret := config.GetString(fmt.Sprintf("filesystems.disks.%s.secret", disk))
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
	}, nil
}

// AllDirectories returns all the directories within a given directory and all its subdirectories.
func (r *Cloudinary) AllDirectories(path string) ([]string, error) {
	var result []string
	folders, err := r.instance.Admin.SubFolders(r.ctx, admin.SubFoldersParams{
		Folder: validPath(path),
	})
	if err != nil {
		return nil, err
	}

	for _, folder := range folders.Folders {
		result = append(result, folder.Path)
		// Recursively call to get directories in the subdirectory
		subdirs, err := r.AllDirectories(folder.Path)
		if err != nil {
			return nil, err
		}
		result = append(result, subdirs...)
	}

	return result, nil
}

// AllFiles returns all the files from the given directory including all its subdirectories.
func (r *Cloudinary) AllFiles(path string) ([]string, error) {
	var result []string
	assetTypes := []api.AssetType{api.Image, api.Video, api.File}
	for _, assetType := range assetTypes {
		nextCursor := ""
		for {
			response, err := r.instance.Admin.Assets(r.ctx, admin.AssetsParams{
				Prefix:       validPath(path),
				DeliveryType: "upload",
				AssetType:    assetType,
				MaxResults:   500,
				NextCursor:   nextCursor,
			})
			if err != nil {
				return nil, err
			}

			for _, folder := range response.Assets {
				result = append(result, folder.PublicID)
			}

			nextCursor = response.NextCursor
			if nextCursor == "" {
				break // Exit the loop when there is no next cursor
			}
		}
	}
	return result, nil
}

// Copy copies a file to a new location.
func (r *Cloudinary) Copy(source, destination string) error {
	result, err := r.instance.Upload.Upload(r.ctx, r.Url(source), uploader.UploadParams{
		PublicID:     r.getPublicId(destination),
		ResourceType: r.getResourceType(destination),
	})
	if err != nil {
		return err
	}
	if result.Error.Message != "" {
		return fmt.Errorf("copy file error: %#v", result.Error)
	}
	return nil
}

// Delete deletes a file.
func (r *Cloudinary) Delete(file ...string) error {
	for _, f := range file {
		result, err := r.instance.Upload.Destroy(r.ctx, uploader.DestroyParams{
			PublicID:     r.getPublicId(f),
			ResourceType: r.getResourceType(f),
			Invalidate:   api.Bool(true),
		})
		if err != nil {
			return err
		}
		if result.Result != "ok" {
			return fmt.Errorf("delete file error: %+v", result.Error)
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

// Directories return all the directories within a given directory.
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
	asset, err := r.getAsset(file)
	if err != nil || asset.Error.Message != "" {
		return false
	}
	return true
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
	rawContent, err := GetRawContent(r.Url(file))
	if err != nil {
		return "", err
	}
	return string(rawContent), nil
}

// LastModified returns the last modified time of a file.
func (r *Cloudinary) LastModified(file string) (time.Time, error) {
	resource, err := r.getAsset(file)
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
		return fmt.Errorf("make directory error: %+v", result.Error)
	}
	return nil
}

// MimeType returns the mime-type of a file.
func (r *Cloudinary) MimeType(file string) (string, error) {
	resource, err := r.getAsset(file)
	if err != nil {
		return "", err
	}
	// Check if the resource format is empty, return only the resource type.
	if resource.Format == "" {
		return resource.ResourceType, nil
	}
	// Replace 'jpg' with 'jpeg' in the format if it is 'jpg'
	format := strings.ReplaceAll(resource.Format, "jpg", "jpeg")

	return resource.ResourceType + "/" + format, nil
}

// Missing checks if a file is missing.
func (r *Cloudinary) Missing(file string) bool {
	return !r.Exists(file)
}

// Move moves a file to a new location.
func (r *Cloudinary) Move(source, destination string) error {
	rename, err := r.instance.Upload.Rename(r.ctx, uploader.RenameParams{
		FromPublicID: r.getPublicId(source),
		ToPublicID:   r.getPublicId(destination),
		ResourceType: r.getResourceType(destination),
	})
	if err != nil {
		return err
	}
	if rename.Error != nil {
		return fmt.Errorf("move file error: %#v", rename.Error)
	}
	return nil
}

// Path returns the full path for a file.
func (r *Cloudinary) Path(file string) string {
	return validPath(file)
}

// Put stores a new file on the disk.
func (r *Cloudinary) Put(file, content string) error {
	tempFile, err := r.tempFile(content)
	defer os.Remove(tempFile.Name())
	if err != nil {
		return err
	}
	_, err = r.instance.Upload.Upload(r.ctx, tempFile.Name(), uploader.UploadParams{
		PublicID:       r.getPublicId(file),
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(false),
		ResourceType:   r.getResourceType(file),
	})
	return err
}

// PutFile stores a new file on the disk.
func (r *Cloudinary) PutFile(path string, source filesystem.File) (string, error) {
	return r.PutFileAs(path, source, source.File())
}

// PutFileAs stores a new file on the disk.
func (r *Cloudinary) PutFileAs(path string, source filesystem.File, name string) (string, error) {
	uploadResult, err := r.instance.Upload.Upload(r.ctx, source.File(), uploader.UploadParams{
		Folder:         validPath(path),
		PublicID:       r.getPublicId(name),
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
	resource, err := r.getAsset(file)
	if err != nil {
		return 0, err
	}
	return int64(resource.Bytes), nil
}

// TemporaryUrl get the temporary url of a file.
func (r *Cloudinary) TemporaryUrl(file string, time time.Time) (string, error) {
	return "", errors.New("cloudinary doesn't support temporary url")
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
	asset, err := r.getAsset(file)
	if err != nil {
		return ""
	}
	return asset.SecureURL
}

func (r *Cloudinary) getAsset(path string) (*uploader.ExplicitResult, error) {
	return r.instance.Upload.Explicit(r.ctx, uploader.ExplicitParams{
		PublicID:     r.getPublicId(path),
		Type:         "upload",
		ResourceType: r.getResourceType(path),
	})
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

func (r *Cloudinary) getResourceType(path string) string {
	extension := strings.TrimPrefix(filepath.Ext(path), ".")
	value := "image"
	resourceTypes := r.config.Get(fmt.Sprintf("filesystems.disks.%s.resource_types", r.disk), defaultResourcesTypes()).(map[string][]string)

	for resourceType, extensions := range resourceTypes {
		for _, ext := range extensions {
			if ext == extension {
				value = resourceType
				break
			}
		}
	}

	return value
}

func (r *Cloudinary) getPublicId(path string) string {
	basename := filepath.Base(path)
	dirname := filepath.Dir(path)
	filename := strings.TrimSuffix(basename, filepath.Ext(basename))
	if r.getResourceType(path) == "raw" {
		filename = basename
	}
	// Prepend the dirname if it exists and is not "."
	if dirname != "." {
		return filepath.Join(dirname, filename)
	}
	return filename
}
