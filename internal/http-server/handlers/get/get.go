package get

import (
	"context"
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "L0/internal/http-server/api/response"
	"L0/internal/storage"
	order "L0/internal/strct"
	"L0/pkg/logger/sl"
)

type Request struct {
	OrderUID string `json:"order_uid" validate:"required,url"`
	// Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	order.Data
	// Alias string `json:"alias,omitempty"`
}

type OrderDataGetter interface {
	Get(ctx context.Context, order_uid string) (order.Data, error)
}

func New(ctx context.Context, log *slog.Logger, orderDataGetter OrderDataGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			// render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		data, err := orderDataGetter.Get(ctx, req.OrderUID)
		if errors.Is(err, storage.OrderNotFound) {
			log.Info("url already exists", slog.String("url", req.OrderUID))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}
		log.Info("url added", slog.Any("data", data))

		responseOK(w, r, data)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data order.Data) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Data:     data,
		// Alias:    alias,
	})
}
