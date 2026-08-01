package usecase

import (
	"context"
	"testing"

	"backend/internal/domain"
	"backend/internal/dto"

	"gorm.io/gorm"
)

type mockProductRepository struct {
	create   func(ctx context.Context, product *domain.Product) error
	findAll  func(ctx context.Context) ([]domain.Product, error)
	findByID func(ctx context.Context, id uint64) (*domain.Product, error)
	update   func(ctx context.Context, product *domain.Product) error
	delete   func(ctx context.Context, id uint64) error
}

func (m *mockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return m.create(ctx, product)
}

func (m *mockProductRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	return m.findAll(ctx)
}

func (m *mockProductRepository) FindByID(ctx context.Context, id uint64) (*domain.Product, error) {
	return m.findByID(ctx, id)
}

func (m *mockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return m.update(ctx, product)
}

func (m *mockProductRepository) Delete(ctx context.Context, id uint64) error {
	return m.delete(ctx, id)
}

func TestProductUsecase_Create_Success(t *testing.T) {
	repo := &mockProductRepository{
		create: func(ctx context.Context, product *domain.Product) error {
			product.ID = 1
			return nil
		},
	}

	uc := NewProductUsecase(repo)

	res, err := uc.Create(context.Background(), dto.CreateProductRequest{
		Name:        "Laptop",
		Description: "Gaming",
		Price:       10000000,
		Stock:       10,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != 1 {
		t.Errorf("expected ID 1, got %d", res.ID)
	}

	if res.Name != "Laptop" {
		t.Errorf("expected Laptop, got %s", res.Name)
	}
}

func TestProductUsecase_FindAll_Success(t *testing.T) {
	repo := &mockProductRepository{
		findAll: func(ctx context.Context) ([]domain.Product, error) {
			return []domain.Product{
				{
					ID:    1,
					Name:  "Laptop",
					Price: 100,
					Stock: 5,
				},
				{
					ID:    2,
					Name:  "Mouse",
					Price: 20,
					Stock: 30,
				},
			}, nil
		},
	}

	uc := NewProductUsecase(repo)

	products, err := uc.FindAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
}

func TestProductUsecase_FindByID_Success(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uint64) (*domain.Product, error) {
			return &domain.Product{
				ID:    id,
				Name:  "Laptop",
				Price: 100,
				Stock: 10,
			}, nil
		},
	}

	uc := NewProductUsecase(repo)

	product, err := uc.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.ID != 1 {
		t.Errorf("expected id 1")
	}
}

func TestProductUsecase_FindByID_NotFound(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uint64) (*domain.Product, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	uc := NewProductUsecase(repo)

	_, err := uc.FindByID(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "product not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProductUsecase_Update_Success(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uint64) (*domain.Product, error) {
			return &domain.Product{
				ID:    id,
				Name:  "Old",
				Price: 10,
				Stock: 1,
			}, nil
		},
		update: func(ctx context.Context, product *domain.Product) error {
			return nil
		},
	}

	uc := NewProductUsecase(repo)

	res, err := uc.Update(context.Background(), 1, dto.UpdateProductRequest{
		Name:        "New",
		Description: "Updated",
		Price:       200,
		Stock:       5,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Name != "New" {
		t.Errorf("expected New, got %s", res.Name)
	}
}

func TestProductUsecase_Update_NotFound(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uint64) (*domain.Product, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	uc := NewProductUsecase(repo)

	_, err := uc.Update(context.Background(), 1, dto.UpdateProductRequest{})

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "product not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProductUsecase_Delete_Success(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uint64) (*domain.Product, error) {
			return &domain.Product{
				ID: id,
			}, nil
		},
		delete: func(ctx context.Context, id uint64) error {
			return nil
		},
	}

	uc := NewProductUsecase(repo)

	if err := uc.Delete(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProductUsecase_Delete_NotFound(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uint64) (*domain.Product, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	uc := NewProductUsecase(repo)

	err := uc.Delete(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "product not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}
