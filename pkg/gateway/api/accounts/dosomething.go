package accounts

import (
	"net/http"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/fernandodr19/library/pkg/instrumentation"
)

func (h Handler) DoSomething(r *http.Request) responses.Response {
	operation := "accounts.Handler.DoSomething"
	ctx := r.Context()
	log := instrumentation.Logger()
	log.Info("doing somethings")

	err := h.Usecase.DoSomething(ctx)
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	return responses.OK(nil)
}
