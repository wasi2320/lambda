# Project Name: S3 Image Compression Lambda

## Description

This project implements a GoLang Lambda function that automatically compresses and resizes images uploaded to a specified S3 bucket when they exceed a predefined size threshold. It accomplishes this by listening for SQS messages triggered by S3 events, ensuring efficient resource utilization and cost optimization.

## Key Features

**Automatic compression:** Reduces image file sizes based on format and quality parameters. Intelligent resizing: Maintains aspect ratio while fitting within constraints.
**Configurable:** Customize size thresholds, compression quality, and resizing behavior.
**Modular and reusable:** Structured for easy integration and customization.

### **GoLang 1.17 or later**

AWS account with S3, SQS, and Lambda services configured
imaging library

```go
go get -u github.com/disintegration/imaging
```

### **Clone the repository**

```bash
git clone https://github.com/tabed23/S3-Image-Compression-Lambda.git
```

***Configure constants***
Open main.go and edit the following variables:

1. **bucketName:** Your S3 bucket name.
2. **maxSize:** Maximum allowed image size in MB (default: 5).
3. **quality:** JPEG compression quality (1-100, default: 75).
4. **resizeMode:** Resizing behavior (e.g., fit, fill, lanczos, default: fit).
5. **resizeMaxWidth:** Maximum image width after resizing (optional).
6. **resizeMaxHeight:** Maximum image height after resizing (optional).

#### Build the GoLang binary

```go
go build -o main main.go
```

### Deploy the Lambda function

Use the AWS Management Console, AWS CLI, or Serverless Framework to deploy the main binary to your Lambda function.
Configure an SQS event source for your Lambda function, triggering it on new messages in your S3 event notification queue.

***How it works:***

When an image is uploaded to your S3 bucket, a notification triggers an SQS message. The Lambda function receives the S3 object key from the message.
It retrieves the image from S3 and checks its size against the maxSize threshold. If the size exceeds the threshold  The image is resized proportionally while maintaining aspect ratio, considering resizeMode, resizeMaxWidth, and resizeMaxHeight for flexibility.

The image is compressed based on its format:

1. **JPEG:** Uses the specified quality.
2. **PNG:** Applies interlace optimization.
3. **GIF:** Remains untouched due to inherent lossless compression.

```The compressed image is uploaded back to the same S3 object path, replacing the original.```

***Notes:***

This script currently supports JPEG, PNG, and GIF formats.
Error handling and logging are included for debugging purposes.
Consider adding environment variables for configuration to enhance security and flexibility.
Explore more advanced image processing options using the imaging library for customized functionality.
