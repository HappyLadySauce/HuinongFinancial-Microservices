package admin

import (
	"net/http"

	"api/internal/logic/admin"
	"api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetLoanApplicationDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewGetLoanApplicationDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetLoanApplicationDetail()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
