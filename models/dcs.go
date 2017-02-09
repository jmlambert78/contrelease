package models

type (
	// Org representation
	Datacenter struct {
		Id         string `json:"_id" bson:"_id,omitempty"  form:"_id"`
		Name       string `json:"name" bson:"name"  form:"name"`
		LongName   string `json:"lname" bson:"lname"  form:"lname"`
		RegUrl     string `json:"regurl" bson:"regurl"  form:"regurl"`
		ZipRepoUrl string `json:"ziprepo" bson:"ziprepo" form:"ziprepo"`
	}
)
