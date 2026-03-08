package product

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Category represents a row in the categories table.
type Category struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name         string      `gorm:"type:varchar(255);not null" json:"name"`
	Slug         string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description  *string     `gorm:"type:text" json:"description,omitempty"`
	ImageURL     *string     `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	ParentID     *uuid.UUID  `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	SortOrder    int         `gorm:"default:0" json:"sort_order"`
	IsActive     bool        `gorm:"default:true" json:"is_active"`
	CategoryType string      `gorm:"type:varchar(50);not null;default:'product'" json:"category_type"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Children     []Category  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName overrides the default GORM table name.
func (Category) TableName() string {
	return "categories"
}

// ProductVendor is a lightweight representation of a vendor used for preloading
// within the product module. It maps to the same "vendors" table.
type ProductVendor struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid" json:"user_id"`
	BusinessName string    `gorm:"type:varchar(255)" json:"business_name"`
	LogoURL      *string   `gorm:"type:varchar(500)" json:"logo_url,omitempty"`
	VendorType   string    `gorm:"type:varchar(50)" json:"vendor_type"`
	Status       string    `gorm:"type:varchar(20)" json:"status"`
	City         string    `gorm:"type:varchar(100)" json:"city"`
	IsOnline     bool      `json:"is_online"`
}

// TableName tells GORM to use the vendors table for ProductVendor.
func (ProductVendor) TableName() string {
	return "vendors"
}

// Product represents a row in the products table.
type Product struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	VendorID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"vendor_id"`
	CategoryID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug              string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description       *string        `gorm:"type:text" json:"description,omitempty"`
	Price             float64        `gorm:"type:float8;not null" json:"price"`
	ComparePrice      *float64       `gorm:"type:float8" json:"compare_price,omitempty"`
	Unit              string         `gorm:"type:varchar(50);not null;default:'piece'" json:"unit"`
	SKU               *string        `gorm:"type:varchar(100)" json:"sku,omitempty"`
	Images            pq.StringArray `gorm:"type:text[]" json:"images,omitempty"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	AvgRating         float64        `gorm:"type:float8;default:0" json:"avg_rating"`
	TotalReviews      int            `gorm:"default:0" json:"total_reviews"`
	StockQuantity     int            `gorm:"not null;default:0" json:"stock_quantity"`
	LowStockThreshold int            `gorm:"not null;default:5" json:"low_stock_threshold"`
	WeightGrams       *int           `json:"weight_grams,omitempty"`
	Tags              pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	Category          Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Vendor            ProductVendor  `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
}

// TableName overrides the default GORM table name.
func (Product) TableName() string {
	return "products"
}
