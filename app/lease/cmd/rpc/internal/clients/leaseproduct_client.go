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

// 请求响应结构定义（与leaseproduct服务的proto定义完全一致）
type CheckInventoryAvailabilityReq struct {
	ProductCode string `protobuf:"bytes,1,opt,name=productCode,proto3" json:"productCode,omitempty"`
	Quantity    int32  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	StartDate   string `protobuf:"bytes,3,opt,name=startDate,proto3" json:"startDate,omitempty"`
	EndDate     string `protobuf:"bytes,4,opt,name=endDate,proto3" json:"endDate,omitempty"`
}

// 修正：CheckInventoryAvailabilityResp应该与leaseproduct.proto保持一致，只包含available和availableCount字段
type CheckInventoryAvailabilityResp struct {
	Available      bool  `protobuf:"varint,1,opt,name=available,proto3" json:"available,omitempty"`
	AvailableCount int32 `protobuf:"varint,2,opt,name=availableCount,proto3" json:"availableCount,omitempty"`
}

type GetLeaseProductReq struct {
	ProductCode string `protobuf:"bytes,1,opt,name=productCode,proto3" json:"productCode,omitempty"`
}

// 修正：GetLeaseProductResp应该与leaseproduct.proto保持一致，只包含data字段
type GetLeaseProductResp struct {
	Data *LeaseProductInfo `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

type LeaseProductInfo struct {
	Id             int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductCode    string  `protobuf:"bytes,2,opt,name=productCode,proto3" json:"productCode,omitempty"`
	Name           string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Type           string  `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Machinery      string  `protobuf:"bytes,5,opt,name=machinery,proto3" json:"machinery,omitempty"`
	Brand          string  `protobuf:"bytes,6,opt,name=brand,proto3" json:"brand,omitempty"`
	Model          string  `protobuf:"bytes,7,opt,name=model,proto3" json:"model,omitempty"`
	DailyRate      float64 `protobuf:"fixed64,8,opt,name=dailyRate,proto3" json:"dailyRate,omitempty"`
	Deposit        float64 `protobuf:"fixed64,9,opt,name=deposit,proto3" json:"deposit,omitempty"`
	MaxDuration    int32   `protobuf:"varint,10,opt,name=maxDuration,proto3" json:"maxDuration,omitempty"`
	MinDuration    int32   `protobuf:"varint,11,opt,name=minDuration,proto3" json:"minDuration,omitempty"`
	Description    string  `protobuf:"bytes,12,opt,name=description,proto3" json:"description,omitempty"`
	InventoryCount int32   `protobuf:"varint,13,opt,name=inventoryCount,proto3" json:"inventoryCount,omitempty"`
	AvailableCount int32   `protobuf:"varint,14,opt,name=availableCount,proto3" json:"availableCount,omitempty"`
	Status         int32   `protobuf:"varint,15,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt      int64   `protobuf:"varint,16,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt      int64   `protobuf:"varint,17,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
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
