package product

import (
	"net/http"

	"api/internal/logic/product"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取贷款产品详情
func GetLoanProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewGetLoanProductLogic(r.Context(), svcCtx)
		resp, err := l.GetLoanProduct()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
