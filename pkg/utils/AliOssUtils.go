package utils

import (
	"fmt"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// ExtractAnySuffix 提取字符串的任意后缀
func extractAnySuffix(s string) string {
	lastIndex := strings.LastIndex(s, ".")
	if lastIndex == -1 {
		return "" // 如果没有点，则返回空字符串
	}
	return s[lastIndex:]
}

func accessUrl(endpoint string, bucketName string, objectKey string) string {
	if strings.HasPrefix(endpoint, "https://") {
		return strBuilder(endpoint[:8], bucketName, ".", endpoint[8:], "/", objectKey)
	} else if strings.HasPrefix(endpoint, "http://") {
		return strBuilder(endpoint[:7], bucketName, ".", endpoint[7:], "/", objectKey)
	}
	return strBuilder("https://", bucketName, ".", endpoint, "/", objectKey)
}

func strBuilder(args ...string) string {
	builder := strings.Builder{}
	for _, k := range args {
		builder.WriteString(k)
	}
	return builder.String()
}

func CheckFileSize(filePath string) error {
	const maxFileSize = 5 * 1024 * 1024 // 5MB

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	if fileInfo.Size() > maxFileSize {
		return fmt.Errorf("file size exceeds the maximum limit")
	}

	// If the file size is within the limit, you can proceed with further processing.
	// For example, reading the file content:
	buffer := make([]byte, maxFileSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Use the read file content
	fmt.Println("Read", n, "bytes from the file")
	// Add your logic here to process the data in buffer
	return nil
}

// UploadImages 阿里OSS对象存储上传图片
// 最大上传限制5MB
func UploadImages(localFilePath string, serviceName string, userId uint32) (string, error) {

	if localFilePath == "" {
		return "", nil
	}

	if err := CheckFileSize(localFilePath); err != nil {
		return "", err
	}

	endpoint := config.Conf.AliyunConfig.Oss.Endpoint
	accessKeyID := config.Conf.AliyunConfig.Oss.AccessKeyId
	accessKeySecret := config.Conf.AliyunConfig.Oss.AccessKeySecret

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Fatalf("Failed to create OSS client: %v", err)
		return "", err
	}

	// 填写存储空间名称，例如examplebucket。
	bucketName := config.Conf.AliyunConfig.Oss.BucketName // 请替换为实际的Bucket名称
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
		return "", err
	}

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	objectKey := serviceName + "_" + strconv.Itoa(int(userId)) + extractAnySuffix(localFilePath) // 请替换为实际的对象Key
	err = bucket.PutObjectFromFile(objectKey, localFilePath)
	if err != nil {
		log.Fatalf("Failed to put object from file: %v", err)
		return "", err
	}

	log.Println("File uploaded successfully.")

	//使用stringBuilder代替原本的直接拼接,优化性能
	//endpoint[:8] + bucketName + "." + endpoint[8:] + "/" + objectKey
	//len(endpoint[:8])+len(bucketName)+len(".")+len(endpoint[8:])+len("/")+len(objectKey),
	return accessUrl(endpoint, bucketName, endpoint), nil
}
