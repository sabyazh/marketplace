package product

import (
	"time"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Category Request DTOs
// ---------------------------------------------------------------------------

// CreateCategoryRequest is the payload for creating a new category.
type CreateCategoryRequest struct {
	Name         string     `json:"name"          validate:"required"`
	Slug         string     `json:"slug"          validate:"required"`
	CategoryType string     `json:"category_type" validate:"required"`
	Description  *string    `json:"description,omitempty"`
	ImageURL     *string    `json:"image_url,omitempty"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder    *int       `json:"sort_order,omitempty"`
}

// UpdateCategoryRequest is the payload for updating an existing category.
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

// ---------------------------------------------------------------------------
// Category Response DTOs
// ---------------------------------------------------------------------------

// CategoryResponse is the public representation of a category.
type CategoryResponse struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name"`
	Slug         string             `json:"slug"`
	Description  *string            `json:"description,omitempty"`
	ImageURL     *string            `json:"image_url,omitempty"`
	ParentID     *uuid.UUID         `json:"parent_id,omitempty"`
	SortOrder    int                `json:"sort_order"`
	IsActive     bool               `json:"is_active"`
	CategoryType string             `json:"category_type"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Children     []CategoryResponse `json:"children,omitempty"`
}

// ---------------------------------------------------------------------------
// Product Request DTOs
// ---------------------------------------------------------------------------

// CreateProductRequest is the payload for creating a new product.
type CreateProductRequest struct {
	CategoryID        uuid.UUID `json:"category_id"        validate:"required"`
	Name              string    `json:"name"               validate:"required"`
	Slug              string    `json:"slug"               validate:"required"`
	Description       *string   `json:"description,omitempty"`
	Price             float64   `json:"price"              validate:"required,gt=0"`
	ComparePrice      *float64  `json:"compare_price,omitempty"`
	Unit              *string   `json:"unit,omitempty"`
	SKU               *string   `json:"sku,omitempty"`
	Images            []string  `json:"images,omitempty"`
	StockQuantity     int       `json:"stock_quantity"     validate:"required,gte=0"`
	LowStockThreshold *int      `json:"low_stock_threshold,omitempty"`
	WeightGrams       *int      `json:"weight_grams,omitempty"`
	Tags              []string  `json:"tags,omitempty"`
}

// UpdateProductRequest is the payload for updating an existing product.
type UpdateProductRequest struct {
	CategoryID        *uuid.UUID `json:"category_id,omitempty"`
	Name              *string    `json:"name,omitempty"`
	Slug              *string    `json:"slug,omitempty"`
	Description       *string    `json:"description,omitempty"`
	Price             *float64   `json:"price,omitempty"     validate:"omitempty,gt=0"`
	ComparePrice      *float64   `json:"compare_price,omitempty"`
	Unit              *string    `json:"unit,omitempty"`
	SKU               *string    `json:"sku,omitempty"`
	Images            []string   `json:"images,omitempty"`
	IsActive          *bool      `json:"is_active,omitempty"`
	StockQuantity     *int       `json:"stock_quantity,omitempty"     validate:"omitempty,gte=0"`
	LowStockThreshold *int       `json:"low_stock_threshold,omitempty"`
	WeightGrams       *int       `json:"weight_grams,omitempty"`
	Tags              []string   `json:"tags,omitempty"`
}

// UpdateStockRequest is the payload for directly setting stock quantity.
type UpdateStockRequest struct {
	StockQuantity int `json:"stock_quantity" validate:"required,gte=0"`
}

// ---------------------------------------------------------------------------
// Product Response DTOs
// ---------------------------------------------------------------------------

// ProductVendorResponse is a lightweight vendor summary embedded in product responses.
type ProductVendorResponse struct {
	ID           uuid.UUID `json:"id"`
	BusinessName string    `json:"business_name"`
	LogoURL      *string   `json:"logo_url,omitempty"`
	VendorType   string    `json:"vendor_type,omitempty"`
	Status       string    `json:"status,omitempty"`
	City         string    `json:"city,omitempty"`
	IsOnline     bool      `json:"is_online"`
}

// ProductCategoryResponse is a lightweight category summary embedded in product responses.
type ProductCategoryResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	CategoryType string    `json:"category_type,omitempty"`
	ImageURL     *string   `json:"image_url,omitempty"`
}

// ProductResponse is the public representation of a product.
type ProductResponse struct {
	ID                uuid.UUID                `json:"id"`
	VendorID          uuid.UUID                `json:"vendor_id"`
	CategoryID        uuid.UUID                `json:"category_id"`
	CategoryName      string                   `json:"category_name"`
	Name              string                   `json:"name"`
	Slug              string                   `json:"slug"`
	Description       *string                  `json:"description,omitempty"`
	Price             float64                  `json:"price"`
	ComparePrice      *float64                 `json:"compare_price,omitempty"`
	Unit              string                   `json:"unit"`
	SKU               *string                  `json:"sku,omitempty"`
	Images            []string                 `json:"images,omitempty"`
	IsActive          bool                     `json:"is_active"`
	AvgRating         float64                  `json:"avg_rating"`
	TotalReviews      int                      `json:"total_reviews"`
	StockQuantity     int                      `json:"stock_quantity"`
	LowStockThreshold int                      `json:"low_stock_threshold"`
	WeightGrams       *int                     `json:"weight_grams,omitempty"`
	Tags              []string                 `json:"tags,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	Vendor            *ProductVendorResponse   `json:"vendor,omitempty"`
	Category          *ProductCategoryResponse `json:"category,omitempty"`
}

// ---------------------------------------------------------------------------
// Product List Filters
// ---------------------------------------------------------------------------

// ProductListFilters holds the optional query parameters for listing products.
type ProductListFilters struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	VendorID   *uuid.UUID `json:"vendor_id,omitempty"`
	MinPrice   *float64   `json:"min_price,omitempty"`
	MaxPrice   *float64   `json:"max_price,omitempty"`
	IsActive   *bool      `json:"is_active,omitempty"`
	Search     string     `json:"search,omitempty"`
	SortBy     string     `json:"sort_by,omitempty"`
	SortOrder  string     `json:"sort_order,omitempty"`
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

// toCategoryResponse converts a Category model to its public response DTO.
func toCategoryResponse(c *Category) CategoryResponse {
	resp := CategoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		ImageURL:     c.ImageURL,
		ParentID:     c.ParentID,
		SortOrder:    c.SortOrder,
		IsActive:     c.IsActive,
		CategoryType: c.CategoryType,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}

	if len(c.Children) > 0 {
		resp.Children = make([]CategoryResponse, len(c.Children))
		for i := range c.Children {
			resp.Children[i] = toCategoryResponse(&c.Children[i])
		}
	}

	return resp
}

