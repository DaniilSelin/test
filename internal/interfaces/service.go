package interfaces

import (
    "context"

    "quotebook/internal/models"
)

type IQuoteService interface {
    CreateQuote(ctx context.Context, b *models.Quote) (int, error)
    QuotesAll(ctx context.Context) (*[]models.Quote, error)
    QuoteByAuthor(ctx context.Context, author string) (*[]models.Quote, error)
    RandQuote(ctx context.Context) (*models.Quote, error)
    DeleteQuote(ctx context.Context, id int) error
}