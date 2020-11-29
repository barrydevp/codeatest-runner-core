package services

import (
	"context"
	"fmt"
	"os"
	// "path"
	// "runtime"

	// "fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	// "time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var GGClient *storage.Client
var DEFAULT_CREDENTIALS string
var DEFAULT_PROJECT_ID string
var DEFAULT_BUCKET_NAME string

func init() {

	credentialsJsonFile := os.Getenv("CREDENTIALS_JSON_FILE")

	if credentialsJsonFile == "" {
		log.Fatal("Missing env: CREDENTIALS_JSON_FILE")
	}

	DEFAULT_CREDENTIALS = credentialsJsonFile

	// DEFAULT_PROJECT_ID = "code-and-t"
	// DEFAULT_BUCKET_NAME = "codeatest"
	DEFAULT_PROJECT_ID = "abiding-weaver-291614c"
	DEFAULT_BUCKET_NAME = "codeatest2"

	GGClient = CreateGGClient(DEFAULT_CREDENTIALS, DEFAULT_PROJECT_ID)

}
func CreateGGClient(jsonPath, projectID string) *storage.Client {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CreateBucket() *storage.BucketHandle {
	return GGClient.Bucket(DEFAULT_BUCKET_NAME)
}

func DownloadFile(file string) (string, error) {
	bucket := CreateBucket()

	rc, err := bucket.Object(file).NewReader(context.TODO())

	if err != nil {
		return "", err
	}

	defer rc.Close()

	slurp, err := ioutil.ReadAll(rc)

	if err != nil {
		return "", err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("./.temp-download/%s", file))

	directory := filepath.Dir(filePath)

	err = os.MkdirAll(directory, 0744)

	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filePath, slurp, 0744)

	return filePath, nil
}
