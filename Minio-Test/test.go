package Minio_Test

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func main() {
	// MinIO 配置信息
	endpoint := "localhost:1900"
	accessKeyID := "ZbK6pEyviB1rd1F54Laz"                         // 替换为您的 access key
	secretAccessKey := "9Cimnodjux8pN79GL4fX4VLSnj8HvL8Ymr4SQFjA" // 替换为您的 secret key
	useSSL := false

	// 初始化 MinIO 客户端
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("创建 MinIO 客户端失败:", err)
	}

	// 测试桶名称
	bucketName := "xaiohua"

	// 测试文件信息
	objectName := "test.txt"
	contentType := "text/plain"
	testContent := []byte("Hello, MinIO!")

	// 运行测试
	fmt.Println("开始 MinIO 连接测试...")

	// 1. 测试创建桶
	err = createBucketIfNotExists(minioClient, bucketName)
	if err != nil {
		log.Fatalln("创建桶失败:", err)
	}
	fmt.Println("✓ 桶创建/检查成功")

	// 2. 测试上传文件
	err = uploadTestFile(minioClient, bucketName, objectName, contentType, testContent)
	if err != nil {
		log.Fatalln("上传文件失败:", err)
	}
	fmt.Println("✓ 文件上传成功")

	// 3. 测试获取文件信息
	err = getFileInfo(minioClient, bucketName, objectName)
	if err != nil {
		log.Fatalln("获取文件信息失败:", err)
	}
	fmt.Println("✓ 获取文件信息成功")

	// 4. 测试下载文件
	err = downloadTestFile(minioClient, bucketName, objectName)
	if err != nil {
		log.Fatalln("下载文件失败:", err)
	}
	fmt.Println("✓ 文件下载成功")

	// 5. 测试删除文件
	err = deleteTestFile(minioClient, bucketName, objectName)
	if err != nil {
		log.Fatalln("删除文件失败:", err)
	}
	fmt.Println("✓ 文件删除成功")

	fmt.Println("\n所有测试完成！MinIO 配置正确且可以正常工作。")
}

// 如果桶不存在则创建
func createBucketIfNotExists(minioClient *minio.Client, bucketName string) error {
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("检查桶是否存在失败: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("创建桶失败: %v", err)
		}
		fmt.Printf("创建了新的桶: %s\n", bucketName)
	} else {
		fmt.Printf("桶已存在: %s\n", bucketName)
	}
	return nil
}

// 上传测试文件
func uploadTestFile(minioClient *minio.Client, bucketName, objectName, contentType string, content []byte) error {
	// 创建一个临时文件
	tmpfile, err := os.CreateTemp("", "minio-test-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// 写入测试内容
	if _, err := tmpfile.Write(content); err != nil {
		return fmt.Errorf("写入临时文件失败: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %v", err)
	}

	// 上传文件
	_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, tmpfile.Name(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}
	return nil
}

// 获取文件信息
func getFileInfo(minioClient *minio.Client, bucketName, objectName string) error {
	info, err := minioClient.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}
	fmt.Printf("文件信息: 大小=%d, 类型=%s\n", info.Size, info.ContentType)
	return nil
}

// 下载测试文件
func downloadTestFile(minioClient *minio.Client, bucketName, objectName string) error {
	tmpfile := "downloaded-test-HandleFile"
	err := minioClient.FGetObject(context.Background(), bucketName, objectName, tmpfile, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}
	defer os.Remove(tmpfile)
	return nil
}

// 删除测试文件
func deleteTestFile(minioClient *minio.Client, bucketName, objectName string) error {
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}
