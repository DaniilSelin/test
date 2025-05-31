// фактически, эти тесты бесполезные, т.к. сложной логики нет
package service

import (
    "context"
    "errors"
    "log"
    "testing"
    "time"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"

    "QuoteBook/config"
    "QuoteBook/internal/errdefs"
    "QuoteBook/internal/models"
    "QuoteBook/internal/interfaces"
)

type MockQuoteRepository struct {
    mock.Mock
}

func (m *MockQuoteRepository) CreateQuote(ctx context.Context, q *models.Quote) (int, error) {
    args := m.Called(ctx, q)
    return args.Int(0), args.Error(1)
}

func (m *MockQuoteRepository) QuotesAll(ctx context.Context) ([]models.Quote, error) {
    args := m.Called(ctx)
    return args.Get(0).([]models.Quote), args.Error(1)
}

func (m *MockQuoteRepository) QuoteByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
    args := m.Called(ctx, author)
    return args.Get(0).([]models.Quote), args.Error(1)
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
    svc := service.NewQuoteService(cfg, mockRepo)

    q := &models.Quote{
        Author:    "Author1",
        Text:      "Sample text",
        CreatedAt: time.Now(),
    }

    // Expect repository CreateQuote to be called once with (ctx, q) and return id=123, nil
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
    svc := service.NewQuoteService(cfg, mockRepo)

    // Empty Author should return ErrInvalidInput without calling repo
    q := &models.Quote{
        Author:    "",
        Text:      "No author",
        CreatedAt: time.Now(),
    }

    id, err := svc.CreateQuote(ctx, q)
    require.Equal(t, 0, id)
    require.ErrorIs(t, err, errdefs.ErrInvalidInput)

    // Repo should not be called
    mockRepo.AssertNotCalled(t, "CreateQuote", mock.Anything, mock.Anything)
}

func TestQuotesAll_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := service.NewQuoteService(cfg, mockRepo)

    expected := []models.Quote{
        {ID: 1, Author: "A1", Text: "T1", CreatedAt: time.Now()},
        {ID: 2, Author: "A2", Text: "T2", CreatedAt: time.Now()},
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
    svc := service.NewQuoteService(cfg, mockRepo)

    mockRepo.On("QuotesAll", ctx).Return([]models.Quote{}, errdefs.ErrDB).Once()

    got, err := svc.QuotesAll(ctx)
    require.ErrorIs(t, err, errdefs.ErrDB)
    require.Nil(t, got)

    mockRepo.AssertExpectations(t)
}

func TestQuoteByAuthor_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := service.NewQuoteService(cfg, mockRepo)

    expected := []models.Quote{
        {ID: 1, Author: "AuthX", Text: "Tx", CreatedAt: time.Now()},
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
    svc := service.NewQuoteService(cfg, mockRepo)

    mockRepo.On("QuoteByAuthor", ctx, "NoAuth").Return([]models.Quote{}, nil).Once()

    got, err := svc.QuoteByAuthor(ctx, "NoAuth")
    require.NoError(t, err)
    require.Empty(t, got)

    mockRepo.AssertExpectations(t)
}

func TestRandQuote_Success(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := service.NewQuoteService(cfg, mockRepo)

    expected := &models.Quote{ID: 42, Author: "RAuthor", Text: "RText", CreatedAt: time.Now()}
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
    svc := service.NewQuoteService(cfg, mockRepo)

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
    svc := service.NewQuoteService(cfg, mockRepo)

    mockRepo.On("DeleteQuote", ctx, 99).Return(nil).Once()

    err := svc.DeleteQuote(ctx, 99)
    require.NoError(t, err)

    mockRepo.AssertExpectations(t)
}

func TestDeleteQuote_NotFound(t *testing.T) {
    ctx := context.Background()
    cfg := loadTestConfig(t)
    mockRepo := new(MockQuoteRepository)
    svc := service.NewQuoteService(cfg, mockRepo)

    mockRepo.On("DeleteQuote", ctx, 100).Return(errdefs.ErrNotFound).Once()

    err := svc.DeleteQuote(ctx, 100)
    require.ErrorIs(t, err, errdefs.ErrNotFound)

    mockRepo.AssertExpectations(t)
}