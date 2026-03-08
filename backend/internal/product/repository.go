package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/prabakarankannan/marketplace-backend/internal/common/pagination"
)

// ProductRepository defines the data-access contract for the product module.
type ProductRepository interface {
	// Category operations
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategory(ctx context.Context, id uuid.UUID) (*Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*Category, error)
	ListCategories(ctx context.Context, categoryType string) ([]Category, error)

	// Product operations
	CreateProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	ListProducts(ctx context.Context, params pagination.Params, filters ProductListFilters) (*pagination.Result[Product], error)
	ListByVendor(ctx context.Context, vendorID uuid.UUID, params pagination.Params) (*pagination.Result[Product], error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID, params pagination.Params) (*pagination.Result[Product], error)

	// Stock operations
	UpdateStock(ctx context.Context, productID uuid.UUID, quantity int) error
	GetLowStockProducts(ctx context.Context, vendorID uuid.UUID, params pagination.Params) (*pagination.Result[Product], error)
	IncrementStock(ctx context.Context, productID uuid.UUID, delta int) error
	DecrementStock(ctx context.Context, productID uuid.UUID, delta int) error
}

// productRepository is the GORM-backed implementation of ProductRepository.
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository returns a new ProductRepository backed by the provided GORM DB.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// ---------------------------------------------------------------------------
// Category operations
// ---------------------------------------------------------------------------

// CreateCategory inserts a new category into the database.
func (r *productRepository) CreateCategory(ctx context.Context, category *Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now
	return r.db.WithContext(ctx).Create(category).Error
}

// UpdateCategory saves all fields of the category back to the database.
func (r *productRepository) UpdateCategory(ctx context.Context, category *Category) error {
	category.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(category).Error
}

// DeleteCategory hard-deletes a category by its ID.
func (r *productRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Category{}).Error
}

// GetCategory retrieves a category by its primary key.
func (r *productRepository) GetCategory(ctx context.Context, id uuid.UUID) (*Category, error) {
	var cat Category
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// GetCategoryBySlug retrieves a category by its slug.
func (r *productRepository) GetCategoryBySlug(ctx context.Context, slug string) (*Category, error) {
	var cat Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// ListCategories returns a tree of categories (parents with preloaded children),
// optionally filtered by category type.
func (r *productRepository) ListCategories(ctx context.Context, categoryType string) ([]Category, error) {
	query := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, name ASC")
		}).
		Order("sort_order ASC, name ASC")

	if categoryType != "" {
		query = query.Where("category_type = ?", categoryType)
	}

	var categories []Category
	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// ---------------------------------------------------------------------------
// Product operations
// ---------------------------------------------------------------------------

// CreateProduct inserts a new product into the database.
func (r *productRepository) CreateProduct(ctx context.Context, product *Product) error {
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	return r.db.WithContext(ctx).Create(product).Error
}

// UpdateProduct saves all fields of the product back to the database.
func (r *productRepository) UpdateProduct(ctx context.Context, product *Product) error {
	product.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(product).Error
}

// DeleteProduct soft-deletes a product by its ID.
func (r *productRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Product{}).Error
}

// GetProduct retrieves a product by its primary key, preloading the Category.
func (r *productRepository) GetProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	var p Product
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Vendor").
		Where("id = ?", id).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// ListProducts returns a paginated list of products with dynamic filtering.
func (r *productRepository) ListProducts(ctx context.Context, params pagination.Params, filters ProductListFilters) (*pagination.Result[Product], error) {
	query := r.db.WithContext(ctx).Model(&Product{})

	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.VendorID != nil {
		query = query.Where("vendor_id = ?", *filters.VendorID)
	}
	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ?", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Determine sort order.
	orderClause := buildOrderClause(filters.SortBy, filters.SortOrder)

	var products []Product
	if err := query.
		Preload("Category").
		Preload("Vendor").
		Order(orderClause).
		Offset(params.Offset()).
		Limit(params.PerPage).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return &pagination.Result[Product]{
		Items: products,
		Total: total,
	}, nil
}

// ListByVendor returns a paginated list of products belonging to a vendor.
func (r *productRepository) ListByVendor(ctx context.Context, vendorID uuid.UUID, params pagination.Params) (*pagination.Result[Product], error) {
	query := r.db.WithContext(ctx).Model(&Product{}).Where("vendor_id = ?", vendorID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var products []Product
	if err := query.
		Preload("Category").
		Preload("Vendor").
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PerPage).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return &pagination.Result[Product]{
		Items: products,
		Total: total,
	}, nil
}

// ListByCategory returns a paginated list of active products in a category.
func (r *productRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID, params pagination.Params) (*pagination.Result[Product], error) {
	query := r.db.WithContext(ctx).Model(&Product{}).
		Where("category_id = ? AND is_active = true", categoryID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var products []Product
	if err := query.
		Preload("Category").
		Preload("Vendor").
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PerPage).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return &pagination.Result[Product]{
		Items: products,
		Total: total,
	}, nil
}

// ---------------------------------------------------------------------------
// Stock operations
// ---------------------------------------------------------------------------

// UpdateStock sets the stock quantity for a product.
func (r *productRepository) UpdateStock(ctx context.Context, productID uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"stock_quantity": quantity,
			"updated_at":    time.Now(),
		}).Error
}

// GetLowStockProducts returns products where stock_quantity <= low_stock_threshold
// for a given vendor.
func (r *productRepository) GetLowStockProducts(ctx context.Context, vendorID uuid.UUID, params pagination.Params) (*pagination.Result[Product], error) {
	query := r.db.WithContext(ctx).Model(&Product{}).
		Where("vendor_id = ? AND stock_quantity <= low_stock_threshold", vendorID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var products []Product
	if err := query.
		Preload("Category").
		Preload("Vendor").
		Order("stock_quantity ASC").
		Offset(params.Offset()).
		Limit(params.PerPage).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return &pagination.Result[Product]{
		Items: products,
		Total: total,
	}, nil
}

// IncrementStock increases the stock quantity by delta.
func (r *productRepository) IncrementStock(ctx context.Context, productID uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).
		Model(&Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"stock_quantity": gorm.Expr("stock_quantity + ?", delta),
			"updated_at":    time.Now(),
		}).Error
}

// DecrementStock decreases the stock quantity by delta. It returns an error
// if the current stock is less than delta to prevent negative stock.
func (r *productRepository) DecrementStock(ctx context.Context, productID uuid.UUID, delta int) error {
	result := r.db.WithContext(ctx).
		Model(&Product{}).
		Where("id = ? AND stock_quantity >= ?", productID, delta).
		Updates(map[string]interface{}{
			"stock_quantity": gorm.Expr("stock_quantity - ?", delta),
			"updated_at":    time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient stock or product not found")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// buildOrderClause constructs a SQL ORDER BY clause from sort parameters.
func buildOrderClause(sortBy, sortOrder string) string {
	column := "created_at"
	switch sortBy {
	case "price":
		column = "price"
	case "rating":
		column = "avg_rating"
	case "created_at":
		column = "created_at"
	}

	direction := "DESC"
	if sortOrder == "asc" {
		direction = "ASC"
	}

	return column + " " + direction
}
