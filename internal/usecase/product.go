package usecase

import (
	"context"
	"errors"

	"backend/internal/domain"
	"backend/internal/dto"
	"backend/internal/repository"

	"gorm.io/gorm"
)

type ProductUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(
	productRepo repository.ProductRepository,
) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (u *ProductUsecase) Create(
	ctx context.Context,
	req dto.CreateProductRequest,
) (*dto.ProductResponse, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	response := dto.ToProductResponse(product)

	return &response, nil
}

func (u *ProductUsecase) FindAll(
	ctx context.Context,
) ([]dto.ProductResponse, error) {
	products, err := u.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToProductResponses(products), nil
}

func (u *ProductUsecase) FindByID(
	ctx context.Context,
	id uint64,
) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	response := dto.ToProductResponse(product)

	return &response, nil
}

func (u *ProductUsecase) Update(
	ctx context.Context,
	id uint64,
	req dto.UpdateProductRequest,
) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	response := dto.ToProductResponse(product)

	return &response, nil
}

func (u *ProductUsecase) Delete(
	ctx context.Context,
	id uint64,
) error {
	_, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return u.productRepo.Delete(ctx, id)
}
