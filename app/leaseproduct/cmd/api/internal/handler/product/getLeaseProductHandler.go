package product

import (
	"net/http"

	"api/internal/logic/product"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取租赁产品详情
func GetLeaseProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewGetLeaseProductLogic(r.Context(), svcCtx)
		resp, err := l.GetLeaseProduct()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
