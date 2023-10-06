package get

import (
	"context"
	"errors"

	// "errors"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "L0/internal/http-server/api/response"
	"L0/internal/storage"

	// "L0/internal/storage"
	order "L0/internal/strct"
	"L0/pkg/logger/sl"
)

type Request struct {
	OrderUID string `json:"order_uid" validate:"required"`
}

type Response struct {
	resp.Response
	order.Data
}

// cache or db rep entity
type OrderDataGetter interface {
	Get(ctx context.Context, order_uid string) (any, error)
}

// make handler
func New(ctx context.Context, log *slog.Logger, orderDataGetter OrderDataGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		req := Request{OrderUID: r.URL.Query().Get("order_uid")}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			// render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		data, err := orderDataGetter.Get(ctx, req.OrderUID)
		if errors.Is(err, storage.OrderNotFound) {
			log.Info("this", slog.String("order_uid", req.OrderUID))
			render.JSON(w, r, resp.Error("can't find such an order"))
			return
		}
		if err != nil {
			log.Error("failed to get data", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get data"))
			return
		}
		log.Debug("got data from cache", slog.Any("data", data))

		responseOK(w, r, data.(order.Data))
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data order.Data) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Data:     data,
	})
}
