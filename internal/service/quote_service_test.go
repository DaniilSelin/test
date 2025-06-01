// фактически, эти тесты бесполезные, т.к. сложной логики нет
package service

import (
    "context"
    "log"
    "testing"
    "time"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "quotebook/config"
    "quotebook/internal/errdefs"
    "quotebook/internal/models"
)

type MockQuoteRepository struct {
    mock.Mock
}

func (m *MockQuoteRepository) CreateQuote(ctx context.Context, q *models.Quote) (int, error) {
    args := m.Called(ctx, q)
    return args.Int(0), args.Error(1)
}

func (m *MockQuoteRepository) QuotesAll(ctx context.Context) (*[]models.Quote, error) {
    args := m.Called(ctx)
    return args.Get(0).(*[]models.Quote), args.Error(1)
}

func (m *MockQuoteRepository) QuoteByAuthor(ctx context.Context, author string) (*[]models.Quote, error) {
    args := m.Called(ctx, author)
    return args.Get(0).(*[]models.Quote), args.Error(1)
}

func (m *MockQuoteRepository) RandQuote(ctx context.Context) (*models.Quote, error) {
    args := m.Called(ctx)
    return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *MockQuoteRepository) DeleteQuote(ctx context.Context, id int) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func loadTestConfig(t *testing.T) *config.Config {
    cfg, err := config.LoadConfig("../../config/config.yml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    return cfg
}

func TestCreateQuote_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    q := &models.Quote{
        Author:    "Author1",
        Quote:      "Sample text",
        CreatedAt: time.Now(),
    }

    mockRepo.On("CreateQuote", ctx, q).Return(123, nil).Once()

    id, err := svc.CreateQuote(ctx, q)
    require.NoError(t, err)
    require.Equal(t, 123, id)

    mockRepo.AssertExpectations(t)
}

func TestCreateQuote_InvalidInput(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    q := &models.Quote{
        Author:    "",
        Quote:      "No author",
        CreatedAt: time.Now(),
    }

    id, err := svc.CreateQuote(ctx, q)
    require.Equal(t, 0, id)
    require.ErrorIs(t, err, errdefs.ErrInvalidInput)

    mockRepo.AssertNotCalled(t, "CreateQuote", mock.Anything, mock.Anything)
}

func TestQuotesAll_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    expected := &[]models.Quote{
        {ID: 1, Author: "A1", Quote: "T1", CreatedAt: time.Now()},
        {ID: 2, Author: "A2", Quote: "T2", CreatedAt: time.Now()},
    }

    mockRepo.On("QuotesAll", ctx).Return(expected, nil).Once()

    got, err := svc.QuotesAll(ctx)
    require.NoError(t, err)
    require.Equal(t, expected, got)

    mockRepo.AssertExpectations(t)
}

func TestQuotesAll_Error(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    // тут требуется конкретный nil
    mockRepo.On("QuotesAll", ctx).Return((*[]models.Quote)(nil), errdefs.ErrDB).Once()

    got, err := svc.QuotesAll(ctx)
    require.ErrorIs(t, err, errdefs.ErrDB)
    require.Nil(t, got)

    mockRepo.AssertExpectations(t)
}

func TestQuoteByAuthor_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    expected := &[]models.Quote{
        {ID: 1, Author: "AuthX", Quote: "Tx", CreatedAt: time.Now()},
    }

    mockRepo.On("QuoteByAuthor", ctx, "AuthX").Return(expected, nil).Once()

    got, err := svc.QuoteByAuthor(ctx, "AuthX")
    require.NoError(t, err)
    require.Equal(t, expected, got)

    mockRepo.AssertExpectations(t)
}

func TestQuoteByAuthor_EmptyResult(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    // та же самая ситуация
    mockRepo.On("QuoteByAuthor", ctx, "NoAuth").Return((*[]models.Quote)(nil), nil).Once()

    got, err := svc.QuoteByAuthor(ctx, "NoAuth")
    require.NoError(t, err)
    require.Empty(t, got)

    mockRepo.AssertExpectations(t)
}

func TestRandQuote_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    expected := &models.Quote{ID: 42, Author: "RAuthor", Quote: "RText", CreatedAt: time.Now()}
    mockRepo.On("RandQuote", ctx).Return(expected, nil).Once()

    got, err := svc.RandQuote(ctx)
    require.NoError(t, err)
    require.Equal(t, expected, got)

    mockRepo.AssertExpectations(t)
}

func TestRandQuote_NotFound(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    mockRepo.On("RandQuote", ctx).Return((*models.Quote)(nil), errdefs.ErrNotFound).Once()

    got, err := svc.RandQuote(ctx)
    require.ErrorIs(t, err, errdefs.ErrNotFound)
    require.Nil(t, got)

    mockRepo.AssertExpectations(t)
}

func TestDeleteQuote_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    mockRepo.On("DeleteQuote", ctx, 99).Return(nil).Once()

    err := svc.DeleteQuote(ctx, 99)
    require.NoError(t, err)

    mockRepo.AssertExpectations(t)
}

func TestDeleteQuote_NotFound(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := NewQuoteService(cfg, mockRepo)

    mockRepo.On("DeleteQuote", ctx, 100).Return(errdefs.ErrNotFound).Once()

    err := svc.DeleteQuote(ctx, 100)
    require.ErrorIs(t, err, errdefs.ErrNotFound)

    mockRepo.AssertExpectations(t)
}