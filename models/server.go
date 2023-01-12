package models

import (
	"time"

	"gorm.io/gorm"
)

var InstanceTypes = map[string]string{
	"t2.nano":    "1核0.5G",
	"t2.micro":   "1核1G",
	"t2.small":   "1核2G",
	"t2.medium":  "2核4G",
	"t2.large":   "2核8G",
	"t2.xlarge":  "4核16G",
	"t2.2xlarge": "8核32G",
	"t3.nano":    "2核0.5G",
	"t3.micro":   "2核1G",
	"t3.small":   "2核2G",
	"t3.medium":  "2核4G",
	"t3.large":   "2核8G",
	"t3.xlarge":  "4核16G",
	"t3.2xlarge": "8核32G",
	"c4.large":   "2核3.75G",
	"c4.xlarge":  "4核7.5G",
	"c4.2xlarge": "8核15G",
	"c4.4xlarge": "16核30G",
	"c4.8xlarge": "36核60G",
	"c5.large":   "2核4G",
	"c5.xlarge":  "4核8G",
	"c5.2xlarge": "8核16G",
	"c5.4xlarge": "16核32G",
}

type Server struct {
	ID             uint      `json:"id"`
	Name           string    `gorm:"size:128" json:"name" validate:"required"` // 机器名称或主机名
	IP             string    `gorm:"size:128;unique" json:"ip" validate:"required,ip"`
	Port           uint      `gorm:"default:22" json:"port" validate:"required,gte=0,lte=65535"`
	CredentialID   uint      `json:"credential_id" validate:"required"` // 关联认证
	Remark         string    `gorm:"type:text" json:"remark"`           // 记录机器用途
	InstanceType   string    `gorm:"size:128" json:"instance_type"`     // 实例类型
	Specifications string    `gorm:"size:128" json:"specifications"`    // 实例配置
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Servers []Server

func (s *Server) BeforeCreate(tx *gorm.DB) (err error) {
	if len(s.InstanceType) > 0 {
		s.Specifications = InstanceTypes[s.InstanceType]
	}
	return nil
}

func (s *Server) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(s.InstanceType) > 0 {
		s.Specifications = InstanceTypes[s.InstanceType]
	}
	return nil
}
