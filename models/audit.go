package models

import "time"

type (
	// Org representation
	AuditTrail struct {
		Id        string    `json:"_id" bson:"_id,omitempty"  form:"_id"`
		Who       string    `json:"who" bson:"who"  form:"who"`
		EventDate time.Time `json:"adate" bson:"adate"  form:"adate"`
		Action    string    `json:"action" bson:"action"  form:"action"`
		Release   Release   `json:"release" bson:"release"  form:"release"`
		Status    bool      `json:"status" bson:"status"  form:"status"`
	}
)
