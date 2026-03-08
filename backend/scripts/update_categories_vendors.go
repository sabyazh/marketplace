//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("🖼️  Updating Category Images & Vendor Logos...")
	fmt.Println("══════════════════════════════════════════════════════")

	// ─────────────────────────────────────────────────────────────
	// 1. Update Category Images
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n📂 Adding images to categories...")

	categoryImages := map[string]string{
		// Original categories
		"fruits-and-vegetables":   "https://images.unsplash.com/photo-1610832958506-aa56368176cf?w=400",
		"fresh-vegetables":        "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=400",
		"fresh-fruits":            "https://images.unsplash.com/photo-1619566636858-adf3ef46400b?w=400",
		"dairy-and-bread":         "https://images.unsplash.com/photo-1628088062854-d1870b4553da?w=400",
		"milk-and-curd":           "https://images.unsplash.com/photo-1563636619-e9143da7973b?w=400",
		"paneer-and-cheese":       "https://images.unsplash.com/photo-1631452180519-c014fe946bc7?w=400",
		"bread-and-eggs":          "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400",
		"rice-and-atta":           "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400",
		"basmati-rice":            "https://images.unsplash.com/photo-1536304993881-460e32f50420?w=400",
		"wheat-atta":              "https://images.unsplash.com/photo-1574323347407-f5e1ad6d020b?w=400",
		"other-grains":            "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400",
		"dals-and-pulses":         "https://images.unsplash.com/photo-1613585064236-185d1e628fb4?w=400",
		"masala-and-spices":       "https://images.unsplash.com/photo-1596040033229-a9821ebd058d?w=400",
		"cooking-oil-and-ghee":    "https://images.unsplash.com/photo-1474979266404-7eaacbcd87c5?w=400",
		"snacks-and-namkeen":      "https://images.unsplash.com/photo-1599490659213-e2b9527b711e?w=400",
		"beverages":               "https://images.unsplash.com/photo-1544145945-f90425340c7e?w=400",
		"dry-fruits-and-nuts":     "https://images.unsplash.com/photo-1605493725784-84a3570385b2?w=400",
		"instant-and-ready-to-eat": "https://images.unsplash.com/photo-1612929633738-8fe44f7ec841?w=400",
		"personal-care":           "https://images.unsplash.com/photo-1556228578-0d85b1a4d571?w=400",
		"cleaning-and-household":  "https://images.unsplash.com/photo-1585421514738-01798e348b17?w=400",

		// TN-specific categories
		"millets-and-traditional-grains": "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400",
		"ragi-products":                  "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400",
		"kambu-and-thinai":               "https://images.unsplash.com/photo-1574323347407-f5e1ad6d020b?w=400",
		"pickles-and-chutneys":           "https://images.unsplash.com/photo-1589135233689-3d3c8f012741?w=400",
		"podi-and-powder-mix":            "https://images.unsplash.com/photo-1596040033229-a9821ebd058d?w=400",
		"traditional-snacks":             "https://images.unsplash.com/photo-1599490659213-e2b9527b711e?w=400",
		"traditional-sweets":             "https://images.unsplash.com/photo-1589301760014-d929f3979dbc?w=400",
		"filter-coffee-and-tea":          "https://images.unsplash.com/photo-1514432324607-a09d9b4aefda?w=400",
		"papad-and-appalam":              "https://images.unsplash.com/photo-1601050690597-df0568f70950?w=400",
		"pooja-items":                    "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400",
		"jaggery-and-sweeteners":         "https://images.unsplash.com/photo-1604431696980-07e518647610?w=400",
		"cold-pressed-oils":              "https://images.unsplash.com/photo-1474979266404-7eaacbcd87c5?w=400",
		"idli-dosa-and-breakfast-mix":    "https://images.unsplash.com/photo-1589301760014-d929f3979dbc?w=400",
		"traditional-rice-varieties":     "https://images.unsplash.com/photo-1536304993881-460e32f50420?w=400",
	}

	catUpdated := 0
	for slug, imageURL := range categoryImages {
		result, err := db.Exec(
			"UPDATE categories SET image_url = $1 WHERE slug = $2 AND (image_url IS NULL OR image_url = '')",
			imageURL, slug,
		)
		if err != nil {
			log.Printf("  ✗ %s: %v", slug, err)
			continue
		}
		rows, _ := result.RowsAffected()
		if rows > 0 {
			catUpdated++
			fmt.Printf("  ✓ %s\n", slug)
		}
	}

	// ─────────────────────────────────────────────────────────────
	// 2. Update Vendor Logos
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n🏪 Adding logos to vendors...")

	vendorLogos := map[string]string{
		// Original vendors
		"Fresh Basket":           "https://images.unsplash.com/photo-1542838132-92c53300491e?w=200",
		"Annapurna Grocery":      "https://images.unsplash.com/photo-1604719312566-8912e9227c6a?w=200",
		"Spice Emporium":         "https://images.unsplash.com/photo-1596040033229-a9821ebd058d?w=200",
		"Desi Mart":              "https://images.unsplash.com/photo-1578916171728-46686eac8d58?w=200",
		"Organic Wala":           "https://images.unsplash.com/photo-1488459716781-31db52582fe9?w=200",

		// TN vendors
		"Nallennai Groceries":    "https://images.unsplash.com/photo-1474979266404-7eaacbcd87c5?w=200",
		"Chettinad Masala House": "https://images.unsplash.com/photo-1596040033229-a9821ebd058d?w=200",
		"Kumbakonam Degree Coffee": "https://images.unsplash.com/photo-1514432324607-a09d9b4aefda?w=200",
		"Tamil Millets Store":    "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=200",
		"Madurai Snacks & Sweets": "https://images.unsplash.com/photo-1589301760014-d929f3979dbc?w=200",
		"Kovai Organics":         "https://images.unsplash.com/photo-1488459716781-31db52582fe9?w=200",
		"Tirunelveli Treats":     "https://images.unsplash.com/photo-1604431696980-07e518647610?w=200",
	}

	vendorUpdated := 0
	for name, logoURL := range vendorLogos {
		result, err := db.Exec(
			"UPDATE vendors SET logo_url = $1 WHERE business_name = $2 AND (logo_url IS NULL OR logo_url = '')",
			logoURL, name,
		)
		if err != nil {
			log.Printf("  ✗ %s: %v", name, err)
			continue
		}
		rows, _ := result.RowsAffected()
		if rows > 0 {
			vendorUpdated++
			fmt.Printf("  ✓ %s\n", name)
		}
	}

	// ── Final stats ──────────────────────────────────────────────
	fmt.Println("\n══════════════════════════════════════════════════════")
	fmt.Printf("  ✅  %d categories updated with images\n", catUpdated)
	fmt.Printf("  ✅  %d vendors updated with logos\n", vendorUpdated)
	fmt.Println("══════════════════════════════════════════════════════")
}
