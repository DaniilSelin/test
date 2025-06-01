package repository

import (
	"context"

	"quotebook/config"
	"quotebook/internal/errdefs"
	"quotebook/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuoteRepository struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewQuoteRepository(db *pgxpool.Pool, cfg *config.Config) QuoteRepository {
	return QuoteRepository{
		db:  db,
		cfg: cfg,
	}
}

func (qr QuoteRepository) CreateQuote(ctx context.Context, q *models.Quote) (int, error) {
	query := `
 		INSERT INTO quotesbook (
 			author, quote
 		) VALUES ($1, $2)
 		RETURNING id
	`
	var id int
    err := qr.db.QueryRow(ctx, query,
        q.Author,
        q.Quote,
    ).Scan(&id)
    if err != nil {
        return 0, errdefs.Wrapf(errdefs.ErrDB, "failed to create quote: %v", err)
    }
	return id, nil
}

func (qr QuoteRepository) QuotesAll(ctx context.Context) (*[]models.Quote, error) {
	query := `
		SELECT id, author, quote, created_at
		FROM quotesbook
	`
	rows, err := qr.db.Query(ctx, query)
	if err != nil {
		return nil, errdefs.Wrapf(errdefs.ErrDB, "failed to list quotes: %v", err)
	}
	defer rows.Close()

	var quotes []models.Quote
    for rows.Next() {
        var quote models.Quote
        if err := rows.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.CreatedAt); err != nil {
            return nil, errdefs.Wrapf(errdefs.ErrDB, "failed to scan quote: %v", err)
        }
        quotes = append(quotes, quote)
    }

	if rows.Err() != nil {
		return nil, errdefs.Wrapf(errdefs.ErrDB, "rows iteration error: %v", rows.Err())
	}

	return &quotes, nil
}

func (qr QuoteRepository) QuoteByAuthor(ctx context.Context, author string) (*[]models.Quote, error) {
    query := `
        SELECT id, author, quote, created_at
        FROM quotesbook
        WHERE author = $1
    `

	rows, err := qr.db.Query(ctx, query, author)
    if err != nil {
        return nil, errdefs.Wrapf(errdefs.ErrDB, "failed to query quotes by author: %v", err)
    }
    defer rows.Close()

	var quotes []models.Quote
    for rows.Next() {
        var quote models.Quote
        if err := rows.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.CreatedAt); err != nil {
            return nil, errdefs.Wrapf(errdefs.ErrDB, "failed to scan quote: %v", err)
        }
        quotes = append(quotes, quote)
    }

    if rows.Err() != nil {
        return nil, errdefs.Wrapf(errdefs.ErrDB, "rows iteration error: %v", rows.Err())
    }

	return &quotes, nil
}

func (qr QuoteRepository) RandQuote(ctx context.Context) (*models.Quote, error) {
	query := `
        SELECT id, author, quote, created_at
        FROM quotesbook
        ORDER BY RANDOM()
        LIMIT 1
    `

    var quote models.Quote
    err := qr.db.QueryRow(ctx, query).Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.CreatedAt)
    if err != nil {
        if errdefs.Is(err, pgx.ErrNoRows) {
		    return nil, errdefs.ErrNotFound
		}
        return nil, errdefs.Wrapf(errdefs.ErrDB, "failed to fetch random quote: %v", err)
    }
    return &quote, nil
}

func (qr QuoteRepository) DeleteQuote(ctx context.Context, id int) error {
    query := `
        DELETE FROM quotesbook
        WHERE id = $1
    `

	tag, err := qr.db.Exec(ctx, query, id)
    if err != nil {
        return errdefs.Wrapf(errdefs.ErrDB, "failed to delete quote %d: %v", id, err)
    }

    if tag.RowsAffected() == 0 {
        return errdefs.ErrNotFound
    }
    return nil
}