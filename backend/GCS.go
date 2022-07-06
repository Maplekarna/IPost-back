package backend

import (
	"around/constants"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

var (
	GCSBackend *GoogleCloudStorageBackend
)

func InitGCSBackend() {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}
	GCSBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: constants.GCS_BUCKET,
	}
}

// SaveToGCS Operate to a pointer backend which type is GoogleCloudStorageBackend.
func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)
	wc := object.NewWriter(ctx)

	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}
