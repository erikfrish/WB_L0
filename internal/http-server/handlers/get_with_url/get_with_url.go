package get_with_url

import (
	resp "L0/internal/http-server/api/response"
	order "L0/internal/strct"
	"L0/pkg/logger/sl"
	"context"
	"log/slog"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	OrderUID string `json:"order_uid" validate:"required"`
}

type Response struct {
	resp.Response
	order.Data
}

type OrderDataGetter interface {
	Get(ctx context.Context, order_uid string) (any, error)
}

func New(ctx context.Context, log *slog.Logger, orderDataGetter OrderDataGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		order_uid := chi.URLParam(r, "order_uid")

		data, err := orderDataGetter.Get(ctx, order_uid)

		if err != nil {
			log.Error("failed to get data/no such order with that id", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get data/no such order with that id"))
			return
		}
		log.Info("got data", slog.Any("data", data))

		responseOK(w, r, data.(order.Data))
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data order.Data) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Data:     data,
	})
}
