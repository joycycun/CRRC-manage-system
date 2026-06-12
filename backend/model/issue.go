package model

type Issue struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`

	DeviceType  string `json:"deviceType"`
	IssueSource string `json:"issueSource"`
	Source      string `json:"source"`

	Level    string `json:"level"`
	Severity string `json:"severity"`

	IssueTitle string `json:"issueTitle"`
	Title      string `json:"title"`
	IssueDesc  string `json:"issueDesc"`

	OwnerID   int64  `json:"ownerId"`
	OwnerName string `json:"ownerName"`
	Owner     string `json:"owner"`

	CreatorID   int64  `json:"creatorId"`
	CreatorName string `json:"creatorName"`
	Creator     string `json:"creator"`

	CreateTime    string `json:"createTime"`
	PlanCloseTime string `json:"planCloseTime"`
	RealCloseTime string `json:"realCloseTime"`
	CloseStatus   string `json:"closeStatus"`

	CloseUserID   int64  `json:"closeUserId"`
	CloseUserName string `json:"closeUserName"`
	CloseUser     string `json:"closeUser"`

	ReopenReason string `json:"reopenReason"`
	UpdatedAt    string `json:"updatedAt"`

	Replies []IssueReply `json:"replies"`
}

type IssueReply struct {
	ID            int64  `json:"id"`
	IssueID       int64  `json:"issueId"`
	ReplyUserID   int64  `json:"replyUserId"`
	ReplyUserName string `json:"replyUserName"`
	ReplyUser     string `json:"replyUser"`
	ReplyTime     string `json:"replyTime"`
	Content       string `json:"content"`
	ReplyContent  string `json:"replyContent"`
}
