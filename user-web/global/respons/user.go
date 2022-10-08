package respons

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32
	CreateeAt *time.Time
	UpdateAt  *time.Time
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

//在这里定义的时候发现的问题，在这里写备注的时候，我们需要注意的是，我们comment后面没有冒号，直接一个单引号就可以
type UserInfo struct {
	ID       int32     `json:"ID"`
	Mobile   string    `json:"Mobile,omitempty"`
	Password string    `json:"Password,omitempty"`
	NickName string    `json:"NickName,omitempty"`
	Birthday time.Time `json:"Birthday,omitempty"`
	Gender   string    `json:"Gender,omitempty"`
	Role     int       `json:"Role,omitempty"`
}

type UserLists []UserInfo

type GetUserInfoListRespon struct {
	Code     int       `json:"Code,omitempty"`
	Desc     string    `json:"Desc,omitempty"`
	UserList UserLists `json:"UserList"`
}
