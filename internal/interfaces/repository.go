package interfaces

import (
	"context"

	"quotebook/internal/models"
)

type IQuoteRepository interface {
    CreateQuote(ctx context.Context, q *models.Quote) (int, error)
    QuotesAll(ctx context.Context) (*[]models.Quote, error)
    QuoteByAuthor(ctx context.Context, author string) (*[]models.Quote, error)
    RandQuote(ctx context.Context) (*models.Quote, error)
    DeleteQuote(ctx context.Context, id int) error
}