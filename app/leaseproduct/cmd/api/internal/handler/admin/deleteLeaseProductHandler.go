package admin

import (
	"net/http"

	"api/internal/logic/admin"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除租赁产品
func DeleteLeaseProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDeleteLeaseProductLogic(r.Context(), svcCtx)
		resp, err := l.DeleteLeaseProduct()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
