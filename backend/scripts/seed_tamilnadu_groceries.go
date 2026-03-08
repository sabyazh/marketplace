//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println("🏛️  Seeding Tamil Nadu Grocery Data...")
	fmt.Println("══════════════════════════════════════════════════════")

	// ─────────────────────────────────────────────────────────────
	// 1. Create TN-specific Vendors
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n📦 Creating Tamil Nadu vendors...")

	passwordHash := "$2a$10$YRz8rU.yrG3RzFxqeGhB8eCwGZxH3Kx9ZLqK.Bz8RLhxN3P5a6Cq2" // Vendor@123

	tnVendors := []struct {
		fullName, email, phone, biz, desc, address, city, state, pincode string
		lat, lng                                                         float64
	}{
		{"Kumar Selvam", "kumar@naallannaigroceries.in", "9876500001", "Nallennai Groceries", "Traditional Tamil Nadu cold-pressed oils and groceries from Erode", "23 Perundurai Road, Erode", "Erode", "Tamil Nadu", "638001", 11.3410, 77.7172},
		{"Meena Lakshmi", "meena@chettinadsupplies.in", "9876500002", "Chettinad Masala House", "Authentic Chettinad spices and masalas from Karaikudi", "45 Palace Road, Karaikudi", "Karaikudi", "Tamil Nadu", "630001", 10.0736, 78.7839},
		{"Ravi Shankar", "ravi@kumbakonamcoffee.in", "9876500003", "Kumbakonam Degree Coffee", "Premium South Indian filter coffee and beverages", "12 Big Street, Kumbakonam", "Kumbakonam", "Tamil Nadu", "612001", 10.9617, 79.3881},
		{"Lakshmi Devi", "lakshmi@tamilmillets.in", "9876500004", "Tamil Millets Store", "Organic millets and traditional grains from Dindigul", "78 Anna Nagar, Dindigul", "Dindigul", "Tamil Nadu", "624001", 10.3673, 77.9803},
		{"Senthil Murugan", "senthil@madurasnacks.in", "9876500005", "Madurai Snacks & Sweets", "Traditional Tamil Nadu snacks, sweets and pickles", "56 Meenakshi Amman Koil Street, Madurai", "Madurai", "Tamil Nadu", "625001", 9.9252, 78.1198},
		{"Anbu Selvan", "anbu@kovaiorganics.in", "9876500006", "Kovai Organics", "Farm-fresh organic produce and traditional items from Coimbatore", "34 RS Puram, Coimbatore", "Coimbatore", "Tamil Nadu", "641002", 11.0168, 76.9558},
		{"Priya Dharshini", "priya@tirunelvelishop.in", "9876500007", "Tirunelveli Treats", "Famous Tirunelveli halwa, banana chips and palm products", "89 High Ground, Tirunelveli", "Tirunelveli", "Tamil Nadu", "627001", 8.7139, 77.7567},
	}

	vendorIDs := make(map[string]string) // business_name -> vendor_id

	for _, v := range tnVendors {
		// Create user (matching existing seed script schema)
		var userID string
		err := db.QueryRow(`
			INSERT INTO users (full_name, email, phone, password_hash, role, is_active)
			VALUES ($1, $2, $3, $4, 'vendor', true)
			ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name
			RETURNING id
		`, v.fullName, v.email, v.phone, passwordHash).Scan(&userID)
		if err != nil {
			log.Printf("  ⚠ User %s: %v", v.email, err)
			db.QueryRow("SELECT id FROM users WHERE email = $1", v.email).Scan(&userID)
		}

		// Create vendor (using ST_MakePoint for PostGIS location)
		var vendorID string
		err = db.QueryRow(`
			INSERT INTO vendors (user_id, business_name, description, vendor_type, status,
				location, address, city, state, pincode,
				service_radius_km, commission_pct, is_online, avg_rating, total_reviews)
			VALUES ($1, $2, $3, 'product', 'approved',
				ST_MakePoint($4, $5)::geography, $6, $7, $8, $9,
				15.0, 8.0, true, $10, $11)
			ON CONFLICT (user_id) DO UPDATE SET business_name = EXCLUDED.business_name
			RETURNING id
		`, userID, v.biz, v.desc, v.lng, v.lat, v.address, v.city, v.state, v.pincode,
			3.5+float64(len(v.biz)%15)/10.0, 20+len(v.biz)%50).Scan(&vendorID)
		if err != nil {
			log.Printf("  ⚠ Vendor %s: %v", v.biz, err)
			db.QueryRow("SELECT id FROM vendors WHERE user_id = $1", userID).Scan(&vendorID)
		}
		vendorIDs[v.biz] = vendorID
		fmt.Printf("  ✓ %s (%s)\n", v.biz, v.city)
	}

	// ─────────────────────────────────────────────────────────────
	// 2. Create TN-specific Categories
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n📂 Creating Tamil Nadu categories...")

	type category struct {
		name, slug, catType string
		parentSlug          string // empty = root
	}

	tnCategories := []category{
		// New root categories
		{"Millets & Traditional Grains", "millets-and-traditional-grains", "product", ""},
		{"Pickles & Chutneys", "pickles-and-chutneys", "product", ""},
		{"Podi & Powder Mix", "podi-and-powder-mix", "product", ""},
		{"Traditional Snacks", "traditional-snacks", "product", ""},
		{"Traditional Sweets", "traditional-sweets", "product", ""},
		{"Filter Coffee & Tea", "filter-coffee-and-tea", "product", ""},
		{"Papad & Appalam", "papad-and-appalam", "product", ""},
		{"Pooja Items", "pooja-items", "product", ""},
		{"Jaggery & Sweeteners", "jaggery-and-sweeteners", "product", ""},
		{"Cold-Pressed Oils", "cold-pressed-oils", "product", ""},
		{"Idli Dosa & Breakfast Mix", "idli-dosa-and-breakfast-mix", "product", ""},
		{"Traditional Rice Varieties", "traditional-rice-varieties", "product", ""},
		// Sub-categories under Millets
		{"Ragi Products", "ragi-products", "product", "millets-and-traditional-grains"},
		{"Kambu & Thinai", "kambu-and-thinai", "product", "millets-and-traditional-grains"},
	}

	catIDs := make(map[string]string) // slug -> id

	// First, load existing categories
	rows, _ := db.Query("SELECT id, slug FROM categories")
	for rows.Next() {
		var id, slug string
		rows.Scan(&id, &slug)
		catIDs[slug] = id
	}
	rows.Close()

	for _, c := range tnCategories {
		var parentID *string
		if c.parentSlug != "" {
			if pid, ok := catIDs[c.parentSlug]; ok {
				parentID = &pid
			}
		}

		var catID string
		err := db.QueryRow(`
			INSERT INTO categories (id, name, slug, category_type, parent_id, is_active, sort_order)
			VALUES (uuid_generate_v4(), $1, $2, $3, $4, true, 0)
			ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, c.name, c.slug, c.catType, parentID).Scan(&catID)
		if err != nil {
			log.Printf("  ⚠ Category %s: %v", c.name, err)
			db.QueryRow("SELECT id FROM categories WHERE slug = $1", c.slug).Scan(&catID)
		}
		catIDs[c.slug] = catID
		fmt.Printf("  ✓ %s\n", c.name)
	}

	// ─────────────────────────────────────────────────────────────
	// 3. Insert TN Products
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n🛒 Creating Tamil Nadu products...")

	type product struct {
		name, desc, unit, catSlug, vendorName string
		price                                 float64
		comparePrice                          float64
		stock, weight                         int
		tags                                  []string
		imageURL                              string
	}

	products := []product{
		// ── Traditional Rice Varieties ───────────────────────────
		{"Ponni Raw Rice (Premium)", "Aged Ponni raw rice from Thanjavur delta, perfect for everyday cooking", "5kg", "traditional-rice-varieties", "Nallennai Groceries", 280, 320, 200, 5000, []string{"ponni", "rice", "thanjavur", "raw"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Ponni Boiled Rice", "Double-boiled Ponni rice, fluffy and aromatic", "5kg", "traditional-rice-varieties", "Nallennai Groceries", 300, 350, 180, 5000, []string{"ponni", "boiled", "rice"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Seeraga Samba Rice", "Premium biryani rice from Tamil Nadu, thin grain with natural aroma", "1kg", "traditional-rice-varieties", "Chettinad Masala House", 180, 220, 100, 1000, []string{"seeraga", "samba", "biryani", "rice"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Mappillai Samba Rice", "Heritage red rice variety, high in iron and nutrients", "1kg", "traditional-rice-varieties", "Tamil Millets Store", 160, 200, 80, 1000, []string{"mappillai", "samba", "heritage", "red-rice"}, "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Thooyamalli Rice", "Traditional aromatic rice from Cauvery delta", "2kg", "traditional-rice-varieties", "Nallennai Groceries", 220, 260, 60, 2000, []string{"thooyamalli", "aromatic", "traditional"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Karuppu Kavuni Rice (Black Rice)", "Antioxidant-rich black rice, used in traditional payasam", "500g", "traditional-rice-varieties", "Tamil Millets Store", 150, 190, 60, 500, []string{"black-rice", "kavuni", "antioxidant"}, "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Idli Rice", "Short grain rice specifically for soft idlis", "5kg", "traditional-rice-varieties", "Nallennai Groceries", 260, 300, 150, 5000, []string{"idli", "rice", "short-grain"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Dosa Rice", "Parboiled rice blend for crispy dosas", "5kg", "traditional-rice-varieties", "Nallennai Groceries", 270, 310, 150, 5000, []string{"dosa", "rice", "parboiled"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Millets & Traditional Grains ─────────────────────────
		{"Ragi (Finger Millet / Kezhvaragu)", "Organic ragi from Dindigul, calcium-rich superfood", "1kg", "millets-and-traditional-grains", "Tamil Millets Store", 85, 110, 200, 1000, []string{"ragi", "finger-millet", "kezhvaragu", "organic"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kambu (Pearl Millet / Bajra)", "Heat-resistant pearl millet, perfect for kanji and dosa", "1kg", "millets-and-traditional-grains", "Tamil Millets Store", 75, 95, 180, 1000, []string{"kambu", "pearl-millet", "bajra"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Thinai (Foxtail Millet)", "Light and easily digestible, great for pongal and upma", "1kg", "millets-and-traditional-grains", "Tamil Millets Store", 90, 120, 150, 1000, []string{"thinai", "foxtail-millet"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Varagu (Kodo Millet)", "Protein-rich kodo millet for rice replacement", "1kg", "millets-and-traditional-grains", "Tamil Millets Store", 95, 125, 120, 1000, []string{"varagu", "kodo-millet"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Samai (Little Millet)", "Fiber-rich little millet, ideal for biryani and pongal", "1kg", "millets-and-traditional-grains", "Tamil Millets Store", 100, 130, 120, 1000, []string{"samai", "little-millet"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kudiraivali (Barnyard Millet)", "Low glycemic index millet for diabetics", "1kg", "millets-and-traditional-grains", "Tamil Millets Store", 110, 140, 100, 1000, []string{"kudiraivali", "barnyard-millet"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Panivaragu (Proso Millet)", "Light and nutritious, great for upma and payasam", "500g", "millets-and-traditional-grains", "Tamil Millets Store", 80, 100, 80, 500, []string{"panivaragu", "proso-millet"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kollu (Horse Gram)", "Protein powerhouse, traditional rasam and soup ingredient", "500g", "millets-and-traditional-grains", "Tamil Millets Store", 65, 85, 150, 500, []string{"kollu", "horse-gram", "ulavalu"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Ragi Flour (Kezhvaragu Maavu)", "Stone-ground ragi flour for puttu, roti and porridge", "1kg", "ragi-products", "Tamil Millets Store", 95, 120, 200, 1000, []string{"ragi", "flour", "kezhvaragu-maavu"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kambu Flour (Bajra Atta)", "Stone-ground pearl millet flour for rotis", "1kg", "kambu-and-thinai", "Tamil Millets Store", 80, 100, 150, 1000, []string{"kambu", "flour", "bajra-atta"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Multi Millet Mix", "Blend of 5 millets for making dosa, idli and pongal", "1kg", "millets-and-traditional-grains", "Tamil Millets Store", 130, 160, 100, 1000, []string{"multi-millet", "mix", "healthy"}, "https://images.unsplash.com/photo-1604329760661-e71dc83f8f26?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Filter Coffee & Tea ──────────────────────────────────
		{"Kumbakonam Degree Coffee Powder", "80:20 coffee-chicory blend, authentic South Indian filter coffee", "500g", "filter-coffee-and-tea", "Kumbakonam Degree Coffee", 280, 350, 200, 500, []string{"kumbakonam", "degree-coffee", "filter-coffee"}, "https://images.unsplash.com/photo-1556742526-795a8eac090e?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Pure Coffee Powder (No Chicory)", "100% Arabica coffee, estate fresh from Yercaud", "250g", "filter-coffee-and-tea", "Kumbakonam Degree Coffee", 220, 280, 150, 250, []string{"pure-coffee", "arabica", "yercaud"}, "https://images.unsplash.com/photo-1556742526-795a8eac090e?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Peaberry Coffee Beans", "Premium single-origin peaberry from Nilgiris", "250g", "filter-coffee-and-tea", "Kumbakonam Degree Coffee", 350, 450, 80, 250, []string{"peaberry", "nilgiris", "single-origin"}, "https://images.unsplash.com/photo-1556742526-795a8eac090e?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Chicory Powder", "Pure roasted chicory for blending with coffee", "200g", "filter-coffee-and-tea", "Kumbakonam Degree Coffee", 60, 80, 200, 200, []string{"chicory", "coffee-blend"}, "https://images.unsplash.com/photo-1556742526-795a8eac090e?w=400&h=400&fit=crop&auto=format&q=80"},
		{"South Indian Filter Coffee Set (Dabara Set)", "Stainless steel tumbler and dabara for authentic filter coffee", "1 set", "filter-coffee-and-tea", "Kumbakonam Degree Coffee", 180, 250, 100, 300, []string{"filter", "dabara", "tumbler", "steel"}, "https://images.unsplash.com/photo-1556742526-795a8eac090e?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Sukku Malli Coffee", "Traditional dry ginger coriander herbal coffee", "200g", "filter-coffee-and-tea", "Kumbakonam Degree Coffee", 120, 150, 120, 200, []string{"sukku", "malli", "herbal-coffee"}, "https://images.unsplash.com/photo-1556742526-795a8eac090e?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Nannari Sherbet Syrup", "Natural sarsaparilla root syrup, refreshing summer drink", "500ml", "filter-coffee-and-tea", "Tirunelveli Treats", 150, 190, 100, 500, []string{"nannari", "sherbet", "sarsaparilla"}, "https://plus.unsplash.com/premium_photo-1666262370578-59f5736c0202?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Rose Milk Syrup", "Traditional rose-flavored milk syrup", "500ml", "filter-coffee-and-tea", "Madurai Snacks & Sweets", 90, 120, 150, 500, []string{"rose-milk", "syrup"}, "https://plus.unsplash.com/premium_photo-1666262370578-59f5736c0202?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Panagam Mix (Jaggery Drink)", "Traditional temple drink mix with dry ginger and cardamom", "200g", "filter-coffee-and-tea", "Madurai Snacks & Sweets", 80, 100, 100, 200, []string{"panagam", "jaggery-drink", "temple"}, "https://plus.unsplash.com/premium_photo-1666262370578-59f5736c0202?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Cold-Pressed Oils (Marachekku / Chekku) ──────────────
		{"Nallennai (Gingelly Oil / Sesame Oil)", "Traditional wood-pressed sesame oil from Erode", "1L", "cold-pressed-oils", "Nallennai Groceries", 380, 450, 150, 1000, []string{"nallennai", "sesame-oil", "gingelly", "chekku"}, "https://plus.unsplash.com/premium_photo-1661454028102-e949326751a3?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Chekku Groundnut Oil", "Cold-pressed groundnut oil, perfect for Tamil cooking", "1L", "cold-pressed-oils", "Nallennai Groceries", 320, 380, 150, 1000, []string{"groundnut-oil", "chekku", "cold-pressed"}, "https://plus.unsplash.com/premium_photo-1667818824583-3be6f268bb13?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Chekku Coconut Oil", "Virgin cold-pressed coconut oil from Pollachi", "1L", "cold-pressed-oils", "Kovai Organics", 350, 420, 120, 1000, []string{"coconut-oil", "chekku", "pollachi", "virgin"}, "https://plus.unsplash.com/premium_photo-1661454028102-e949326751a3?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Castor Oil (Amanakku Ennai)", "Cold-pressed castor oil for cooking and medicinal use", "500ml", "cold-pressed-oils", "Nallennai Groceries", 200, 250, 80, 500, []string{"castor-oil", "amanakku", "cold-pressed"}, "https://plus.unsplash.com/premium_photo-1661454028102-e949326751a3?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Iluppa Ennai (Mahua Oil)", "Traditional oil for temple lamps and cooking", "500ml", "cold-pressed-oils", "Nallennai Groceries", 180, 220, 60, 500, []string{"iluppa", "mahua-oil", "traditional"}, "https://plus.unsplash.com/premium_photo-1661454028102-e949326751a3?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Podi & Powder Mix ────────────────────────────────────
		{"Idli Milagai Podi (Gun Powder)", "Spicy lentil powder for idli, dosa with gingelly oil", "200g", "podi-and-powder-mix", "Chettinad Masala House", 90, 120, 200, 200, []string{"milagai-podi", "gun-powder", "idli"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Ellu Podi (Sesame Powder)", "Nutritious sesame and cumin powder for rice", "200g", "podi-and-powder-mix", "Chettinad Masala House", 80, 100, 180, 200, []string{"ellu-podi", "sesame", "til"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Paruppu Podi (Dal Powder)", "Roasted lentil powder with red chillies", "200g", "podi-and-powder-mix", "Chettinad Masala House", 85, 110, 180, 200, []string{"paruppu-podi", "dal-powder"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Curry Leaves Podi (Karuveppilai Podi)", "Aromatic curry leaf powder with urad dal", "150g", "podi-and-powder-mix", "Chettinad Masala House", 70, 90, 150, 150, []string{"karuveppilai", "curry-leaves", "podi"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Sundakkai Podi (Turkey Berry Powder)", "Traditional digestive powder from dried turkey berries", "100g", "podi-and-powder-mix", "Chettinad Masala House", 75, 95, 100, 100, []string{"sundakkai", "turkey-berry", "digestive"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Coconut Chutney Powder", "Instant coconut chutney powder, just add water", "200g", "podi-and-powder-mix", "Chettinad Masala House", 95, 120, 150, 200, []string{"coconut", "chutney", "instant"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Sambar Powder (Homemade Style)", "Stone-ground sambar masala with Chettinad recipe", "250g", "podi-and-powder-mix", "Chettinad Masala House", 110, 140, 200, 250, []string{"sambar", "powder", "chettinad"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Rasam Powder (Homemade Style)", "Aromatic rasam powder with pepper and cumin", "200g", "podi-and-powder-mix", "Chettinad Masala House", 95, 120, 200, 200, []string{"rasam", "powder", "pepper"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Chettinad Chicken Masala", "Authentic Chettinad spice blend for chicken curry", "100g", "podi-and-powder-mix", "Chettinad Masala House", 85, 110, 150, 100, []string{"chettinad", "chicken", "masala"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Chettinad Fish Fry Masala", "Spicy fish fry masala with fennel and star anise", "100g", "podi-and-powder-mix", "Chettinad Masala House", 80, 100, 150, 100, []string{"chettinad", "fish-fry", "masala"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kulambu Masala Powder", "All-purpose gravy masala for kuzhambu and kootu", "200g", "podi-and-powder-mix", "Chettinad Masala House", 90, 115, 180, 200, []string{"kulambu", "kuzhambu", "masala"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Puttu Powder", "Steam-roasted rice flour for making puttu", "500g", "podi-and-powder-mix", "Nallennai Groceries", 60, 80, 120, 500, []string{"puttu", "rice-flour", "steamed"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Pickles & Chutneys ───────────────────────────────────
		{"Avakkai Mango Pickle", "Andhra-TN style raw mango pickle with mustard and fenugreek", "300g", "pickles-and-chutneys", "Chettinad Masala House", 130, 160, 150, 300, []string{"avakkai", "mango", "pickle"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Narthangai Pickle (Citron)", "Tangy citron pickle, classic Tamil accompaniment", "300g", "pickles-and-chutneys", "Chettinad Masala House", 120, 150, 120, 300, []string{"narthangai", "citron", "pickle"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Maavadu (Baby Mango Pickle)", "Tender baby mango pickle in brine", "300g", "pickles-and-chutneys", "Madurai Snacks & Sweets", 140, 170, 100, 300, []string{"maavadu", "baby-mango", "pickle"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Lemon Pickle (Elumichai Oorugai)", "Traditional lemon pickle with gingelly oil", "250g", "pickles-and-chutneys", "Chettinad Masala House", 100, 130, 150, 250, []string{"lemon", "elumichai", "pickle"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Garlic Pickle (Poondu Oorugai)", "Spicy garlic pickle with red chillies", "250g", "pickles-and-chutneys", "Chettinad Masala House", 110, 140, 120, 250, []string{"garlic", "poondu", "pickle"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Gongura Pickle", "Tangy sorrel leaves pickle, Andhra-TN favorite", "300g", "pickles-and-chutneys", "Chettinad Masala House", 140, 170, 100, 300, []string{"gongura", "sorrel", "pickle"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kadarangai (Lime Rind Pickle)", "Dried lime rind pickle, unique TN specialty", "200g", "pickles-and-chutneys", "Madurai Snacks & Sweets", 110, 140, 80, 200, []string{"kadarangai", "lime-rind", "pickle"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Inji Marappa (Ginger Candy Pickle)", "Sweet and spicy ginger preserve", "200g", "pickles-and-chutneys", "Madurai Snacks & Sweets", 95, 120, 100, 200, []string{"inji", "ginger", "candy", "marappa"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Manga Thokku (Mango Paste)", "Cooked mango chutney paste for rice", "300g", "pickles-and-chutneys", "Madurai Snacks & Sweets", 120, 150, 120, 300, []string{"manga-thokku", "mango", "paste"}, "https://images.unsplash.com/photo-1633321702518-7fecdafb94d5?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Traditional Snacks ───────────────────────────────────
		{"Murukku (Chakli)", "Crunchy spiral rice flour snack with cumin", "250g", "traditional-snacks", "Madurai Snacks & Sweets", 90, 120, 200, 250, []string{"murukku", "chakli", "rice"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Thattai", "Crispy flattened rice and dal crackers", "200g", "traditional-snacks", "Madurai Snacks & Sweets", 80, 100, 200, 200, []string{"thattai", "crispy", "traditional"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Seedai (Sweet & Salt)", "Traditional deep-fried rice flour balls", "200g", "traditional-snacks", "Madurai Snacks & Sweets", 85, 110, 150, 200, []string{"seedai", "uppu-seedai", "vella-seedai"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Ribbon Pakoda", "Thin crispy ribbon-shaped gram flour snack", "200g", "traditional-snacks", "Madurai Snacks & Sweets", 75, 95, 200, 200, []string{"ribbon", "pakoda", "gram-flour"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Omapodi", "Spicy thin sev with ajwain flavor", "200g", "traditional-snacks", "Madurai Snacks & Sweets", 70, 90, 200, 200, []string{"omapodi", "sev", "ajwain"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Karasev", "Thick spicy besan sev, Tamil Nadu favorite", "200g", "traditional-snacks", "Madurai Snacks & Sweets", 75, 95, 200, 200, []string{"karasev", "besan", "spicy"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Banana Chips (Nendran)", "Crispy Kerala-style banana chips with coconut oil", "250g", "traditional-snacks", "Tirunelveli Treats", 110, 140, 180, 250, []string{"banana-chips", "nendran", "coconut-oil"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Mixture (Madras Mix)", "Classic South Indian snack mixture", "250g", "traditional-snacks", "Madurai Snacks & Sweets", 85, 110, 200, 250, []string{"mixture", "madras-mix", "south-indian"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Athirasam", "Traditional jaggery and rice flour sweet snack", "6 pcs", "traditional-snacks", "Madurai Snacks & Sweets", 120, 150, 100, 300, []string{"athirasam", "adhirasam", "jaggery"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kuzhi Paniyaram Mix", "Ready mix for sweet and savory paniyaram", "500g", "traditional-snacks", "Kovai Organics", 95, 120, 120, 500, []string{"paniyaram", "kuzhi", "mix"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Traditional Sweets ───────────────────────────────────
		{"Tirunelveli Halwa (Wheat Halwa)", "Famous Tirunelveli iruttu kadai style wheat halwa", "250g", "traditional-sweets", "Tirunelveli Treats", 200, 250, 100, 250, []string{"tirunelveli", "halwa", "wheat", "iruttu-kadai"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Mysore Pak", "Ghee-rich gram flour sweet, soft and melt-in-mouth", "250g", "traditional-sweets", "Madurai Snacks & Sweets", 180, 220, 120, 250, []string{"mysore-pak", "besan", "ghee"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Jangiri", "Saffron-soaked urad dal sweet, festival favorite", "6 pcs", "traditional-sweets", "Madurai Snacks & Sweets", 120, 150, 100, 200, []string{"jangiri", "imarti", "saffron"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Badusha", "Flaky ghee-soaked pastry with sugar syrup", "6 pcs", "traditional-sweets", "Madurai Snacks & Sweets", 130, 160, 100, 200, []string{"badusha", "balushahi", "ghee"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Palkova (Milk Sweet)", "Reduced milk fudge, Srivilliputhur specialty", "200g", "traditional-sweets", "Tirunelveli Treats", 160, 200, 80, 200, []string{"palkova", "milk-sweet", "srivilliputhur"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Laddu (Boondi Laddu)", "Sweet boondi laddu with cardamom and cashews", "6 pcs", "traditional-sweets", "Madurai Snacks & Sweets", 140, 170, 120, 300, []string{"laddu", "boondi", "cardamom"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Nei Appam (Ghee Appam)", "Golden ghee-fried jaggery appam", "6 pcs", "traditional-sweets", "Madurai Snacks & Sweets", 100, 130, 100, 200, []string{"nei-appam", "ghee", "jaggery"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kozhukattai Mix (Modak)", "Ready mix for sweet kozhukattai with jaggery filling", "300g", "traditional-sweets", "Kovai Organics", 85, 110, 100, 300, []string{"kozhukattai", "modak", "mix"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Karupatti Mittai (Palm Jaggery Candy)", "Traditional palm jaggery candy with peanuts", "200g", "traditional-sweets", "Tirunelveli Treats", 90, 110, 100, 200, []string{"karupatti", "palm-jaggery", "candy"}, "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Idli Dosa & Breakfast Mix ────────────────────────────
		{"Idli Rava (Idli Rice Broken)", "Coarse broken rice for soft fluffy idlis", "1kg", "idli-dosa-and-breakfast-mix", "Nallennai Groceries", 70, 90, 200, 1000, []string{"idli-rava", "broken-rice"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Urad Dal (Ulundu) - Whole White", "Premium quality whole white urad dal for idli batter", "1kg", "idli-dosa-and-breakfast-mix", "Nallennai Groceries", 160, 200, 200, 1000, []string{"urad-dal", "ulundu", "whole-white"}, "https://plus.unsplash.com/premium_photo-1701064865162-db655bfb99c3?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Aachi Idli Mix", "Instant idli mix, just add water", "500g", "idli-dosa-and-breakfast-mix", "Kovai Organics", 75, 95, 200, 500, []string{"aachi", "idli", "instant-mix"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Aachi Dosa Mix", "Instant crispy dosa mix", "500g", "idli-dosa-and-breakfast-mix", "Kovai Organics", 80, 100, 200, 500, []string{"aachi", "dosa", "instant-mix"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Rava Dosa Mix", "Instant crispy rava dosa with cashews and pepper", "500g", "idli-dosa-and-breakfast-mix", "Kovai Organics", 85, 110, 150, 500, []string{"rava-dosa", "instant-mix"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Pongal Mix (Ven Pongal)", "Ready mix for creamy ven pongal", "300g", "idli-dosa-and-breakfast-mix", "Kovai Organics", 70, 90, 150, 300, []string{"pongal", "ven-pongal", "mix"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Upma Rava (Bombay Rava)", "Medium-grain semolina for perfect upma", "1kg", "idli-dosa-and-breakfast-mix", "Nallennai Groceries", 55, 70, 200, 1000, []string{"upma", "rava", "semolina"}, "https://plus.unsplash.com/premium_photo-1725878608875-ca527315bb46?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Adai Mix (Multi-Dal Dosa)", "Protein-rich multi-lentil dosa mix", "500g", "idli-dosa-and-breakfast-mix", "Tamil Millets Store", 95, 120, 120, 500, []string{"adai", "multi-dal", "dosa"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Pesarattu Mix (Green Gram Dosa)", "Andhra-style green gram dosa mix", "500g", "idli-dosa-and-breakfast-mix", "Tamil Millets Store", 90, 115, 100, 500, []string{"pesarattu", "green-gram", "dosa"}, "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Jaggery & Sweeteners ─────────────────────────────────
		{"Karupatti (Palm Jaggery)", "Natural palm jaggery from Tirunelveli, unprocessed", "500g", "jaggery-and-sweeteners", "Tirunelveli Treats", 120, 150, 200, 500, []string{"karupatti", "palm-jaggery", "tirunelveli"}, "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Vellam (Cane Jaggery Block)", "Traditional cane jaggery block for pongal and payasam", "1kg", "jaggery-and-sweeteners", "Nallennai Groceries", 80, 100, 200, 1000, []string{"vellam", "cane-jaggery", "block"}, "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Palm Jaggery Powder", "Powdered palm jaggery for easy cooking", "500g", "jaggery-and-sweeteners", "Tirunelveli Treats", 140, 170, 150, 500, []string{"palm-jaggery", "powder", "natural"}, "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Nattu Sakkarai (Country Sugar)", "Unrefined country sugar, less processed", "500g", "jaggery-and-sweeteners", "Kovai Organics", 70, 90, 150, 500, []string{"nattu-sakkarai", "country-sugar", "unrefined"}, "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Panangkarkandu (Palm Candy)", "Crystallized palm sugar, traditional sweetener", "250g", "jaggery-and-sweeteners", "Tirunelveli Treats", 110, 140, 100, 250, []string{"panangkarkandu", "palm-candy", "crystal"}, "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Papad & Appalam ──────────────────────────────────────
		{"Appalam (Madras Papad)", "Thin crispy urad dal papad, large size", "200g", "papad-and-appalam", "Madurai Snacks & Sweets", 60, 80, 200, 200, []string{"appalam", "papad", "urad"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Vadam (Sun-Dried Rice Crackers)", "Traditional sun-dried rice and sago crackers", "200g", "papad-and-appalam", "Madurai Snacks & Sweets", 70, 90, 150, 200, []string{"vadam", "rice-crackers", "sun-dried"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Vathal (Dried Vegetables)", "Assorted dried vegetables for frying - sundakkai, manathakkali", "100g", "papad-and-appalam", "Chettinad Masala House", 90, 120, 100, 100, []string{"vathal", "dried-vegetables", "fry"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Fried Gram (Pottukadalai)", "Roasted split chickpeas, essential chutney ingredient", "500g", "papad-and-appalam", "Nallennai Groceries", 75, 95, 200, 500, []string{"pottukadalai", "fried-gram", "roasted"}, "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Pooja Items ──────────────────────────────────────────
		{"Vibhuti (Sacred Ash)", "Pure vibhuti from temple source", "100g", "pooja-items", "Madurai Snacks & Sweets", 30, 40, 200, 100, []string{"vibhuti", "sacred-ash", "pooja"}, "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kumkum Powder", "Bright red kumkum for daily pooja", "50g", "pooja-items", "Madurai Snacks & Sweets", 25, 35, 200, 50, []string{"kumkum", "sindoor", "pooja"}, "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Camphor (Karpooram)", "Pure edible camphor for aarti and cooking", "50g", "pooja-items", "Madurai Snacks & Sweets", 45, 60, 200, 50, []string{"camphor", "karpooram", "aarti"}, "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Incense Sticks (Agarbatti) - Sandalwood", "Premium sandalwood fragrance incense sticks", "50 sticks", "pooja-items", "Madurai Snacks & Sweets", 40, 55, 200, 100, []string{"agarbatti", "incense", "sandalwood"}, "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Sambrani Cups (Benzoin)", "Dhoop sambrani cups for traditional fumigation", "12 cups", "pooja-items", "Madurai Snacks & Sweets", 35, 50, 200, 100, []string{"sambrani", "benzoin", "dhoop"}, "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Cotton Wicks (Thiri)", "Hand-rolled cotton wicks for oil lamps", "100 pcs", "pooja-items", "Madurai Snacks & Sweets", 20, 30, 300, 50, []string{"cotton-wicks", "thiri", "oil-lamp"}, "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Turmeric Sticks (Manjal Kizhangu)", "Whole turmeric rhizomes for pooja and cooking", "200g", "pooja-items", "Kovai Organics", 50, 65, 150, 200, []string{"manjal", "turmeric", "whole", "pooja"}, "https://images.unsplash.com/photo-1604608672516-f1b9b1d37076?w=400&h=400&fit=crop&auto=format&q=80"},

		// ── Tamarind & Essentials ────────────────────────────────
		{"Puli (Tamarind Block)", "Seedless pressed tamarind from Tamil Nadu", "500g", "masala-and-spices", "Chettinad Masala House", 65, 85, 200, 500, []string{"puli", "tamarind", "seedless"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Tamarind Paste (Puli Paste)", "Ready-to-use tamarind concentrate", "300g", "masala-and-spices", "Chettinad Masala House", 55, 70, 180, 300, []string{"tamarind", "paste", "concentrate"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Curry Leaves (Karuveppilai) - Fresh", "Aromatic fresh curry leaves from organic farm", "100g", "masala-and-spices", "Kovai Organics", 15, 25, 300, 100, []string{"curry-leaves", "karuveppilai", "fresh"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Vendhayam (Fenugreek Seeds)", "Whole fenugreek seeds for tempering and sprouting", "200g", "masala-and-spices", "Chettinad Masala House", 40, 55, 200, 200, []string{"vendhayam", "fenugreek", "methi"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Sombu (Fennel Seeds)", "Aromatic fennel seeds for cooking and mouth freshener", "200g", "masala-and-spices", "Chettinad Masala House", 55, 70, 200, 200, []string{"sombu", "fennel", "saunf"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Star Anise (Thakkolam)", "Whole star anise for biryani and Chettinad cooking", "50g", "masala-and-spices", "Chettinad Masala House", 60, 80, 150, 50, []string{"star-anise", "thakkolam", "biryani"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Kalpasi (Stone Flower)", "Rare Chettinad spice for authentic chicken curry", "25g", "masala-and-spices", "Chettinad Masala House", 80, 100, 80, 25, []string{"kalpasi", "stone-flower", "chettinad"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
		{"Marathi Mokku (Dried Flower)", "Dried flower buds for Chettinad masala", "25g", "masala-and-spices", "Chettinad Masala House", 70, 90, 80, 25, []string{"marathi-mokku", "dried-flower", "chettinad"}, "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817?w=400&h=400&fit=crop&auto=format&q=80"},
	}

	inserted := 0
	for _, p := range products {
		vendorID, ok := vendorIDs[p.vendorName]
		if !ok {
			log.Printf("  ⚠ Vendor not found: %s", p.vendorName)
			continue
		}
		catID, ok := catIDs[p.catSlug]
		if !ok {
			log.Printf("  ⚠ Category not found: %s", p.catSlug)
			continue
		}

		slug := strings.ToLower(strings.ReplaceAll(p.name, " ", "-"))
		slug = strings.ReplaceAll(slug, "(", "")
		slug = strings.ReplaceAll(slug, ")", "")
		slug = strings.ReplaceAll(slug, "/", "-")
		slug = strings.ReplaceAll(slug, "&", "and")
		slug = strings.ReplaceAll(slug, ",", "")
		slug = strings.ReplaceAll(slug, "'", "")
		slug = strings.ReplaceAll(slug, ".", "")

		tags := "{" + strings.Join(p.tags, ",") + "}"

		var comparePrice *float64
		if p.comparePrice > 0 {
			comparePrice = &p.comparePrice
		}

		// Check if product with this slug already exists
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM products WHERE slug = $1 AND deleted_at IS NULL", slug).Scan(&exists)
		if exists > 0 {
			log.Printf("  ⏭ %s (slug already exists)", p.name)
			inserted++ // count as success since it's already there
			continue
		}

		_, err := db.Exec(`
			INSERT INTO products
				(id, vendor_id, category_id, name, slug, description, price, compare_price,
				 unit, stock_quantity, low_stock_threshold, weight_grams, tags, images, is_active)
			VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8, $9, 5, $10, $11, ARRAY[$12], true)
		`, vendorID, catID, p.name, slug, p.desc, p.price, comparePrice,
			p.unit, p.stock, p.weight, tags, p.imageURL)
		if err != nil {
			log.Printf("  ✗ %s: %v", p.name, err)
			continue
		}
		inserted++
	}

	// ── Final stats ──────────────────────────────────────────────
	var totalProducts, totalCategories, totalVendors int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalProducts)
	db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&totalCategories)
	db.QueryRow("SELECT COUNT(*) FROM vendors").Scan(&totalVendors)

	fmt.Println("\n══════════════════════════════════════════════════════")
	fmt.Printf("  ✅  %d Tamil Nadu products inserted\n", inserted)
	fmt.Printf("  📊  Total: %d products | %d categories | %d vendors\n", totalProducts, totalCategories, totalVendors)
	fmt.Println("══════════════════════════════════════════════════════")
}
