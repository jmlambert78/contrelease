package models

import "time"

type (
	// Release represents the structure of one new release for a given target
	Release struct {
		Id             string    `json:"_id" bson:"_id,omitempty"  form:"_id"`
		ServiceName    string    `json:"sname" bson:"sname"  form:"sname"`
		ServiceVersion string    `json:"svers" bson:"svers" form:"svers"`
		Destination    string    `json:"dest" bson:"dest"  form:"dest"`
		CentralZipURL  string    `json:"zipurl" bson:"zipurl"  form:"zipurl"`
		CentralImage   string    `json:"cimage" bson:"cimage"  form:"cimage"`
		LocalImage     string    `json:"limage" bson:"limage"  form:"limage"`
		InsertDate     time.Time `json:"date" bson:"date"  form:"date"`
		ReleaseStatus  bool      `json:"rmok" bson:"rmok"  form:"rmok"`
	}
)
