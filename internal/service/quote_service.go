package service

import (
    "context"

    "quotebook/internal/interfaces"
    _ "quotebook/internal/logger"
    "quotebook/internal/models"
    "quotebook/internal/errdefs"
    "quotebook/config"

    _ "go.uber.org/zap"
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

func (qs QuoteService) QuotesAll(ctx context.Context) (*[]models.Quote, error) {
    return qs.repo.QuotesAll(ctx)
}

func (qs QuoteService) QuoteByAuthor(ctx context.Context, author string) (*[]models.Quote, error) {
    return qs.repo.QuoteByAuthor(ctx, author)
}

func (qs QuoteService) RandQuote(ctx context.Context) (*models.Quote, error) { 
    return qs.repo.RandQuote(ctx)
}

func (qs QuoteService) DeleteQuote(ctx context.Context, id int) error {
    return qs.repo.DeleteQuote(ctx, id)
}