package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	model "github.com/lin-snow/ech0/internal/model/common"
	userModel "github.com/lin-snow/ech0/internal/model/user"
)

type CommonServiceInterface interface {
	// CommonGetUserByUserId 根据用户ID获取用户信息
	CommonGetUserByUserId(userId uint) (userModel.User, error)

	// UploadImage 上传图片
	UploadImage(userid uint, file *multipart.FileHeader) (string, error)

	// DeleteImage 删除图片
	DeleteImage(userid uint, url, source string) error

	// DirectDeleteImage 直接根据URL和来源删除图片
	DirectDeleteImage(url, source string) error

	// GetSysAdmin 获取系统管理员
	GetSysAdmin() (userModel.User, error)

	// GetStatus 获取系统状态
	GetStatus() (model.Status, error)

	// GetHeatMap 获取热力图数据
	GetHeatMap() ([]model.Heatmap, error)

	// GenerateRSS 生成RSS订阅链接
	GenerateRSS(ctx *gin.Context) (string, error)

	// UploadMusic 上传音乐文件
	UploadMusic(userId uint, file *multipart.FileHeader) (string, error)

	// DeleteMusic 删除音乐文件
	DeleteMusic(userid uint) error

	// GetPlayMusicUrl 获取可播放的音乐URL
	GetPlayMusicUrl() string

	// PlayMusic 播放音乐
	PlayMusic(ctx *gin.Context)
}
