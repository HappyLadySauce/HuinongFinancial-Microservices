package admin

import (
	"net/http"

	"api/internal/logic/admin"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取租赁产品详情
func GetLeaseProductDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewGetLeaseProductDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetLeaseProductDetail()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
