package loan

import (
	"net/http"

	"api/internal/logic/loan"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMyLoanApplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数中获取 applicationId
		applicationId := r.URL.Query().Get("id")
		if applicationId == "" {
			applicationId = r.PathValue("id") // go 1.22+ 方式
		}

		l := loan.NewGetMyLoanApplicationLogic(r.Context(), svcCtx)
		resp, err := l.GetMyLoanApplication(applicationId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
