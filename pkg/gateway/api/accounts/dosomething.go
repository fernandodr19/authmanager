package accounts

import (
	"net/http"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/fernandodr19/library/pkg/instrumentation"
	"github.com/google/uuid"
)

// Does something
// @Summary Does something
// @Description Does something incredible.
// @Tags Something
// @Param Authorization header string true "Bearer Authorization Token"
// @Accept json
// @Produce json
// @Success 200 {object} SomethingResponse
// @Header 200 {string} Token "X-Request-Id"
// @Failure 500 "Internal server error"
// @Router /do-something [get]
func (h Handler) DoSomething(r *http.Request) responses.Response {
	operation := "accounts.Handler.DoSomething"
	ctx := r.Context()
	log := instrumentation.Logger()
	log.Info("doing somethings")

	err := h.Usecase.DoSomething(ctx)
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	return responses.OK(SomethingResponse{
		SomethinID: uuid.NewString(),
	})
}

type SomethingResponse struct {
	SomethinID string `json:"something_id"`
}
