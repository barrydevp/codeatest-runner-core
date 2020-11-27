package services

import (
	"context"
	"log"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var GGClient *storage.Client
var DEFAULT_CREDENTIALS string
var DEFAULT_PROJECT_ID string
var DEFAULT_BUCKET_NAME string

func init() {
	DEFAULT_CREDENTIALS, _ = filepath.Abs("./services/credentials.json")
	DEFAULT_PROJECT_ID = "code-and-t"
	DEFAULT_BUCKET_NAME = "codeatest"

	GGClient = CreateGGClient(DEFAULT_CREDENTIALS, DEFAULT_PROJECT_ID)
}

func CreateGGClient(jsonPath, projectID string) *storage.Client {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(jsonPath)))

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CreateBucket() *storage.BucketHandle {
	return GGClient.Bucket(DEFAULT_BUCKET_NAME)
}

func GenerateSignedUrl(file) string {
	bucket := CreateBucket()

}

func generateV4GetObjectSignedURL(w io.Writer, bucket, object, serviceAccount string) (string, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	// serviceAccount := "service_account.json"
	jsonKey, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}
	u, err := storage.SignedURL(bucket, object, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}

	fmt.Fprintln(w, "Generated GET signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl %q\n", u)
	return u, nil
}
