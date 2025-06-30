package clients

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// AppUser服务客户端接口定义（基于appuser.proto）
type AppUserClient interface {
	GetUserById(ctx context.Context, in *GetUserByIdReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
}

// 请求响应结构定义（与appuser服务的proto定义完全一致）
type GetUserByIdReq struct {
	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

// 修正：GetUserInfoResp应该与appuser.proto保持一致，只包含UserInfo字段
type GetUserInfoResp struct {
	UserInfo *UserInfo `protobuf:"bytes,1,opt,name=user_info,json=userInfo,proto3" json:"user_info,omitempty"`
}

type UserInfo struct {
	Id         int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                 // 用户ID
	Phone      string  `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`                            // 手机号
	Name       string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`                              // 姓名
	Nickname   string  `protobuf:"bytes,4,opt,name=nickname,proto3" json:"nickname,omitempty"`                      // 昵称
	Age        int32   `protobuf:"varint,5,opt,name=age,proto3" json:"age,omitempty"`                               // 年龄
	Gender     int32   `protobuf:"varint,6,opt,name=gender,proto3" json:"gender,omitempty"`                         // 0:未知 1:男 2:女
	Occupation string  `protobuf:"bytes,7,opt,name=occupation,proto3" json:"occupation,omitempty"`                  // 职业
	Address    string  `protobuf:"bytes,8,opt,name=address,proto3" json:"address,omitempty"`                        // 地址
	Income     float64 `protobuf:"fixed64,9,opt,name=income,proto3" json:"income,omitempty"`                        // 收入 单位:元
	CreatedAt  int64   `protobuf:"varint,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"` // 创建时间
	UpdatedAt  int64   `protobuf:"varint,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"` // 更新时间
}

// AppUser客户端实现
type appUserClient struct {
	cc zrpc.Client
}

// 创建AppUser客户端
func NewAppUserClient(cc zrpc.Client) AppUserClient {
	return &appUserClient{cc}
}

func (c *appUserClient) GetUserById(ctx context.Context, in *GetUserByIdReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	out := new(GetUserInfoResp)
	err := c.cc.Conn().Invoke(ctx, "/appuser.AppUser/GetUserById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
