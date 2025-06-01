package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"quotebook/config"
	"quotebook/internal/database"
	"quotebook/internal/errdefs"
	"quotebook/internal/models"
)

var db *pgxpool.Pool
var cfg *config.Config

func TestMain(m *testing.M) {
	var err error

	ctx := context.Background()

	// Загружаем конфиг
	cfg, err = config.LoadConfig("../../config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к БД
	db, err = database.Connect(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	database.MigrationPath = "../../internal/database/migrations"
    if err := database.RunMigrations(ctx, cfg, db); err != nil {
        log.Fatalf("Failed run migrationы to database: %v", err)
    }

	code := m.Run()
	os.Exit(code)
}

func clearTable(t *testing.T) {
	ctx := context.Background()
	_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.quotesbook RESTART IDENTITY CASCADE", cfg.DB.Schema))
	require.NoError(t, err, "Failed to clear quotesbook table")
}

func TestQuoteRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewQuoteRepository(db, cfg)

	t.Run("CreateGetAllDelete", func(t *testing.T) {
		clearTable(t)

		quote := &models.Quote{
			Author:    "Author1",
			Quote:      "First quote text",
			CreatedAt: time.Now(),
		}

		// Create
		id, err := repo.CreateQuote(ctx, quote)
		require.NoError(t, err, "Error when creating quote")
		require.Greater(t, id, 0, "Expected positive ID after creation")

		// QuotesAll
		allQuotes, err := repo.QuotesAll(ctx)
		require.NoError(t, err, "Error when fetching all quotes")
		require.Len(t, *allQuotes, 1, "Expected exactly one quote in table")
		require.Equal(t, quote.Author, (*allQuotes)[0].Author)
		require.Equal(t, quote.Quote, (*allQuotes)[0].Quote)

		// Delete
		err = repo.DeleteQuote(ctx, id)
		require.NoError(t, err, "Error when deleting quote")

		allQuotesAfter, err := repo.QuotesAll(ctx)
		require.NoError(t, err)
		require.Len(t, *allQuotesAfter, 0, "Expected table to be empty after deletion")
	})

	t.Run("QuoteByAuthor", func(t *testing.T) {
		clearTable(t)

		quote1 := &models.Quote{
			Author:    "CommonAuthor",
			Quote:      "Text A",
			CreatedAt: time.Now(),
		}
		quote2 := &models.Quote{
			Author:    "CommonAuthor",
			Quote:      "Text B",
			CreatedAt: time.Now(),
		}
		quote3 := &models.Quote{
			Author:    "OtherAuthor",
			Quote:      "Text C",
			CreatedAt: time.Now(),
		}

		_, err := repo.CreateQuote(ctx, quote1)
		require.NoError(t, err)
		_, err = repo.CreateQuote(ctx, quote2)
		require.NoError(t, err)
		_, err = repo.CreateQuote(ctx, quote3)
		require.NoError(t, err)

		byAuthor, err := repo.QuoteByAuthor(ctx, "CommonAuthor")
		require.NoError(t, err, "Error when fetching quotes by author")
		require.Len(t, *byAuthor, 2, "Expected two quotes for author CommonAuthor")
		for _, q := range *byAuthor {
			require.Equal(t, "CommonAuthor", q.Author)
		}

		emptyBy, err := repo.QuoteByAuthor(ctx, "NoSuchAuthor")
		require.NoError(t, err)
		require.Len(t, *emptyBy, 0, "Expected empty result for non-existent author")
	})

	t.Run("RandQuoteEmpty", func(t *testing.T) {
		clearTable(t)

		// Пустая таблица → ErrNotFound
		_, err := repo.RandQuote(ctx)
		require.Error(t, err, "Expected an error when fetching random quote from empty table")
		require.Equal(t, errdefs.ErrNotFound, err, "Expected ErrNotFound")
	})

	t.Run("RandQuoteOne", func(t *testing.T) {
		clearTable(t)

		quote := &models.Quote{
			Author:    "SoloAuthor",
			Quote:      "Only quote",
			CreatedAt: time.Now(),
		}
		id, err := repo.CreateQuote(ctx, quote)
		require.NoError(t, err)
		require.Greater(t, id, 0)

		got, err := repo.RandQuote(ctx)
		require.NoError(t, err, "Error when fetching random quote")
		require.Equal(t, quote.Author, got.Author)
		require.Equal(t, quote.Quote, got.Quote)
	})

	t.Run("DeleteNotFound", func(t *testing.T) {
		clearTable(t)

		err := repo.DeleteQuote(ctx, 9999)
		require.Error(t, err, "Expected an error when deleting non-existent ID")
		require.Equal(t, errdefs.ErrNotFound, err, "Expected ErrNotFound")
	})
}
