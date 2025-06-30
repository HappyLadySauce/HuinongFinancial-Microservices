package admin

import (
	"net/http"

	"api/internal/logic/admin"
	"api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetLeaseApplicationDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数中获取 applicationId
		applicationId := r.URL.Query().Get("id")
		if applicationId == "" {
			// 尝试从路径值获取（go 1.22+ 方式）
			applicationId = r.PathValue("id")
		}

		l := admin.NewGetLeaseApplicationDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetLeaseApplicationDetail(applicationId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
