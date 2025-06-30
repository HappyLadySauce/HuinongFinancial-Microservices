package clients

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// LeaseProduct服务客户端接口定义（基于leaseproduct.proto）
type LeaseProductClient interface {
	CheckInventoryAvailability(ctx context.Context, in *CheckInventoryAvailabilityReq, opts ...grpc.CallOption) (*CheckInventoryAvailabilityResp, error)
	GetLeaseProduct(ctx context.Context, in *GetLeaseProductReq, opts ...grpc.CallOption) (*GetLeaseProductResp, error)
}

// 请求响应结构定义（复制自leaseproduct服务的proto定义）
type CheckInventoryAvailabilityReq struct {
	ProductCode string `protobuf:"bytes,1,opt,name=productCode,proto3" json:"productCode,omitempty"`
	Quantity    int32  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	StartDate   string `protobuf:"bytes,3,opt,name=startDate,proto3" json:"startDate,omitempty"`
	EndDate     string `protobuf:"bytes,4,opt,name=endDate,proto3" json:"endDate,omitempty"`
}

type CheckInventoryAvailabilityResp struct {
	Code           int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message        string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Available      bool   `protobuf:"varint,3,opt,name=available,proto3" json:"available,omitempty"`
	AvailableCount int32  `protobuf:"varint,4,opt,name=availableCount,proto3" json:"availableCount,omitempty"`
}

type GetLeaseProductReq struct {
	ProductCode string `protobuf:"bytes,1,opt,name=productCode,proto3" json:"productCode,omitempty"`
}

type GetLeaseProductResp struct {
	Code    int32             `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string            `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *LeaseProductInfo `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

type LeaseProductInfo struct {
	Id          int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductCode string  `protobuf:"bytes,2,opt,name=productCode,proto3" json:"productCode,omitempty"`
	Name        string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	DailyRate   float64 `protobuf:"fixed64,8,opt,name=dailyRate,proto3" json:"dailyRate,omitempty"`
	Status      int32   `protobuf:"varint,15,opt,name=status,proto3" json:"status,omitempty"`
}

// LeaseProduct客户端实现
type leaseProductClient struct {
	cc zrpc.Client
}

// 创建LeaseProduct客户端
func NewLeaseProductClient(cc zrpc.Client) LeaseProductClient {
	return &leaseProductClient{cc}
}

func (c *leaseProductClient) CheckInventoryAvailability(ctx context.Context, in *CheckInventoryAvailabilityReq, opts ...grpc.CallOption) (*CheckInventoryAvailabilityResp, error) {
	out := new(CheckInventoryAvailabilityResp)
	err := c.cc.Conn().Invoke(ctx, "/leaseproduct.LeaseProductService/CheckInventoryAvailability", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaseProductClient) GetLeaseProduct(ctx context.Context, in *GetLeaseProductReq, opts ...grpc.CallOption) (*GetLeaseProductResp, error) {
	out := new(GetLeaseProductResp)
	err := c.cc.Conn().Invoke(ctx, "/leaseproduct.LeaseProductService/GetLeaseProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
