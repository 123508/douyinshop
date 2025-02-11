package utils

import (
	"github.com/123508/douyinshop/pkg/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
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

// UploadImages 阿里OSS对象存储上传图片
func UploadImages(localFilePath string, serviceName string, userId uint32) string {

	endpoint := config.Conf.AliyunConfig.Oss.Endpoint
	accessKeyID := config.Conf.AliyunConfig.Oss.AccessKeyId
	accessKeySecret := config.Conf.AliyunConfig.Oss.AccessKeySecret

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Fatalf("Failed to create OSS client: %v", err)
	}

	// 填写存储空间名称，例如examplebucket。
	bucketName := config.Conf.AliyunConfig.Oss.BucketName // 请替换为实际的Bucket名称
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
	}

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	objectKey := serviceName + "_" + strconv.Itoa(int(userId)) + extractAnySuffix(localFilePath) // 请替换为实际的对象Key
	err = bucket.PutObjectFromFile(objectKey, localFilePath)
	if err != nil {
		log.Fatalf("Failed to put object from file: %v", err)
	}

	log.Println("File uploaded successfully.")

	return endpoint[:8] + bucketName + "." + endpoint[8:] + "/" + objectKey
}
