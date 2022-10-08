package request

//获取用户列表参数
type GetUserListRequst struct {
	Pn    uint32 `json:"Pn,omitempty" form:"Pn"`
	PSize uint32 `json:"PSize,omitempty" form:"PSize"`
}