// toProductResponse converts a Product model to its public response DTO.
func toProductResponse(p *Product) ProductResponse {
	resp := ProductResponse{
		ID:                p.ID,
		VendorID:          p.VendorID,
		CategoryID:        p.CategoryID,
		CategoryName:      p.Category.Name,
		Name:              p.Name,
		Slug:              p.Slug,
		Description:       p.Description,
		Price:             p.Price,
		ComparePrice:      p.ComparePrice,
		Unit:              p.Unit,
		SKU:               p.SKU,
		Images:            p.Images,
		IsActive:          p.IsActive,
		AvgRating:         p.AvgRating,
		TotalReviews:      p.TotalReviews,
		StockQuantity:     p.StockQuantity,
		LowStockThreshold: p.LowStockThreshold,
		WeightGrams:       p.WeightGrams,
		Tags:              p.Tags,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}

	// Populate nested vendor if preloaded
	if p.Vendor.ID != uuid.Nil {
		resp.Vendor = &ProductVendorResponse{
			ID:           p.Vendor.ID,
			BusinessName: p.Vendor.BusinessName,
			LogoURL:      p.Vendor.LogoURL,
			VendorType:   p.Vendor.VendorType,
			Status:       p.Vendor.Status,
			City:         p.Vendor.City,
			IsOnline:     p.Vendor.IsOnline,
		}
	}

	// Populate nested category if preloaded
	if p.Category.ID != uuid.Nil {
		resp.Category = &ProductCategoryResponse{
			ID:           p.Category.ID,
			Name:         p.Category.Name,
			Slug:         p.Category.Slug,
			CategoryType: p.Category.CategoryType,
			ImageURL:     p.Category.ImageURL,
		}
	}

	return resp
}
