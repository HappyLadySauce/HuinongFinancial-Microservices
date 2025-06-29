package admin

import (
	"net/http"

	"api/internal/logic/admin"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取贷款产品详情
func GetLoanProductDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewGetLoanProductDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetLoanProductDetail()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
