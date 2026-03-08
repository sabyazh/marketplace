package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// ─── helpers ────────────────────────────────────────────────────────────────

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " & ", "-and-")
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

func mustExec(db *sql.DB, query string, args ...interface{}) {
	if _, err := db.Exec(query, args...); err != nil {
		log.Printf("WARN: %v  |  query: %.80s", err, query)
	}
}

// ─── main ───────────────────────────────────────────────────────────────────

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
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte("Vendor@123"), bcrypt.DefaultCost)
	pwHash := string(pw)

	// ═══════════════════════════════════════════════════════════════════════
	// 1. CATEGORIES  (Indian grocery taxonomy)
	// ═══════════════════════════════════════════════════════════════════════
	type cat struct{ name, slug, desc, parent string; sort int }
	categories := []cat{
		// Top-level
		{"Fruits & Vegetables", "fruits-and-vegetables", "Fresh fruits and vegetables", "", 1},
		{"Dairy & Bread", "dairy-and-bread", "Milk, curd, paneer, bread & eggs", "", 2},
		{"Rice & Atta", "rice-and-atta", "Rice, wheat flour, rava, poha & grains", "", 3},
		{"Dals & Pulses", "dals-and-pulses", "Toor dal, chana dal, moong dal & more", "", 4},
		{"Masala & Spices", "masala-and-spices", "Ground spices, whole spices & blends", "", 5},
		{"Cooking Oil & Ghee", "cooking-oil-and-ghee", "Sunflower oil, mustard oil, ghee & more", "", 6},
		{"Snacks & Namkeen", "snacks-and-namkeen", "Chips, biscuits, namkeen & sweets", "", 7},
		{"Beverages", "beverages", "Tea, coffee, juices & soft drinks", "", 8},
		{"Dry Fruits & Nuts", "dry-fruits-and-nuts", "Almonds, cashews, raisins & more", "", 9},
		{"Instant & Ready to Eat", "instant-and-ready-to-eat", "Noodles, ready meals, soup mixes", "", 10},
		{"Personal Care", "personal-care", "Soaps, shampoos, skincare & hygiene", "", 11},
		{"Cleaning & Household", "cleaning-and-household", "Detergents, cleaners & kitchen supplies", "", 12},
		// Sub-categories
		{"Fresh Vegetables", "fresh-vegetables", "Daily vegetables", "fruits-and-vegetables", 1},
		{"Fresh Fruits", "fresh-fruits", "Seasonal fruits", "fruits-and-vegetables", 2},
		{"Milk & Curd", "milk-and-curd", "Fresh milk, curd & buttermilk", "dairy-and-bread", 1},
		{"Paneer & Cheese", "paneer-and-cheese", "Paneer, cheese & tofu", "dairy-and-bread", 2},
		{"Bread & Eggs", "bread-and-eggs", "Bread, bun, pav & eggs", "dairy-and-bread", 3},
		{"Basmati Rice", "basmati-rice", "Premium basmati rice", "rice-and-atta", 1},
		{"Wheat Atta", "wheat-atta", "Whole wheat flour", "rice-and-atta", 2},
		{"Other Grains", "other-grains", "Rava, poha, maida & semolina", "rice-and-atta", 3},
	}

	catIDs := map[string]string{}
	for _, c := range categories {
		var parentID interface{} = nil
		if c.parent != "" {
			if pid, ok := catIDs[c.parent]; ok {
				parentID = pid
			}
		}
		var id string
		err := db.QueryRow(`
			INSERT INTO categories (name, slug, description, parent_id, sort_order, is_active, category_type)
			VALUES ($1, $2, $3, $4, $5, true, 'product')
			ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, c.name, c.slug, c.desc, parentID, c.sort).Scan(&id)
		if err != nil {
			log.Printf("category %s: %v", c.name, err)
			// try to get existing
			db.QueryRow("SELECT id FROM categories WHERE slug=$1", c.slug).Scan(&id)
		}
		catIDs[c.slug] = id
		fmt.Printf("  ✓ Category: %-30s  %s\n", c.name, id)
	}

	// ═══════════════════════════════════════════════════════════════════════
	// 2. VENDOR USERS & VENDORS
	// ═══════════════════════════════════════════════════════════════════════
	type vendorDef struct {
		name, email, phone, biz, desc, city, state, pin, addr string
		lat, lng float64
	}
	vendors := []vendorDef{
		{
			"Rajesh Kumar", "rajesh@freshbasket.in", "9876543210",
			"Fresh Basket", "Your neighbourhood fresh grocery store with farm-fresh vegetables and fruits daily",
			"Chennai", "Tamil Nadu", "600001", "42 Anna Salai, T Nagar",
			13.0418, 80.2341,
		},
		{
			"Priya Sharma", "priya@annapurnastore.in", "9876543211",
			"Annapurna Grocery", "Premium quality groceries, spices and daily essentials at best prices",
			"Mumbai", "Maharashtra", "400001", "15 MG Road, Fort",
			18.9322, 72.8347,
		},
		{
			"Mohammed Irfan", "irfan@spiceemporium.in", "9876543212",
			"Spice Emporium", "Authentic Indian spices sourced directly from Kerala and Karnataka farms",
			"Bangalore", "Karnataka", "560001", "88 Commercial Street",
			12.9812, 77.6094,
		},
		{
			"Sunita Patel", "sunita@desimart.in", "9876543213",
			"Desi Mart", "One-stop shop for dals, rice, atta and all kitchen staples",
			"Hyderabad", "Telangana", "500001", "23 Begumpet Road",
			17.4432, 78.4982,
		},
		{
			"Amit Verma", "amit@organicwala.in", "9876543214",
			"Organic Wala", "100% organic and chemical-free groceries from certified farms",
			"Delhi", "Delhi", "110001", "56 Chandni Chowk",
			28.6562, 77.2310,
		},
	}

	vendorIDs := map[string]string{} // biz name → vendor id
	for _, v := range vendors {
		// Create user
		var userID string
		err := db.QueryRow(`
			INSERT INTO users (full_name, email, phone, password_hash, role, is_active)
			VALUES ($1, $2, $3, $4, 'vendor', true)
			ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name
			RETURNING id
		`, v.name, v.email, v.phone, pwHash).Scan(&userID)
		if err != nil {
			db.QueryRow("SELECT id FROM users WHERE email=$1", v.email).Scan(&userID)
		}

		// Create vendor
		var vendorID string
		err = db.QueryRow(`
			INSERT INTO vendors (user_id, business_name, description, vendor_type, status,
				location, address, city, state, pincode, is_online, avg_rating, total_reviews, commission_pct)
			VALUES ($1, $2, $3, 'product', 'approved',
				ST_MakePoint($4, $5)::geography, $6, $7, $8, $9, true, $10, $11, 10.00)
			ON CONFLICT (user_id) DO UPDATE SET business_name = EXCLUDED.business_name
			RETURNING id
		`, userID, v.biz, v.desc, v.lng, v.lat, v.addr, v.city, v.state, v.pin,
			3.5+float64(len(v.biz)%15)/10.0, 20+len(v.biz)%50).Scan(&vendorID)
		if err != nil {
			db.QueryRow("SELECT id FROM vendors WHERE user_id=$1", userID).Scan(&vendorID)
		}
		vendorIDs[v.biz] = vendorID
		fmt.Printf("  ✓ Vendor:   %-30s  %s  (%s)\n", v.biz, vendorID, v.city)
	}

	// ═══════════════════════════════════════════════════════════════════════
	// 3. PRODUCTS  (~120 Indian grocery products)
	// ═══════════════════════════════════════════════════════════════════════
	type prod struct {
		name     string
		price    float64
		compare  float64
		unit     string
		cat      string // category slug
		vendor   string // vendor biz name
		stock    int
		weight   int // grams
		desc     string
		tags     string // comma-separated
	}

	products := []prod{
		// ── Fresh Vegetables ─────────────────────────────────────
		{"Tomato (Tamatar)", 40, 50, "kg", "fresh-vegetables", "Fresh Basket", 200, 1000, "Farm-fresh red tomatoes, perfect for curries and chutneys", "tomato,tamatar,sabzi"},
		{"Onion (Pyaz)", 35, 45, "kg", "fresh-vegetables", "Fresh Basket", 300, 1000, "Medium-sized onions, essential for every Indian kitchen", "onion,pyaz,kanda"},
		{"Potato (Aloo)", 30, 40, "kg", "fresh-vegetables", "Fresh Basket", 250, 1000, "Fresh potatoes ideal for aloo gobi, paratha and fries", "potato,aloo,batata"},
		{"Green Chilli", 60, 80, "kg", "fresh-vegetables", "Fresh Basket", 100, 250, "Spicy green chillies for tadka and curries", "chilli,mirchi,hari-mirchi"},
		{"Ginger (Adrak)", 120, 150, "kg", "fresh-vegetables", "Fresh Basket", 80, 250, "Fresh ginger root for tea, curries and pickles", "ginger,adrak"},
		{"Garlic (Lehsun)", 160, 200, "kg", "fresh-vegetables", "Fresh Basket", 80, 250, "Aromatic garlic bulbs for everyday cooking", "garlic,lehsun"},
		{"Coriander Leaves (Dhaniya)", 15, 20, "bunch", "fresh-vegetables", "Fresh Basket", 150, 100, "Fresh coriander for garnishing and chutneys", "coriander,dhaniya,cilantro"},
		{"Lady Finger (Bhindi)", 50, 65, "kg", "fresh-vegetables", "Fresh Basket", 120, 500, "Tender bhindi for fry and curry preparations", "bhindi,ladyfinger,okra"},
		{"Brinjal (Baingan)", 40, 55, "kg", "fresh-vegetables", "Fresh Basket", 100, 500, "Purple brinjal for bhartha and South Indian dishes", "brinjal,baingan,eggplant"},
		{"Cauliflower (Gobi)", 35, 50, "piece", "fresh-vegetables", "Fresh Basket", 80, 500, "Fresh white cauliflower for gobi manchurian and paratha", "cauliflower,gobi,phool-gobi"},
		{"Spinach (Palak)", 25, 35, "bunch", "fresh-vegetables", "Fresh Basket", 100, 250, "Fresh spinach leaves for palak paneer and soup", "spinach,palak,saag"},
		{"Drumstick (Moringa)", 45, 60, "kg", "fresh-vegetables", "Fresh Basket", 60, 500, "Long drumsticks for sambar and curry", "drumstick,moringa,sahjan"},

		// ── Fresh Fruits ──────────────────────────────────────────
		{"Banana (Kela)", 45, 55, "dozen", "fresh-fruits", "Fresh Basket", 200, 1200, "Sweet ripe bananas, rich in potassium", "banana,kela"},
		{"Apple (Seb) - Shimla", 180, 220, "kg", "fresh-fruits", "Fresh Basket", 100, 1000, "Premium Shimla apples, sweet and crunchy", "apple,seb,shimla"},
		{"Mango - Alphonso", 350, 450, "kg", "fresh-fruits", "Fresh Basket", 50, 1000, "King of mangoes - premium Ratnagiri Alphonso", "mango,alphonso,aam,hapus"},
		{"Papaya", 40, 55, "kg", "fresh-fruits", "Fresh Basket", 80, 1000, "Ripe papaya, great for digestion", "papaya,papita"},
		{"Pomegranate (Anar)", 160, 200, "kg", "fresh-fruits", "Fresh Basket", 60, 1000, "Ruby red pomegranate seeds, sweet and juicy", "pomegranate,anar"},
		{"Guava (Amrood)", 80, 100, "kg", "fresh-fruits", "Fresh Basket", 80, 1000, "Fresh guavas rich in Vitamin C", "guava,amrood,peru"},

		// ── Milk & Curd ──────────────────────────────────────────
		{"Amul Toned Milk", 29, 0, "500ml", "milk-and-curd", "Annapurna Grocery", 200, 500, "Amul toned milk 500ml pouch, 3% fat", "milk,amul,toned"},
		{"Amul Full Cream Milk", 35, 0, "500ml", "milk-and-curd", "Annapurna Grocery", 200, 500, "Amul full cream milk, 6% fat content", "milk,amul,full-cream"},
		{"Mother Dairy Curd", 40, 0, "400g", "milk-and-curd", "Annapurna Grocery", 150, 400, "Thick and creamy set curd", "curd,dahi,yogurt"},
		{"Amul Buttermilk (Chaas)", 20, 25, "200ml", "milk-and-curd", "Annapurna Grocery", 100, 200, "Masala chaas with cumin and mint flavour", "buttermilk,chaas,mattha"},
		{"Nandini Curd", 30, 0, "500g", "milk-and-curd", "Annapurna Grocery", 120, 500, "Karnataka's favourite thick curd", "curd,nandini,dahi"},

		// ── Paneer & Cheese ──────────────────────────────────────
		{"Amul Malai Paneer", 90, 110, "200g", "paneer-and-cheese", "Annapurna Grocery", 100, 200, "Fresh malai paneer block, soft and creamy", "paneer,amul,cottage-cheese"},
		{"Amul Processed Cheese", 55, 65, "200g", "paneer-and-cheese", "Annapurna Grocery", 80, 200, "Amul cheese slices for sandwiches", "cheese,amul,processed"},
		{"Britannia Cream Cheese", 75, 90, "180g", "paneer-and-cheese", "Annapurna Grocery", 60, 180, "Smooth cream cheese spread", "cream-cheese,britannia"},

		// ── Bread & Eggs ─────────────────────────────────────────
		{"Britannia White Bread", 40, 0, "pack", "bread-and-eggs", "Annapurna Grocery", 100, 400, "Soft white bread, 400g loaf", "bread,britannia,white"},
		{"Whole Wheat Bread", 45, 55, "pack", "bread-and-eggs", "Annapurna Grocery", 80, 400, "100% whole wheat bread, high fibre", "bread,atta,whole-wheat"},
		{"Farm Fresh Eggs", 75, 90, "12 pcs", "bread-and-eggs", "Annapurna Grocery", 150, 720, "Free-range farm fresh eggs, pack of 12", "eggs,anda,protein"},
		{"Pav Bread", 30, 0, "8 pcs", "bread-and-eggs", "Annapurna Grocery", 80, 320, "Soft pav for pav bhaji and vada pav", "pav,bread,bun"},

		// ── Basmati Rice ─────────────────────────────────────────
		{"India Gate Basmati Rice", 450, 550, "5kg", "basmati-rice", "Desi Mart", 100, 5000, "Premium aged basmati rice, extra long grain", "rice,basmati,india-gate"},
		{"Daawat Rozana Basmati", 320, 400, "5kg", "basmati-rice", "Desi Mart", 80, 5000, "Everyday basmati rice for daily cooking", "rice,basmati,daawat"},
		{"Fortune Biryani Special Rice", 220, 280, "1kg", "basmati-rice", "Desi Mart", 120, 1000, "Long grain rice perfect for biryani", "rice,biryani,fortune"},
		{"Sona Masoori Rice", 280, 350, "5kg", "basmati-rice", "Desi Mart", 100, 5000, "South Indian favourite Sona Masoori rice", "rice,sona-masoori,south-indian"},
		{"Ponni Raw Rice", 240, 300, "5kg", "basmati-rice", "Desi Mart", 80, 5000, "Traditional Tamil Nadu ponni rice", "rice,ponni,raw-rice"},

		// ── Wheat Atta ───────────────────────────────────────────
		{"Aashirvaad Whole Wheat Atta", 280, 330, "5kg", "wheat-atta", "Desi Mart", 150, 5000, "India's most trusted atta for soft rotis", "atta,wheat,aashirvaad"},
		{"Pillsbury Chakki Atta", 260, 310, "5kg", "wheat-atta", "Desi Mart", 100, 5000, "Freshly ground chakki atta for chapatis", "atta,wheat,pillsbury"},
		{"Fortune Chakki Atta", 240, 290, "5kg", "wheat-atta", "Desi Mart", 80, 5000, "Multi-grain atta with added nutrients", "atta,fortune,multigrain"},

		// ── Other Grains ─────────────────────────────────────────
		{"Sooji (Rava/Semolina)", 45, 55, "500g", "other-grains", "Desi Mart", 100, 500, "Fine rava for upma, halwa and rava dosa", "sooji,rava,semolina,suji"},
		{"Thick Poha (Flattened Rice)", 40, 50, "500g", "other-grains", "Desi Mart", 100, 500, "Thick beaten rice flakes for poha and chivda", "poha,flattened-rice,chura"},
		{"Maida (All Purpose Flour)", 35, 45, "500g", "other-grains", "Desi Mart", 80, 500, "Refined flour for naan, cakes and bhaturas", "maida,flour,all-purpose"},
		{"Besan (Gram Flour)", 55, 70, "500g", "other-grains", "Desi Mart", 100, 500, "Pure gram flour for pakoras, ladoo and dhokla", "besan,gram-flour,chickpea"},
		{"Ragi Flour (Finger Millet)", 65, 80, "500g", "other-grains", "Desi Mart", 60, 500, "Nutritious ragi flour for dosa and porridge", "ragi,finger-millet,nachni"},

		// ── Dals & Pulses ────────────────────────────────────────
		{"Toor Dal (Arhar)", 130, 160, "1kg", "dals-and-pulses", "Desi Mart", 120, 1000, "Premium toor dal for sambar and dal tadka", "toor,arhar,dal"},
		{"Moong Dal (Yellow)", 120, 150, "1kg", "dals-and-pulses", "Desi Mart", 100, 1000, "Split yellow moong dal, easy to digest", "moong,dal,yellow"},
		{"Chana Dal", 95, 120, "1kg", "dals-and-pulses", "Desi Mart", 100, 1000, "Bengal gram dal for chana dal fry and puran poli", "chana,dal,bengal-gram"},
		{"Masoor Dal (Red Lentil)", 100, 130, "1kg", "dals-and-pulses", "Desi Mart", 100, 1000, "Quick cooking masoor dal for everyday meals", "masoor,lentil,dal"},
		{"Urad Dal (Black Gram)", 140, 170, "1kg", "dals-and-pulses", "Desi Mart", 80, 1000, "Whole urad dal for dal makhani and medu vada", "urad,black-gram,dal"},
		{"Rajma (Kidney Beans)", 150, 180, "1kg", "dals-and-pulses", "Desi Mart", 80, 1000, "Jammu rajma for authentic rajma chawal", "rajma,kidney-beans"},
		{"Kabuli Chana (Chickpeas)", 120, 150, "1kg", "dals-and-pulses", "Desi Mart", 80, 1000, "White chickpeas for chole bhature and salads", "chana,chickpeas,kabuli,chole"},
		{"Moong Whole (Green Gram)", 110, 140, "1kg", "dals-and-pulses", "Desi Mart", 60, 1000, "Whole green moong for sprouts and curries", "moong,green-gram,whole"},

		// ── Masala & Spices ──────────────────────────────────────
		{"MDH Garam Masala", 75, 90, "100g", "masala-and-spices", "Spice Emporium", 200, 100, "Aromatic blend of 12 spices for rich curries", "garam-masala,mdh,spice-blend"},
		{"Everest Turmeric Powder (Haldi)", 45, 55, "100g", "masala-and-spices", "Spice Emporium", 200, 100, "Pure turmeric powder with high curcumin content", "haldi,turmeric,everest"},
		{"MDH Red Chilli Powder", 55, 70, "100g", "masala-and-spices", "Spice Emporium", 200, 100, "Hot red chilli powder for spicy dishes", "chilli,lal-mirch,mdh"},
		{"Everest Coriander Powder (Dhaniya)", 40, 50, "100g", "masala-and-spices", "Spice Emporium", 200, 100, "Fine coriander seed powder for curries", "dhaniya,coriander,everest"},
		{"MDH Cumin Powder (Jeera)", 65, 80, "100g", "masala-and-spices", "Spice Emporium", 150, 100, "Roasted cumin powder for raita and chaat", "jeera,cumin,mdh"},
		{"Catch Black Pepper Powder", 90, 110, "100g", "masala-and-spices", "Spice Emporium", 120, 100, "Freshly ground black pepper for seasoning", "black-pepper,kali-mirch"},
		{"Whole Cumin Seeds (Jeera)", 70, 85, "200g", "masala-and-spices", "Spice Emporium", 150, 200, "Premium cumin seeds for tadka and tempering", "jeera,cumin-seeds,whole"},
		{"Mustard Seeds (Rai)", 35, 45, "200g", "masala-and-spices", "Spice Emporium", 150, 200, "Black mustard seeds for South Indian tempering", "rai,mustard,sarson"},
		{"Whole Cloves (Laung)", 95, 120, "50g", "masala-and-spices", "Spice Emporium", 100, 50, "Aromatic cloves for biryani and pulao", "laung,cloves,whole-spice"},
		{"Cinnamon Sticks (Dalchini)", 85, 105, "50g", "masala-and-spices", "Spice Emporium", 100, 50, "Ceylon cinnamon sticks for garam masala", "dalchini,cinnamon,whole-spice"},
		{"Cardamom Green (Elaichi)", 180, 220, "50g", "masala-and-spices", "Spice Emporium", 80, 50, "Premium green cardamom for tea and desserts", "elaichi,cardamom,green"},
		{"MDH Kitchen King Masala", 65, 80, "100g", "masala-and-spices", "Spice Emporium", 150, 100, "All-purpose masala for vegetables and curries", "kitchen-king,mdh,masala"},
		{"MDH Chhole Masala", 55, 70, "100g", "masala-and-spices", "Spice Emporium", 120, 100, "Special blend for authentic chole", "chole-masala,mdh,chana"},
		{"MTR Sambar Powder", 60, 75, "200g", "masala-and-spices", "Spice Emporium", 120, 200, "Authentic South Indian sambar masala", "sambar,mtr,south-indian"},
		{"MTR Rasam Powder", 50, 65, "200g", "masala-and-spices", "Spice Emporium", 100, 200, "Tangy rasam masala blend from MTR", "rasam,mtr,south-indian"},
		{"Kashmiri Red Chilli Powder", 110, 140, "100g", "masala-and-spices", "Spice Emporium", 100, 100, "Mild red chilli for colour in tandoori and curries", "kashmiri,chilli,colour"},
		{"Asafoetida (Hing)", 95, 120, "50g", "masala-and-spices", "Spice Emporium", 100, 50, "Strong aromatic hing for dal and sambar tadka", "hing,asafoetida"},

		// ── Cooking Oil & Ghee ───────────────────────────────────
		{"Fortune Sunflower Oil", 180, 210, "1L", "cooking-oil-and-ghee", "Annapurna Grocery", 120, 1000, "Light and healthy refined sunflower oil", "sunflower,oil,fortune"},
		{"Saffola Gold Oil", 200, 240, "1L", "cooking-oil-and-ghee", "Annapurna Grocery", 100, 1000, "Blended edible oil with oryzanol for heart health", "saffola,oil,blended"},
		{"KPL Shudh Mustard Oil", 170, 200, "1L", "cooking-oil-and-ghee", "Annapurna Grocery", 80, 1000, "Pure kachi ghani mustard oil for North Indian cooking", "mustard,oil,sarson,kachi-ghani"},
		{"Amul Ghee", 560, 650, "1L", "cooking-oil-and-ghee", "Annapurna Grocery", 80, 1000, "Pure cow ghee from Amul", "ghee,amul,cow-ghee,desi-ghee"},
		{"Patanjali Cow Ghee", 520, 600, "1L", "cooking-oil-and-ghee", "Annapurna Grocery", 60, 1000, "Patanjali pure desi cow ghee", "ghee,patanjali,cow-ghee"},
		{"Fortune Groundnut Oil", 210, 250, "1L", "cooking-oil-and-ghee", "Annapurna Grocery", 60, 1000, "Cold-pressed groundnut oil for Gujarat-style cooking", "groundnut,oil,peanut"},
		{"KLF Coconut Oil", 190, 230, "1L", "cooking-oil-and-ghee", "Annapurna Grocery", 80, 1000, "Pure coconut oil for South Indian cooking and hair care", "coconut,oil,nariyal,klf"},

		// ── Snacks & Namkeen ─────────────────────────────────────
		{"Haldiram's Aloo Bhujia", 55, 70, "200g", "snacks-and-namkeen", "Annapurna Grocery", 150, 200, "Crispy potato sev namkeen from Haldiram's", "bhujia,haldirams,namkeen"},
		{"Haldiram's Moong Dal", 50, 65, "200g", "snacks-and-namkeen", "Annapurna Grocery", 120, 200, "Crunchy fried moong dal snack", "moong-dal,haldirams,namkeen"},
		{"Lay's Classic Salted Chips", 20, 0, "52g", "snacks-and-namkeen", "Annapurna Grocery", 200, 52, "Classic salted potato chips", "lays,chips,potato"},
		{"Kurkure Masala Munch", 20, 0, "80g", "snacks-and-namkeen", "Annapurna Grocery", 200, 80, "Crunchy corn puffs with Indian masala flavour", "kurkure,masala,snack"},
		{"Parle-G Glucose Biscuits", 10, 0, "80g", "snacks-and-namkeen", "Annapurna Grocery", 300, 80, "India's favourite glucose biscuit since 1939", "parle-g,biscuit,glucose"},
		{"Britannia Good Day Cashew", 30, 0, "75g", "snacks-and-namkeen", "Annapurna Grocery", 200, 75, "Buttery cashew cookies from Britannia", "good-day,britannia,cookies"},
		{"Bikaner Sev Bhujia", 45, 55, "200g", "snacks-and-namkeen", "Annapurna Grocery", 100, 200, "Traditional Rajasthani besan sev", "sev,bikaner,namkeen"},
		{"Balaji Wafers", 20, 0, "40g", "snacks-and-namkeen", "Annapurna Grocery", 200, 40, "Crunchy tomato flavour wafers from Gujarat", "balaji,wafers,chips"},

		// ── Beverages ────────────────────────────────────────────
		{"Tata Tea Gold", 210, 250, "500g", "beverages", "Annapurna Grocery", 100, 500, "Premium Assam tea with 15% long leaves", "tea,tata,gold"},
		{"Brooke Bond Red Label", 180, 220, "500g", "beverages", "Annapurna Grocery", 100, 500, "India's favourite CTC tea for strong chai", "tea,red-label,brooke-bond"},
		{"Nescafe Classic Coffee", 280, 330, "200g", "beverages", "Annapurna Grocery", 80, 200, "100% pure instant coffee granules", "coffee,nescafe,instant"},
		{"Bru Instant Coffee", 250, 300, "200g", "beverages", "Annapurna Grocery", 80, 200, "Smooth aromatic instant coffee-chicory blend", "coffee,bru,instant"},
		{"Tropicana Orange Juice", 90, 110, "1L", "beverages", "Annapurna Grocery", 80, 1000, "100% orange juice, no added sugar", "juice,tropicana,orange"},
		{"Paper Boat Aam Panna", 30, 0, "200ml", "beverages", "Annapurna Grocery", 100, 200, "Traditional raw mango drink", "paper-boat,aam-panna,mango"},
		{"Rooh Afza Sharbat", 120, 150, "750ml", "beverages", "Annapurna Grocery", 60, 750, "Iconic rose-flavoured drink syrup", "rooh-afza,sharbat,rose"},

		// ── Dry Fruits & Nuts ────────────────────────────────────
		{"California Almonds (Badam)", 280, 350, "250g", "dry-fruits-and-nuts", "Organic Wala", 100, 250, "Premium California almonds, raw and unprocessed", "almonds,badam,california"},
		{"Cashew Whole (Kaju)", 320, 400, "250g", "dry-fruits-and-nuts", "Organic Wala", 80, 250, "W320 grade whole cashew nuts from Goa", "cashew,kaju,nuts"},
		{"Kishmish (Raisins)", 120, 150, "250g", "dry-fruits-and-nuts", "Organic Wala", 100, 250, "Seedless green raisins from Afghanistan", "raisins,kishmish,dry-fruit"},
		{"Walnuts (Akhrot)", 350, 430, "250g", "dry-fruits-and-nuts", "Organic Wala", 60, 250, "Premium Kashmiri akhrot halves", "walnuts,akhrot,kashmiri"},
		{"Pistachios (Pista)", 380, 460, "250g", "dry-fruits-and-nuts", "Organic Wala", 60, 250, "Salted and roasted Iranian pistachios", "pista,pistachios,roasted"},
		{"Dates Medjool (Khajoor)", 450, 550, "500g", "dry-fruits-and-nuts", "Organic Wala", 50, 500, "Jumbo Medjool dates from Saudi Arabia", "dates,khajoor,medjool"},
		{"Anjeer (Dried Figs)", 280, 350, "200g", "dry-fruits-and-nuts", "Organic Wala", 60, 200, "Soft dried figs from Afghanistan", "anjeer,figs,dried"},
		{"Mixed Dry Fruits Pack", 500, 600, "500g", "dry-fruits-and-nuts", "Organic Wala", 40, 500, "Premium mix of almonds, cashews, raisins and pistachios", "mixed,dry-fruits,gift-pack"},

		// ── Instant & Ready to Eat ───────────────────────────────
		{"Maggi 2-Minute Noodles", 14, 0, "70g", "instant-and-ready-to-eat", "Annapurna Grocery", 300, 70, "India's favourite instant masala noodles", "maggi,noodles,instant"},
		{"MTR Ready to Eat Dal Makhani", 90, 110, "300g", "instant-and-ready-to-eat", "Annapurna Grocery", 80, 300, "Authentic restaurant-style dal makhani, just heat and eat", "mtr,dal-makhani,ready-to-eat"},
		{"MTR Rava Idli Mix", 55, 70, "500g", "instant-and-ready-to-eat", "Annapurna Grocery", 80, 500, "Instant rava idli mix - soft idlis in minutes", "mtr,idli,instant-mix"},
		{"Gits Gulab Jamun Mix", 85, 100, "200g", "instant-and-ready-to-eat", "Annapurna Grocery", 60, 200, "Easy gulab jamun mix for perfect round jamuns", "gits,gulab-jamun,sweet"},
		{"MTR Bisibelebath Paste", 65, 80, "200g", "instant-and-ready-to-eat", "Annapurna Grocery", 60, 200, "Karnataka's favourite spicy rice dish paste", "mtr,bisibelebath,karnataka"},
		{"Knorr Tomato Soup", 45, 55, "53g", "instant-and-ready-to-eat", "Annapurna Grocery", 100, 53, "Classic tomato soup, serves 4", "knorr,soup,tomato"},

		// ── Personal Care ────────────────────────────────────────
		{"Dettol Original Soap", 40, 50, "125g", "personal-care", "Annapurna Grocery", 200, 125, "Trusted antibacterial bath soap", "dettol,soap,antibacterial"},
		{"Dove Cream Beauty Bar", 55, 65, "100g", "personal-care", "Annapurna Grocery", 150, 100, "Moisturising beauty bathing bar", "dove,soap,moisturising"},
		{"Head & Shoulders Shampoo", 180, 220, "340ml", "personal-care", "Annapurna Grocery", 80, 340, "Anti-dandruff shampoo with cool menthol", "shampoo,head-shoulders,dandruff"},
		{"Colgate MaxFresh Toothpaste", 95, 115, "150g", "personal-care", "Annapurna Grocery", 150, 150, "Cooling crystal toothpaste with whitening", "colgate,toothpaste,maxfresh"},
		{"Himalaya Neem Face Wash", 135, 165, "150ml", "personal-care", "Annapurna Grocery", 80, 150, "Purifying neem face wash for pimple-free skin", "himalaya,face-wash,neem"},

		// ── Cleaning & Household ─────────────────────────────────
		{"Surf Excel Easy Wash", 120, 145, "1kg", "cleaning-and-household", "Annapurna Grocery", 100, 1000, "Detergent powder for tough stain removal", "surf-excel,detergent,washing"},
		{"Vim Dishwash Bar", 30, 0, "200g", "cleaning-and-household", "Annapurna Grocery", 200, 200, "Lemon fresh dishwash bar", "vim,dishwash,utensil"},
		{"Lizol Floor Cleaner", 120, 145, "500ml", "cleaning-and-household", "Annapurna Grocery", 80, 500, "Disinfectant floor cleaner, citrus fragrance", "lizol,floor-cleaner,disinfectant"},
		{"Scotch-Brite Scrub Pad", 30, 0, "3 pcs", "cleaning-and-household", "Annapurna Grocery", 150, 50, "Pack of 3 kitchen scrub pads", "scotch-brite,scrub,kitchen"},

		// ── Organic Specials (Organic Wala) ──────────────────────
		{"Organic Toor Dal", 180, 220, "1kg", "dals-and-pulses", "Organic Wala", 60, 1000, "Certified organic toor dal, pesticide-free", "organic,toor,dal"},
		{"Organic Jaggery (Gur)", 90, 110, "500g", "other-grains", "Organic Wala", 80, 500, "Chemical-free sugarcane jaggery from Maharashtra", "jaggery,gur,organic"},
		{"Cold-Pressed Coconut Oil", 350, 420, "1L", "cooking-oil-and-ghee", "Organic Wala", 50, 1000, "100% organic cold-pressed virgin coconut oil", "coconut-oil,organic,virgin"},
		{"A2 Cow Ghee (Bilona)", 850, 1000, "500ml", "cooking-oil-and-ghee", "Organic Wala", 40, 500, "Traditional bilona method A2 cow ghee", "ghee,a2,bilona,organic"},
		{"Organic Honey", 320, 400, "500g", "dry-fruits-and-nuts", "Organic Wala", 60, 500, "Raw unprocessed organic honey from Himalayan farms", "honey,organic,raw,himalayan"},
		{"Organic Turmeric Powder", 95, 120, "200g", "masala-and-spices", "Organic Wala", 80, 200, "High curcumin organic turmeric from Lakadong", "turmeric,organic,lakadong"},
	}

	inserted := 0
	for _, p := range products {
		catID, ok := catIDs[p.cat]
		if !ok {
			log.Printf("SKIP %s: category %s not found", p.name, p.cat)
			continue
		}
		vendorID, ok := vendorIDs[p.vendor]
		if !ok {
			log.Printf("SKIP %s: vendor %s not found", p.name, p.vendor)
			continue
		}

		slug := slugify(p.name)
		tags := "{" + p.tags + "}"

		var comparePrice interface{} = nil
		if p.compare > 0 {
			comparePrice = p.compare
		}

		_, err := db.Exec(`
			INSERT INTO products
				(vendor_id, category_id, name, slug, description, price, compare_price,
				 unit, stock_quantity, low_stock_threshold, weight_grams, tags, images, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 5, $10, $11, '{}', true)
			ON CONFLICT (vendor_id, slug) DO UPDATE
				SET price = EXCLUDED.price, stock_quantity = EXCLUDED.stock_quantity,
				    description = EXCLUDED.description, compare_price = EXCLUDED.compare_price
		`, vendorID, catID, p.name, slug, p.desc, p.price, comparePrice,
			p.unit, p.stock, p.weight, tags)
		if err != nil {
			log.Printf("product %s: %v", p.name, err)
		} else {
			inserted++
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// Summary
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n══════════════════════════════════════════════════════")
	fmt.Printf("  ✅  %d categories created\n", len(categories))
	fmt.Printf("  ✅  %d vendors created\n", len(vendors))
	fmt.Printf("  ✅  %d products inserted\n", inserted)
	fmt.Println("══════════════════════════════════════════════════════")
	fmt.Println("\n  Vendor login credentials:")
	for _, v := range vendors {
		fmt.Printf("    %-25s  %s / Vendor@123\n", v.biz, v.email)
	}
	fmt.Println()
}
