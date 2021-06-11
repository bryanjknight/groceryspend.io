package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadBase64ImageToS3 uploads a base64 encoded image to s3
func UploadBase64ImageToS3(session *session.Session, bucketName string, key string, base64image string) error {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	if base64image == "" || !strings.Contains(base64image, "base64,") {
		return fmt.Errorf("Failed to get data as base64 data")
	}

	// get the data header info (e.g. data:image/jpeg;base64,)
	// TODO: do something with this information? perhaps not send it at all
	base64key := "base64,"
	base64idx := strings.Index(base64image, base64key)
	base64data := strings.TrimSpace(base64image[base64idx+len(base64key):])

	byteArr, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(byteArr)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	if result == nil {
		return fmt.Errorf("s3 upload result is null")
	}

	return nil
}

// UploadObjectToS3AsJSON uploads the object to S3 as JSON. Requires the object to be marshallable
func UploadObjectToS3AsJSON(session *session.Session, bucketName string, key string, object interface{}) error {

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	resp, err := json.Marshal(object)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(resp))

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	if result == nil {
		return fmt.Errorf("s3 upload result is null")
	}

	return nil
}
