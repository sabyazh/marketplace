package product

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/prabakarankannan/marketplace-backend/internal/common/pagination"
	"github.com/prabakarankannan/marketplace-backend/internal/common/response"
	"github.com/prabakarankannan/marketplace-backend/internal/common/validator"
	"github.com/prabakarankannan/marketplace-backend/internal/middleware"
)

// ProductHandler exposes HTTP endpoints for the product module.
type ProductHandler struct {
	service ProductService
}

// NewProductHandler returns a new ProductHandler wired to the given service.
func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// RegisterRoutes mounts all product and category routes onto the Fiber application.
func RegisterRoutes(app fiber.Router, h *ProductHandler, jwtSecret string) {
	// -----------------------------------------------------------------
	// Public routes (no authentication required)
	// -----------------------------------------------------------------
	products := app.Group("/products")
	products.Get("/", h.ListProducts)
	products.Get("/vendor/:vendorId", h.ListByVendor)
	products.Get("/:id", h.GetProduct)

	categories := app.Group("/categories")
	categories.Get("/", h.ListCategories)
	categories.Get("/:slug/products", h.GetCategoryProducts)

	// -----------------------------------------------------------------
	// Vendor routes (require authentication + vendor role)
	// -----------------------------------------------------------------
	vendorProducts := app.Group("/vendors/me/products", middleware.Auth(jwtSecret), middleware.RBAC("vendor", "admin"))
	vendorProducts.Post("/", h.CreateProduct)
	vendorProducts.Get("/", h.ListMyProducts)
	vendorProducts.Put("/:id", h.UpdateProduct)
	vendorProducts.Delete("/:id", h.DeleteProduct)
	vendorProducts.Put("/:id/stock", h.UpdateStock)

	vendorInventory := app.Group("/vendors/me/inventory", middleware.Auth(jwtSecret), middleware.RBAC("vendor", "admin"))
	vendorInventory.Get("/", h.GetInventory)
	vendorInventory.Get("/low-stock", h.GetLowStockAlerts)

	// -----------------------------------------------------------------
	// Admin routes (require authentication + admin role)
	// -----------------------------------------------------------------
	adminCategories := app.Group("/admin/categories", middleware.Auth(jwtSecret), middleware.RBAC("admin"))
	adminCategories.Post("/", h.CreateCategory)
	adminCategories.Put("/:id", h.UpdateCategory)
	adminCategories.Delete("/:id", h.DeleteCategory)

	adminProducts := app.Group("/admin/products", middleware.Auth(jwtSecret), middleware.RBAC("admin"))
	adminProducts.Get("/", h.AdminListProducts)
}

// ---------------------------------------------------------------------------
// Public product endpoints
// ---------------------------------------------------------------------------

// ListProducts handles GET /products.
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	params := pagination.FromQuery(c)
	filters := parseProductFilters(c)

	result, err := h.service.ListProducts(c.Context(), params, filters)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, result.Items, result.ToMeta(params))
}

// GetProduct handles GET /products/:id.
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product ID")
	}

	result, err := h.service.GetProduct(c.Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, result)
}

// ListByVendor handles GET /products/vendor/:vendorId.
func (h *ProductHandler) ListByVendor(c *fiber.Ctx) error {
	vendorID, err := uuid.Parse(c.Params("vendorId"))
	if err != nil {
		return response.BadRequest(c, "invalid vendor ID")
	}

	params := pagination.FromQuery(c)

	result, err := h.service.ListByVendor(c.Context(), vendorID, params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, result.Items, result.ToMeta(params))
}

// ---------------------------------------------------------------------------
// Public category endpoints
// ---------------------------------------------------------------------------

// ListCategories handles GET /categories.
func (h *ProductHandler) ListCategories(c *fiber.Ctx) error {
	categoryType := c.Query("type")

	result, err := h.service.ListCategories(c.Context(), categoryType)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, result)
}

// GetCategoryProducts handles GET /categories/:slug/products.
func (h *ProductHandler) GetCategoryProducts(c *fiber.Ctx) error {
	slug := c.Params("slug")
	params := pagination.FromQuery(c)

	result, err := h.service.ListByCategory(c.Context(), slug, params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, result.Items, result.ToMeta(params))
}

// ---------------------------------------------------------------------------
// Vendor product endpoints
// ---------------------------------------------------------------------------

// CreateProduct handles POST /vendors/me/products.
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return response.Unauthorized(c, "authentication required")
	}

	req, err := validator.ParseAndValidate[CreateProductRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	result, err := h.service.CreateProduct(c.Context(), userID, *req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, result)
}

