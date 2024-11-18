package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
)

type Extension string

const bucketName = "" // add your s3 bucket name

const (
	JPG  Extension = ".jpg"
	JPEG Extension = ".jpeg"
	PNG  Extension = ".png"
	GIF  Extension = ".gif"
)

// Event from S3
type s3Event = struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func main() {
	lambda.Start(HandleRequest)

}
func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) {
	//Get the SQS events Log
	message := sqsEvent.Records[0]
	event := s3Event{}
	err := json.Unmarshal([]byte(message.Body), &event)
	if err != nil {
		fmt.Printf("Error parsing message body: %v\n", err)
		return
	}
	key := event.Records[0].S3.Object.Key
	processReSize(key)
}
func calculateMaxDimensions(bounds image.Rectangle, maxSize int) (int, int) {
	fmt.Printf("Checking dimensions")
	width, height := bounds.Max.X, bounds.Max.Y
	aspectRatio := float64(width) / float64(height)

	if maxWidth := int(math.Sqrt(float64(maxSize) * aspectRatio)); maxWidth < width {
		width = maxWidth
	}

	maxHeight := int(float64(width) / aspectRatio)
	fmt.Printf("New dimensions  Max Height %d  And Max Widthh %d", maxHeight, width)
	return width, maxHeight
}

func getExtension(filePath string) Extension {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return JPG
	case ".png":
		return PNG
	case ".gif":
		return GIF
	default:
		return ""
	}
}

func processReSize(bucketKey string) {
	fmt.Println("Processing Re Size bucket Image...")
	fmt.Println("--------------------------------")
	fmt.Printf("Bucket Key: %s\n", bucketKey)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}
	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
	}
	res, err := svc.GetObject(input)
	if err != nil {
		panic(err)
	}
	maxCompressedSize := 5 * 1024 // Check if the file size is larger then 5MB
	contentReader := res.Body
	defer contentReader.Close()
	content, err := io.ReadAll(contentReader)
	if err != nil {
		panic(err)
	}
	filesize := *res.ContentLength // Get the file size from the S3 Response

	originalImage, err := imaging.Decode(bytes.NewReader(content))
	if err != nil {
		fmt.Printf("Decode image error: %v\n", err)
		return
	}

	if filesize <= int64(maxCompressedSize) {

		fmt.Println("Image size is already within the limit, no resizing needed.")

		return
	} else if filesize > int64(maxCompressedSize) {


		fmt.Printf("Image size is greater than the limit, resizing")

		maxWidth, maxHeight := calculateMaxDimensions(originalImage.Bounds(), maxCompressedSize)

		extension := getExtension(bucketKey)

		resizedImage := imaging.Fit(originalImage, maxWidth, maxHeight, imaging.Lanczos)

		var buf bytes.Buffer
		switch extension {
		case JPG, JPEG:

			err = imaging.Encode(&buf, resizedImage, imaging.JPEG, imaging.JPEGQuality(95))
		case PNG:

			err = imaging.Encode(&buf, resizedImage, imaging.PNG)
		case GIF:

			err = imaging.Encode(&buf, resizedImage, imaging.GIF)
		}

		if err != nil {

			fmt.Printf("Encode image error: %v\n", err)
			return
		}

		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(bucketKey),
			Body:   bytes.NewReader(buf.Bytes()),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Compressed image uploaded back to S3 successfully!")
		return
	} else {
		fmt.Println("Function did not execute successfully")
	}

}
