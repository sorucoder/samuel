package payloads

import "time"

type ViewAudit struct {
	Date       time.Time `form:"date" time_format:"2006-01-02" binding:"required"`
	Page       int       `form:"page" binding:"min=0"`
	Count      int       `form:"count" binding:"min=1"`
	Sort       string    `form:"sort" binding:"omitempty,oneof=description user"`
	Descending bool      `form:"descending"`
}
