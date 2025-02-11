package HandleFile

import (
	"Gin-WebSocket/conf"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"path/filepath"
	"time"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
	Error   string `json:"error,omitempty"`
}
type UploadResponse struct {
	Status   int    `json:"status"`
	FileID   string `json:"fileId"`
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
}

func UploadFile(file multipart.File, header *multipart.FileHeader) (UploadResponse, error) {
	// 检查文件类型
	fileExt := filepath.Ext(header.Filename)
	if !isAllowedFileType(fileExt) {

		return UploadResponse{}, errors.New("不支持该类型文件！")

	}

	// 生成唯一文件名
	fileID := uuid.New().String()
	fileName := fileID + fileExt

	// 设置文件元数据
	metadata := map[string]string{
		"original-name": header.Filename,
		//"uploaded-by":   c.GetString("userId"), // 假设您在认证中间件中设置了userId
	}

	// 上传文件到MinIO
	if _, err := conf.MinioClient.PutObject(
		context.Background(),
		conf.MinioBucketName,
		fileName,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType:  header.Header.Get("Content-Type"),
			UserMetadata: metadata,
		},
	); err != nil {
		return UploadResponse{}, errors.New("上传文件到MinIO失败！")
	}

	// 生成预签名URL（7天有效）
	url, err := conf.MinioClient.PresignedGetObject(
		context.Background(),
		conf.MinioBucketName,
		fileName,
		time.Hour*24*7,
		nil,
	)
	if err != nil {
		return UploadResponse{}, errors.New("生成预签名URL出错！")
	}

	return UploadResponse{
		Status:   200,
		FileID:   fileID,
		FileName: header.Filename,
		FileURL:  url.String(),
		FileSize: header.Size,
		FileType: header.Header.Get("Content-Type"),
	}, nil

}

// 检查文件类型是否允许
func isAllowedFileType(fileType string) bool {
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".txt":  true,
		".zip":  true,
		".rar":  true,
		".mp4":  true,
		// 添加其他允许的文件类型
	}
	return allowedTypes[fileType]
}
