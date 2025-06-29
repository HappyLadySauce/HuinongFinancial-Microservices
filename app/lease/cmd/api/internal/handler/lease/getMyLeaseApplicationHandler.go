package lease

import (
	"net/http"

	"api/internal/logic/lease"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMyLeaseApplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := lease.NewGetMyLeaseApplicationLogic(r.Context(), svcCtx)
		resp, err := l.GetMyLeaseApplication()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
