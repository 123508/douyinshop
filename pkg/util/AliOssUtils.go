package util

import (
	"fmt"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
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
	const maxFileSize = 5 << 20 // 5MB

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭文件错误")
		}
	}(file)

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

// UploadImages 阿里OSS对象存储上传图片(本地上传)
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
	return accessUrl(endpoint, bucketName, objectKey), nil
}

// UploadImagesByIO 阿里OSS对象存储上传图片(通过网络流上传)
// 最大上传限制5MB
func UploadImagesByIO(file *multipart.FileHeader) (string, error) {

	reader, err := file.Open()

	if err != nil {
		log.Println("文件打开错误!")

		return "", &errorno.BasicMessageError{Message: "文件打开错误", Code: 400}
	}

	defer func(reader multipart.File) {
		err := reader.Close()
		if err != nil {
			log.Println("关闭文件错误")
		}
	}(reader)

	// 检查文件大小是否超过限制
	const maxFileSize = 5 << 20 // 5 MB

	if file.Size > maxFileSize {

		log.Println("文件过大")

		return "", &errorno.BasicMessageError{Message: "文件过大", Code: 400}
	}

	endpoint := config.Conf.AliyunConfig.Oss.Endpoint
	accessKeyID := config.Conf.AliyunConfig.Oss.AccessKeyId
	accessKeySecret := config.Conf.AliyunConfig.Oss.AccessKeySecret
	suffix := extractAnySuffix(file.Filename)

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

	s := uuid.New().String()

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	objectKey := s + suffix // 请替换为实际的对象Key
	err = bucket.PutObject(objectKey, reader)
	if err != nil {
		log.Fatalf("Failed to put object from file: %v", err)
		return "", err
	}

	log.Println("File uploaded successfully.")

	return accessUrl(endpoint, bucketName, objectKey), nil
}

// DownloadImages 将图片保存到服务器中
// 最大上传限制5MB
func DownloadImages(file *multipart.FileHeader) (string, error) {
	reader, err := file.Open()

	if err != nil {
		log.Println("文件打开错误!")

		return "", &errorno.BasicMessageError{Message: "文件打开错误", Code: 400}
	}

	defer func(reader multipart.File) {
		err := reader.Close()
		if err != nil {
			log.Println("文件关闭错误")
		}
	}(reader)

	// 检查文件大小是否超过限制
	const maxFileSize = 5 << 20 // 5 MB

	if file.Size > maxFileSize {

		log.Println("文件过大")

		return "", &errorno.BasicMessageError{Message: "文件过大", Code: 400}
	}

	// 创建目标文件

	s := uuid.New().String()

	path := "../../static/imageStore/" + s + extractAnySuffix(file.Filename)

	dir := filepath.Dir(path)

	err = os.MkdirAll(dir, 755)

	if err != nil {
		return "", err
	}

	dst, err := os.Create(path)
	if err != nil {
		return "", &errorno.BasicMessageError{Message: "无法创建目标文件,请重试", Code: 500}
	}

	abs, err := filepath.Abs(dst.Name())

	if err != nil {
		log.Println("获取文件绝对路径失败")
		return "", &errorno.BasicMessageError{Message: "获取文件绝对路径失败", Code: 404}
	}

	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Println("关闭文件错误")
		}
	}(dst)

	// 将上传的文件内容复制到目标文件
	if _, err := io.Copy(dst, reader); err != nil {

		return "", &errorno.BasicMessageError{Message: "无法复制目标文件,请重试", Code: 500}
	}
	return abs, nil

}
