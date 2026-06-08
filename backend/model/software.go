package model

import (
	"database/sql"
	"time"
)

type SoftwareVersion struct {
	ID              int64        `json:"id"`
	ProjectID       int64        `json:"projectId"`
	SoftwareVersion string       `json:"softwareVersion"`
	DeviceType      string       `json:"deviceType"`
	HardwareID      int64        `json:"hardwareId"`
	HardwareVersion string       `json:"hardwareVersion"`
	OwnerID         int64        `json:"ownerId"`
	OwnerName       string       `json:"ownerName"`
	ReleaseDate     sql.NullTime `json:"-"`
	DownloadURL     string       `json:"downloadUrl"`
	BusinessDesc    string       `json:"businessDesc"`
	Description     string       `json:"description"`
	SoftwareStatus  string       `json:"softwareStatus"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	IsDeleted       bool         `json:"isDeleted"`
}

type SoftwareVersionResponse struct {
	ID              int64      `json:"id"`
	ProjectID       int64      `json:"projectId"`
	SoftwareVersion string     `json:"softwareVersion"`
	DeviceType      string     `json:"deviceType"`
	HardwareID      int64      `json:"hardwareId"`
	HardwareVersion string     `json:"hardwareVersion"`
	OwnerID         int64      `json:"ownerId"`
	OwnerName       string     `json:"ownerName"`
	ReleaseDate     *time.Time `json:"releaseDate"`
	DownloadURL     string     `json:"downloadUrl"`
	BusinessDesc    string     `json:"businessDesc"`
	Description     string     `json:"description"`
	SoftwareStatus  string     `json:"softwareStatus"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	IsDeleted       bool       `json:"isDeleted"`
}

func SoftwareNullTimeToPtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func (s SoftwareVersion) ToResponse() SoftwareVersionResponse {
	return SoftwareVersionResponse{
		ID:              s.ID,
		ProjectID:       s.ProjectID,
		SoftwareVersion: s.SoftwareVersion,
		DeviceType:      s.DeviceType,
		HardwareID:      s.HardwareID,
		HardwareVersion: s.HardwareVersion,
		OwnerID:         s.OwnerID,
		OwnerName:       s.OwnerName,
		ReleaseDate:     SoftwareNullTimeToPtr(s.ReleaseDate),
		DownloadURL:     s.DownloadURL,
		BusinessDesc:    s.BusinessDesc,
		Description:     s.Description,
		SoftwareStatus:  s.SoftwareStatus,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		IsDeleted:       s.IsDeleted,
	}
}

type ProjectBranch struct {
	ID         int64     `json:"id"`
	ProjectID  int64     `json:"projectId"`
	DeviceType string    `json:"deviceType"`
	RepoName   string    `json:"repoName"`
	RepoURL    string    `json:"repoUrl"`
	BranchName string    `json:"branchName"`
	CloneURL   string    `json:"cloneUrl"`
	OwnerID    int64     `json:"ownerId"`
	OwnerName  string    `json:"ownerName"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	IsDeleted  bool      `json:"isDeleted"`
}
