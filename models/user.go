package models

type (
	// User represents the structure of one user
	User struct {
		Name string `json:"name" bson:"name"  form:"name"`
		Pass string `json:"pass" bson:"pass" form:"pass"`
	}
)
type (
	// User represents the structure of one user
	Roles struct {
		UserName     string `json:"name" bson:"name"  form:"name"`
		DC           string `json:"dc" bson:"dc" form:"dc"`
		DevRole      bool   `json:"devrole" bson:"devrole"  form:"devrole"`
		RmRole       bool   `json:"rmrole" bson:"rmrole"  form:"rmrole"`
		OpsRole      bool   `json:"opsrole" bson:"opsrole"  form:"opsrole"`
		AdmRole      bool   `json:"admrole" bson:"admrole"  form:"admrole"`
		SuperAdmRole bool   `json:"superadmrole" bson:"superadmrole"  form:"superadmrole"`
	}
)
