package clients

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// LoanProduct服务客户端接口定义（基于loanproduct.proto）
type LoanProductClient interface {
	GetLoanProduct(ctx context.Context, in *GetLoanProductReq, opts ...grpc.CallOption) (*GetLoanProductResp, error)
}

// 请求响应结构定义（与loanproduct服务的proto定义完全一致）
type GetLoanProductReq struct {
	Id          int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductCode string `protobuf:"bytes,2,opt,name=productCode,proto3" json:"productCode,omitempty"`
}

// 修正：GetLoanProductResp应该与loanproduct.proto保持一致，只包含data字段
type GetLoanProductResp struct {
	Data *LoanProductInfo `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

type LoanProductInfo struct {
	Id           int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductCode  string  `protobuf:"bytes,2,opt,name=productCode,proto3" json:"productCode,omitempty"`
	Name         string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Type         string  `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	MaxAmount    float64 `protobuf:"fixed64,5,opt,name=maxAmount,proto3" json:"maxAmount,omitempty"`
	MinAmount    float64 `protobuf:"fixed64,6,opt,name=minAmount,proto3" json:"minAmount,omitempty"`
	MaxDuration  int32   `protobuf:"varint,7,opt,name=maxDuration,proto3" json:"maxDuration,omitempty"`
	MinDuration  int32   `protobuf:"varint,8,opt,name=minDuration,proto3" json:"minDuration,omitempty"`
	InterestRate float64 `protobuf:"fixed64,9,opt,name=interestRate,proto3" json:"interestRate,omitempty"`
	Description  string  `protobuf:"bytes,10,opt,name=description,proto3" json:"description,omitempty"`
	Status       int32   `protobuf:"varint,11,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt    int64   `protobuf:"varint,12,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt    int64   `protobuf:"varint,13,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

// LoanProduct客户端实现
type loanProductClient struct {
	cc zrpc.Client
}

// 创建LoanProduct客户端
func NewLoanProductClient(cc zrpc.Client) LoanProductClient {
	return &loanProductClient{cc}
}

func (c *loanProductClient) GetLoanProduct(ctx context.Context, in *GetLoanProductReq, opts ...grpc.CallOption) (*GetLoanProductResp, error) {
	out := new(GetLoanProductResp)
	err := c.cc.Conn().Invoke(ctx, "/loanproduct.LoanProductService/GetLoanProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
