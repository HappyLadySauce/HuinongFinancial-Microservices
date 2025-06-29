package loan

import (
	"net/http"

	"api/internal/logic/loan"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMyLoanApplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := loan.NewGetMyLoanApplicationLogic(r.Context(), svcCtx)
		resp, err := l.GetMyLoanApplication()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
