package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex" validate:"required,min=2,max=100"`
	Description string         `json:"description" gorm:"type:text"`
	ParentID    *uint          `json:"parent_id,omitempty" gorm:"index"`
	Parent      *Category      `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []Category     `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Product represents a product
type Product struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null" validate:"required,min=2,max=200"`
	Description  string         `json:"description" gorm:"type:text"`
	SKU          string         `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Price        float64        `json:"price" gorm:"not null;check:price >= 0" validate:"required,min=0"`
	ComparePrice *float64       `json:"compare_price,omitempty" gorm:"check:compare_price >= 0"`
	CategoryID   uint           `json:"category_id" gorm:"not null;index" validate:"required"`
	Category     Category       `json:"category" gorm:"foreignKey:CategoryID"`
	Brand        string         `json:"brand,omitempty" validate:"max=100"`
	Weight       *float64       `json:"weight,omitempty" gorm:"check:weight >= 0"`
	Dimensions   string         `json:"dimensions,omitempty"`
	Stock        int            `json:"stock" gorm:"not null;default:0;check:stock >= 0" validate:"min=0"`
	MinStock     int            `json:"min_stock" gorm:"default:0;check:min_stock >= 0"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	IsFeatured   bool           `json:"is_featured" gorm:"default:false"`
	Tags         string         `json:"tags,omitempty"`
	ImageURL     string         `json:"image_url,omitempty" validate:"omitempty,url"`
	Images       []ProductImage `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID uint           `json:"product_id" gorm:"not null;index"`
	URL       string         `json:"url" gorm:"not null" validate:"required,url"`
	AltText   string         `json:"alt_text,omitempty"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ProductCreateRequest represents a product creation request
type ProductCreateRequest struct {
	Name         string   `json:"name" validate:"required,min=2,max=200"`
	Description  string   `json:"description,omitempty"`
	SKU          string   `json:"sku" validate:"required"`
	Price        float64  `json:"price" validate:"required,min=0"`
	ComparePrice *float64 `json:"compare_price,omitempty" validate:"omitempty,min=0"`
	CategoryID   uint     `json:"category_id" validate:"required"`
	Brand        string   `json:"brand,omitempty" validate:"max=100"`
	Weight       *float64 `json:"weight,omitempty" validate:"omitempty,min=0"`
	Dimensions   string   `json:"dimensions,omitempty"`
	Stock        int      `json:"stock" validate:"min=0"`
	MinStock     int      `json:"min_stock" validate:"min=0"`
	IsActive     bool     `json:"is_active"`
	IsFeatured   bool     `json:"is_featured"`
	Tags         string   `json:"tags,omitempty"`
	ImageURL     string   `json:"image_url,omitempty" validate:"omitempty,url"`
	Images       []string `json:"images,omitempty"`
}

// ProductUpdateRequest represents a product update request
type ProductUpdateRequest struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description  *string  `json:"description,omitempty"`
	SKU          *string  `json:"sku,omitempty"`
	Price        *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	ComparePrice *float64 `json:"compare_price,omitempty" validate:"omitempty,min=0"`
	CategoryID   *uint    `json:"category_id,omitempty"`
	Brand        *string  `json:"brand,omitempty" validate:"omitempty,max=100"`
	Weight       *float64 `json:"weight,omitempty" validate:"omitempty,min=0"`
	Dimensions   *string  `json:"dimensions,omitempty"`
	Stock        *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	MinStock     *int     `json:"min_stock,omitempty" validate:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active,omitempty"`
	IsFeatured   *bool    `json:"is_featured,omitempty"`
	Tags         *string  `json:"tags,omitempty"`
	ImageURL     *string  `json:"image_url,omitempty" validate:"omitempty,url"`
}

// ProductResponse represents public product data
type ProductResponse struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	SKU          string           `json:"sku"`
	Price        float64          `json:"price"`
	ComparePrice *float64         `json:"compare_price,omitempty"`
	Category     CategoryResponse `json:"category"`
	Brand        string           `json:"brand,omitempty"`
	Weight       *float64         `json:"weight,omitempty"`
	Dimensions   string           `json:"dimensions,omitempty"`
	Stock        int              `json:"stock"`
	MinStock     int              `json:"min_stock"`
	IsActive     bool             `json:"is_active"`
	IsFeatured   bool             `json:"is_featured"`
	Tags         string           `json:"tags,omitempty"`
	ImageURL     string           `json:"image_url,omitempty"`
	Images       []ProductImage   `json:"images,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// CategoryResponse represents public category data
type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// ToResponse converts Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		SKU:          p.SKU,
		Price:        p.Price,
		ComparePrice: p.ComparePrice,
		Category:     p.Category.ToResponse(),
		Brand:        p.Brand,
		Weight:       p.Weight,
		Dimensions:   p.Dimensions,
		Stock:        p.Stock,
		MinStock:     p.MinStock,
		IsActive:     p.IsActive,
		IsFeatured:   p.IsFeatured,
		Tags:         p.Tags,
		ImageURL:     p.ImageURL,
		Images:       p.Images,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

// ToResponse converts Category to CategoryResponse
func (c *Category) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ParentID:    c.ParentID,
		IsActive:    c.IsActive,
	}
}

// TableName returns the table name
func (Product) TableName() string {
	return "products"
}

// TableName returns the table name
func (Category) TableName() string {
	return "categories"
}

// TableName returns the table name
func (ProductImage) TableName() string {
	return "product_images"
}
