package admin

import (
	"net/http"

	"api/internal/logic/admin"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除贷款产品
func DeleteLoanProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDeleteLoanProductLogic(r.Context(), svcCtx)
		resp, err := l.DeleteLoanProduct()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
