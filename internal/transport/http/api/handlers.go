package api

import (
    "quotebook/config"
	"quotebook/internal/models"
	"quotebook/internal/errdefs"
	"quotebook/internal/logger"
    "quotebook/internal/interfaces"

    "fmt"
	"context"
    "strconv"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
)

type Handler struct {
	logger *logger.Logger
    cfg *config.Config
	qbs interfaces.IQuoteService
}

func NewHandler(lg *logger.Logger, cfg *config.Config, qbs interfaces.IQuoteService) *Handler {
	return &Handler{
		qbs: qbs,
		logger: lg,
        cfg: cfg,
	}
}

// удобно
func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// так как frontend фактически нет, то я тут генерирую RequstID
func (h *Handler) GenerateRequestID(r *http.Request) context.Context {
    ctx := r.Context()
    ctx = logger.CtxWWithLogger(ctx, h.logger)
    return context.WithValue(ctx, logger.RequestID, uuid.New().String())
}

// HandlePostQuote обрабатывает POST /quotes
func (h *Handler) HandlePostQuote() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := h.GenerateRequestID(r)

        h.logger.Info(ctx, "incoming request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )

        payload, err := decode[models.Quote](r)
        if err != nil {
            h.logger.Info(ctx, "invalid JSON payload", zap.Error(err))
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }
        id, err := h.qbs.CreateQuote(ctx, &payload);
        if err != nil {
            handleServiceError(ctx, w, err)
            return
        }

        h.logger.Info(ctx, "Quote created",
            zap.Int("id", id),
        )

        encode(w, r, http.StatusCreated, id)
    })
}

// HandleGetQuotes обрабатывает GET /quotes
func (h *Handler) HandleGetQuotes() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := h.GenerateRequestID(r)

        h.logger.Info(ctx, "incoming request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )

        quotes, err := h.qbs.QuotesAll(ctx)
        if err != nil {
            handleServiceError(ctx, w, err)
            return
        }

        h.logger.Info(ctx, "listed quotes",
            zap.Int("returned", len(*quotes)),
        )
        encode(w, r, http.StatusOK, quotes)
    })
}

// HandleGetQuoteByAuthor обрабатывает GET /quotes/{author}
func (h *Handler) HandleGetQuoteByAuthor() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := h.GenerateRequestID(r)

        h.logger.Info(ctx, "incoming request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )

        vars := mux.Vars(r)
        author := vars["author"]
        quotes, err := h.qbs.QuoteByAuthor(ctx, author)
        if err != nil {
            handleServiceError(ctx, w, err)
            return
        }

        h.logger.Info(ctx, "return quotes", zap.String("author", author))
        encode(w, r, http.StatusOK, quotes)
    })
}

// HandleGetRrote обрабатывает GET /quotes/random
func (h *Handler) HandleGetRandQuote() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := h.GenerateRequestID(r)

        h.logger.Info(ctx, "incoming request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )

        quote, err := h.qbs.RandQuote(ctx)
        if err != nil {
            handleServiceError(ctx, w, err)
            return
        }

        h.logger.Info(ctx, "return random quote",
            zap.Int("id", quote.ID),
        )
        encode(w, r, http.StatusOK, quote)
    })
}

// HandleDeleteQuote обрабатывает DELETE /quotes
func (h *Handler) HandleDeleteQuote() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := h.GenerateRequestID(r)

        h.logger.Info(ctx, "incoming request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )

        vars := mux.Vars(r)
        idStr := vars["id"]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            handleServiceError(ctx, w, errdefs.ErrInvalidInput)
            return
        }

        if err := h.qbs.DeleteQuote(ctx, id); err != nil {
            handleServiceError(ctx, w, err)
            return
        }


        h.logger.Info(ctx, "quote delete",
            zap.Int("id", id),
        )
        w.WriteHeader(http.StatusNoContent)
    })
}

// возвращает нужную ошибку
// чуть медленее чем на месте (много лишних проверок)
// зато код более компактный и читаемый
func handleServiceError(ctx context.Context, w http.ResponseWriter, err error) {
    switch {
    case errdefs.Is(err, errdefs.ErrNotFound):
        http.Error(w, "Not Found", http.StatusNotFound)
    case errdefs.Is(err, errdefs.ErrInvalidInput):
        http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
    case errdefs.Is(err, errdefs.ErrConflict):
        http.Error(w, "Conflict: "+err.Error(), http.StatusConflict)
    default:
        logger.GetLoggerFromCtx(ctx).Error(ctx, "internal error", zap.Error(err))
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}