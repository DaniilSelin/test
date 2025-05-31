package interfaces

import (
    "context"

    "quotebook/internal/models"
)

type IQuoteService interface {
    CreateQuote(ctx context.Context, b *models.Quote) (int, error)
    QuotesAll(ctx context.Context) ([]models.Quote, error)
    QuoteByAuthor(ctx context.Context, author string)
    RandQuote(ctx context.Context)
    DeleteQuote(ctx context.Context, id int)
}