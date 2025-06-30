package product

import (
	"net/http"

	"api/internal/logic/product"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取贷款产品详情
func GetLoanProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLoanProductReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewGetLoanProductLogic(r.Context(), svcCtx)
		resp, err := l.GetLoanProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
