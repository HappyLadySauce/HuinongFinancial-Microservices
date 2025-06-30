package lease

import (
	"net/http"

	"api/internal/logic/lease"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMyLeaseApplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数中获取 applicationId
		applicationId := r.URL.Query().Get("id")
		if applicationId == "" {
			applicationId = r.PathValue("id") // go 1.22+ 方式
		}

		l := lease.NewGetMyLeaseApplicationLogic(r.Context(), svcCtx)
		resp, err := l.GetMyLeaseApplication(applicationId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