// UpdateProduct handles PUT /vendors/me/products/:id.
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return response.Unauthorized(c, "authentication required")
	}

	productID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product ID")
	}

	req, parseErr := validator.ParseAndValidate[UpdateProductRequest](c)
	if parseErr != nil {
		return response.Error(c, parseErr)
	}

	if err := h.service.UpdateProduct(c.Context(), userID, productID, *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// DeleteProduct handles DELETE /vendors/me/products/:id.
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return response.Unauthorized(c, "authentication required")
	}

	productID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product ID")
	}

	if err := h.service.DeleteProduct(c.Context(), userID, productID); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// UpdateStock handles PUT /vendors/me/products/:id/stock.
func (h *ProductHandler) UpdateStock(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return response.Unauthorized(c, "authentication required")
	}

	productID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product ID")
	}

	req, parseErr := validator.ParseAndValidate[UpdateStockRequest](c)
	if parseErr != nil {
		return response.Error(c, parseErr)
	}

	if err := h.service.UpdateStock(c.Context(), userID, productID, req.StockQuantity); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// ListMyProducts handles GET /vendors/me/products.
func (h *ProductHandler) ListMyProducts(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return response.Unauthorized(c, "authentication required")
	}

	params := pagination.FromQuery(c)

	result, err := h.service.GetInventory(c.Context(), userID, params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, result.Items, result.ToMeta(params))
}

// GetInventory handles GET /vendors/me/inventory.
func (h *ProductHandler) GetInventory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return response.Unauthorized(c, "authentication required")
	}

	params := pagination.FromQuery(c)

	result, err := h.service.GetInventory(c.Context(), userID, params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, result.Items, result.ToMeta(params))
}

// GetLowStockAlerts handles GET /vendors/me/inventory/low-stock.
func (h *ProductHandler) GetLowStockAlerts(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return response.Unauthorized(c, "authentication required")
	}

	params := pagination.FromQuery(c)

	result, err := h.service.GetLowStockAlerts(c.Context(), userID, params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, result.Items, result.ToMeta(params))
}

// ---------------------------------------------------------------------------
// Admin endpoints
// ---------------------------------------------------------------------------

// CreateCategory handles POST /admin/categories.
func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	req, err := validator.ParseAndValidate[CreateCategoryRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	result, err := h.service.CreateCategory(c.Context(), *req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, result)
}

// UpdateCategory handles PUT /admin/categories/:id.
func (h *ProductHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid category ID")
	}

	req, parseErr := validator.ParseAndValidate[UpdateCategoryRequest](c)
	if parseErr != nil {
		return response.Error(c, parseErr)
	}

	if err := h.service.UpdateCategory(c.Context(), id, *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// DeleteCategory handles DELETE /admin/categories/:id.
func (h *ProductHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid category ID")
	}

	if err := h.service.DeleteCategory(c.Context(), id); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// AdminListProducts handles GET /admin/products.
func (h *ProductHandler) AdminListProducts(c *fiber.Ctx) error {
	params := pagination.FromQuery(c)
	filters := parseProductFilters(c)

	result, err := h.service.ListProducts(c.Context(), params, filters)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, result.Items, result.ToMeta(params))
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseProductFilters extracts ProductListFilters from the request query string.
func parseProductFilters(c *fiber.Ctx) ProductListFilters {
	var filters ProductListFilters

	if catID := c.Query("category_id"); catID != "" {
		if parsed, err := uuid.Parse(catID); err == nil {
			filters.CategoryID = &parsed
		}
	}

	if venID := c.Query("vendor_id"); venID != "" {
		if parsed, err := uuid.Parse(venID); err == nil {
			filters.VendorID = &parsed
		}
	}

	if minP := c.Query("min_price"); minP != "" {
		if parsed, err := strconv.ParseFloat(minP, 64); err == nil {
			filters.MinPrice = &parsed
		}
	}

	if maxP := c.Query("max_price"); maxP != "" {
		if parsed, err := strconv.ParseFloat(maxP, 64); err == nil {
			filters.MaxPrice = &parsed
		}
	}

	if active := c.Query("is_active"); active != "" {
		if parsed, err := strconv.ParseBool(active); err == nil {
			filters.IsActive = &parsed
		}
	}

	filters.Search = c.Query("search")
	filters.SortBy = c.Query("sort_by")
	filters.SortOrder = c.Query("sort_order")

	return filters
}
