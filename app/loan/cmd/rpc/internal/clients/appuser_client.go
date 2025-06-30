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

// 请求响应结构定义（复制自appuser服务的proto定义）
type GetUserByIdReq struct {
	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type GetUserInfoResp struct {
	Code     int32     `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message  string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	UserInfo *UserInfo `protobuf:"bytes,3,opt,name=user_info,json=userInfo,proto3" json:"user_info,omitempty"`
}

type UserInfo struct {
	Id     int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Phone  string `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	Status int32  `protobuf:"varint,10,opt,name=status,proto3" json:"status,omitempty"`
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
