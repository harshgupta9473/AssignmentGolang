package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)


func ExtractFromField[T any](data []T, fieldname string) string {
	var details []string

	for _, entry := range data {
		val := ""

		if v, ok := any(entry).(map[string]interface{}); ok {
			if fieldvalue, exists := v[fieldname]; exists {
				if strValue, ok := fieldvalue.(string); ok {
					val = strValue
				}
			}
		}
		details = append(details, val)
	}
	return strings.Join(details, ", ")
}

var uploader *s3manager.Uploader

// took help but still doubt

func InitAWS() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	accessKey := os.Getenv("accessKey")
	secretKey := os.Getenv("secretKeyAWS")
	region := os.Getenv("region")
	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretKey,
				"",
			),
		},
	})

	if err != nil {
		panic(err)
	}

	uploader = s3manager.NewUploader(awsSession)
}

func SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	bucketName := os.Getenv("bucketName")
	// Upload the file to S3 using the fileReader
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileHeader.Filename),
		Body:   f,
	})
	if err != nil {
		return "", err
	}

	// Get the URL of the uploaded file
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileHeader.Filename)

	return url, nil
}
