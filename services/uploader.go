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
var CREDENTIALS string
var PROJECT_ID string
var BUCKET_NAME string

func init() {

	CREDENTIALS = os.Getenv("CREDENTIALS_JSON_FILE")
	PROJECT_ID = os.Getenv("PROJECT_ID")
	BUCKET_NAME = os.Getenv("BUCKET_NAME")

	if CREDENTIALS == "" {
		log.Fatal("Missing env: CREDENTIALS_JSON_FILE")
	}

	if PROJECT_ID == "" {
		PROJECT_ID = "blissful-star-290200"
	}

	if BUCKET_NAME == "" {
		BUCKET_NAME = "codeatest-2"
	}

	GGClient = CreateGGClient(CREDENTIALS, PROJECT_ID)

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
	return GGClient.Bucket(BUCKET_NAME)
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
