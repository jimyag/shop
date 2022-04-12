package response

//
// GetUserInfoResponse
//  @Description: 用户信息的响应参数
//
type GetUserInfoResponse struct {
	Id        int32  `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	Role      int32  `json:"role"`
}
