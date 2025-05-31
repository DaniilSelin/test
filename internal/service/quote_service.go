package service

import (
    "context"
    "time"

    "QuoteBook/internal/interfaces"
    "QuoteBook/internal/logger"
    "QuoteBook/internal/models"
    "QuoteBook/internal/errdefs"
    "QuoteBook/config"

    "go.uber.org/zap"
)

type QuoteService struct {
    repo interfaces.IQuoteRepository
    cfg *config.Config
}

func NewQuoteService(cfg *config.Config, repo interfaces.IQuoteRepository) QuoteService {
    return QuoteService{
        repo: repo, 
        cfg: cfg,
    }
}

func (qs QuoteService) CreateQuote(ctx context.Context, q *models.Quote) (int, error) {
    if q.Author== "" {
        return 0, errdefs.Wrap(errdefs.ErrInvalidInput, "Author reqiured")
    }
    return qs.repo.CreateQuote(ctx, q)
}

func QuotesAll(ctx context.Context) (*[]models.Quote, error) {
    return qs.repo.QuotesAll(ctx)
}

func QuoteByAuthor(ctx context.Context, author string) (*[]models.Quote, error) {
    return qs.repo.QuoteByAuthor(ctx, author)
}

func RandQuote(ctx context.Context) (*models.Quote, error) { 
    return qs.repo.RandQuote(ctx)
}

func DeleteQuote(ctx context.Context, id string) error {
    return qs.repo.DeleteQuote(ctx, id)
}