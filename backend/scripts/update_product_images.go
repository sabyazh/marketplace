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

	// Unsplash image URLs (curated from search results)
	// Format: product name substring → image URL
	type mapping struct {
		nameContains string
		imageURL     string
	}

	const qs = "?w=400&h=400&fit=crop&auto=format&q=80"

	mappings := []mapping{
		// ── Fresh Vegetables ─────────────────────────────────────
		{"Tomato", "https://images.unsplash.com/photo-1560433802-62c9db426a4d" + qs},
		{"Onion", "https://images.unsplash.com/photo-1585849834908-3481231155e8" + qs},
		{"Potato (Aloo)", "https://plus.unsplash.com/premium_photo-1724256031338-b5bfba816cfd" + qs},
		{"Green Chilli", "https://plus.unsplash.com/premium_photo-1700028097487-87e40431bd3a" + qs},
		{"Ginger", "https://plus.unsplash.com/premium_photo-1675364893053-180a3c6e0119" + qs},
		{"Garlic", "https://plus.unsplash.com/premium_photo-1675731118551-79b3da05a5d4" + qs},
		{"Coriander Leaves", "https://plus.unsplash.com/premium_photo-1763467051482-ce108eb2d727" + qs},
		{"Lady Finger", "https://plus.unsplash.com/premium_photo-1666877059056-f42ada662ccc" + qs},
		{"Brinjal", "https://plus.unsplash.com/premium_photo-1666270423836-864dfa7071e5" + qs},
		{"Cauliflower", "https://plus.unsplash.com/premium_photo-1711684803510-6f05fa515378" + qs},
		{"Spinach", "https://plus.unsplash.com/premium_photo-1701699257548-8261a687236f" + qs},
		{"Drumstick", "https://images.unsplash.com/photo-1771643033515-0028fd03b708" + qs},

		// ── Fresh Fruits ──────────────────────────────────────────
		{"Banana", "https://plus.unsplash.com/premium_photo-1724250081102-cab0e5cb314c" + qs},
		{"Apple", "https://plus.unsplash.com/premium_photo-1724249989963-9286e126af81" + qs},
		{"Mango", "https://plus.unsplash.com/premium_photo-1724255863045-2ad716767715" + qs},
		{"Papaya", "https://plus.unsplash.com/premium_photo-1722938907181-08d806f7b9a6" + qs},
		{"Pomegranate", "https://plus.unsplash.com/premium_photo-1668076515507-c5bc223c99a4" + qs},
		{"Guava", "https://plus.unsplash.com/premium_photo-1675040829892-d6de2f0fb422" + qs},

		// ── Milk & Curd ──────────────────────────────────────────
		{"Milk", "https://plus.unsplash.com/premium_photo-1695042864936-b28290dd54b4" + qs},
		{"Curd", "https://plus.unsplash.com/premium_photo-1666275003961-67e497ad41c6" + qs},
		{"Buttermilk", "https://plus.unsplash.com/premium_photo-1666275003961-67e497ad41c6" + qs},

		// ── Paneer & Cheese ──────────────────────────────────────
		{"Paneer", "https://plus.unsplash.com/premium_photo-1723730426108-1bb37a500d5c" + qs},
		{"Cheese", "https://plus.unsplash.com/premium_photo-1691939610797-aba18030c15f" + qs},
		{"Cream Cheese", "https://plus.unsplash.com/premium_photo-1691939610797-aba18030c15f" + qs},

		// ── Bread & Eggs ─────────────────────────────────────────
		{"Bread", "https://plus.unsplash.com/premium_photo-1668772632934-9f00276753ad" + qs},
		{"Eggs", "https://plus.unsplash.com/premium_photo-1676686126965-cb536e2328c3" + qs},
		{"Pav", "https://plus.unsplash.com/premium_photo-1689247409836-a9d05ca56d35" + qs},

		// ── Rice ──────────────────────────────────────────────────
		{"Rice", "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21" + qs},
		{"Ponni", "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21" + qs},

		// ── Wheat & Grains ───────────────────────────────────────
		{"Atta", "https://plus.unsplash.com/premium_photo-1666174323664-9733133a6f1c" + qs},
		{"Sooji", "https://plus.unsplash.com/premium_photo-1725878608875-ca527315bb46" + qs},
		{"Poha", "https://plus.unsplash.com/premium_photo-1700750447788-9b54d4ba6734" + qs},
		{"Maida", "https://plus.unsplash.com/premium_photo-1666174323664-9733133a6f1c" + qs},
		{"Besan", "https://plus.unsplash.com/premium_photo-1666174323664-9733133a6f1c" + qs},
		{"Ragi", "https://plus.unsplash.com/premium_photo-1666174323664-9733133a6f1c" + qs},

		// ── Dals & Pulses ────────────────────────────────────────
		{"Toor Dal", "https://plus.unsplash.com/premium_photo-1701064865162-db655bfb99c3" + qs},
		{"Moong Dal", "https://plus.unsplash.com/premium_photo-1674025751520-01102eca6352" + qs},
		{"Chana Dal", "https://plus.unsplash.com/premium_photo-1674025748666-9ac6c8a7ab3d" + qs},
		{"Masoor", "https://plus.unsplash.com/premium_photo-1694506374757-632d818eb023" + qs},
		{"Urad Dal", "https://plus.unsplash.com/premium_photo-1701064865162-db655bfb99c3" + qs},
		{"Rajma", "https://images.unsplash.com/photo-1733562256707-c13261bfd0a9" + qs},
		{"Kabuli Chana", "https://plus.unsplash.com/premium_photo-1675237624857-7d995e29897d" + qs},
		{"Moong Whole", "https://plus.unsplash.com/premium_photo-1674025751520-01102eca6352" + qs},

		// ── Masala & Spices ──────────────────────────────────────
		{"Garam Masala", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},
		{"Turmeric", "https://plus.unsplash.com/premium_photo-1726862790171-0d6208559224" + qs},
		{"Haldi", "https://plus.unsplash.com/premium_photo-1726862790171-0d6208559224" + qs},
		{"Red Chilli Powder", "https://plus.unsplash.com/premium_photo-1726880501641-c7072313efd2" + qs},
		{"Kashmiri Red Chilli", "https://plus.unsplash.com/premium_photo-1726880501641-c7072313efd2" + qs},
		{"Coriander Powder", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},
		{"Cumin", "https://plus.unsplash.com/premium_photo-1722686499744-59e1bcf902a6" + qs},
		{"Jeera", "https://plus.unsplash.com/premium_photo-1722686499744-59e1bcf902a6" + qs},
		{"Black Pepper", "https://plus.unsplash.com/premium_photo-1668446314011-301c7a98b6a9" + qs},
		{"Mustard Seeds", "https://plus.unsplash.com/premium_photo-1722686499744-59e1bcf902a6" + qs},
		{"Cloves", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},
		{"Cinnamon", "https://plus.unsplash.com/premium_photo-1666174330841-7575f42c7cd9" + qs},
		{"Cardamom", "https://plus.unsplash.com/premium_photo-1669842468224-364235e9c9ef" + qs},
		{"Kitchen King", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},
		{"Chhole Masala", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},
		{"Sambar Powder", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},
		{"Rasam Powder", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},
		{"Asafoetida", "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs},

		// ── Cooking Oil & Ghee ───────────────────────────────────
		{"Sunflower Oil", "https://plus.unsplash.com/premium_photo-1667818824583-3be6f268bb13" + qs},
		{"Saffola", "https://plus.unsplash.com/premium_photo-1667818824583-3be6f268bb13" + qs},
		{"Mustard Oil", "https://plus.unsplash.com/premium_photo-1667818824583-3be6f268bb13" + qs},
		{"Ghee", "https://plus.unsplash.com/premium_photo-1673965777407-9ae155fde644" + qs},
		{"Groundnut Oil", "https://plus.unsplash.com/premium_photo-1667818824583-3be6f268bb13" + qs},
		{"Coconut Oil", "https://plus.unsplash.com/premium_photo-1661454028102-e949326751a3" + qs},

		// ── Snacks & Namkeen ─────────────────────────────────────
		{"Bhujia", "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa" + qs},
		{"Moong Dal", "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa" + qs},
		{"Lay's", "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa" + qs},
		{"Kurkure", "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa" + qs},
		{"Parle-G", "https://plus.unsplash.com/premium_photo-1667621221004-e344ae82ad7e" + qs},
		{"Good Day", "https://plus.unsplash.com/premium_photo-1667621221004-e344ae82ad7e" + qs},
		{"Sev", "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa" + qs},
		{"Balaji", "https://plus.unsplash.com/premium_photo-1672753747124-2bd4da9931fa" + qs},

		// ── Beverages ────────────────────────────────────────────
		{"Tea", "https://plus.unsplash.com/premium_photo-1672076780874-82fe06198ef7" + qs},
		{"Coffee", "https://images.unsplash.com/photo-1556742526-795a8eac090e" + qs},
		{"Tropicana", "https://plus.unsplash.com/premium_photo-1666262370578-59f5736c0202" + qs},
		{"Aam Panna", "https://plus.unsplash.com/premium_photo-1666262370578-59f5736c0202" + qs},
		{"Rooh Afza", "https://plus.unsplash.com/premium_photo-1666262370578-59f5736c0202" + qs},

		// ── Dry Fruits & Nuts ────────────────────────────────────
		{"Almonds", "https://plus.unsplash.com/premium_photo-1675237625910-e5d354c03987" + qs},
		{"Cashew", "https://plus.unsplash.com/premium_photo-1723978744235-cd52cec7ad38" + qs},
		{"Kishmish", "https://plus.unsplash.com/premium_photo-1669205434519-a042ba09fbdd" + qs},
		{"Walnuts", "https://plus.unsplash.com/premium_photo-1668445743008-0d84ffc370e0" + qs},
		{"Pistachios", "https://plus.unsplash.com/premium_photo-1725874816737-8918f48c41a4" + qs},
		{"Dates", "https://plus.unsplash.com/premium_photo-1676208753932-6e8bc83a0b0d" + qs},
		{"Anjeer", "https://plus.unsplash.com/premium_photo-1676208753932-6e8bc83a0b0d" + qs},
		{"Mixed Dry Fruits", "https://plus.unsplash.com/premium_photo-1675237625910-e5d354c03987" + qs},
		{"Honey", "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9" + qs},

		// ── Instant & Ready to Eat ───────────────────────────────
		{"Maggi", "https://plus.unsplash.com/premium_photo-1738105946995-5cb9d1a8656b" + qs},
		{"Noodles", "https://plus.unsplash.com/premium_photo-1738105946995-5cb9d1a8656b" + qs},
		{"Dal Makhani", "https://plus.unsplash.com/premium_photo-1701064865162-db655bfb99c3" + qs},
		{"Idli Mix", "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21" + qs},
		{"Gulab Jamun", "https://plus.unsplash.com/premium_photo-1677956021545-986eaa5f6e65" + qs},
		{"Bisibelebath", "https://plus.unsplash.com/premium_photo-1694141252026-3df1de888a21" + qs},
		{"Soup", "https://plus.unsplash.com/premium_photo-1666262370578-59f5736c0202" + qs},

		// ── Personal Care ────────────────────────────────────────
		{"Dettol", "https://plus.unsplash.com/premium_photo-1677776518705-70b21cbc4d47" + qs},
		{"Dove", "https://plus.unsplash.com/premium_photo-1677776518705-70b21cbc4d47" + qs},
		{"Soap", "https://plus.unsplash.com/premium_photo-1677776518705-70b21cbc4d47" + qs},
		{"Shampoo", "https://plus.unsplash.com/premium_photo-1679106764781-1387196b5945" + qs},
		{"Toothpaste", "https://plus.unsplash.com/premium_photo-1676909663912-b251c1468a74" + qs},
		{"Face Wash", "https://plus.unsplash.com/premium_photo-1679106764781-1387196b5945" + qs},

		// ── Cleaning & Household ─────────────────────────────────
		{"Surf Excel", "https://plus.unsplash.com/premium_photo-1664372899356-3d9dc566fba2" + qs},
		{"Detergent", "https://plus.unsplash.com/premium_photo-1664372899356-3d9dc566fba2" + qs},
		{"Vim", "https://plus.unsplash.com/premium_photo-1664443944738-357d3f7d1c42" + qs},
		{"Lizol", "https://plus.unsplash.com/premium_photo-1664443944738-357d3f7d1c42" + qs},
		{"Scotch-Brite", "https://plus.unsplash.com/premium_photo-1664443944738-357d3f7d1c42" + qs},

		// ── Organic Specials ─────────────────────────────────────
		{"Organic Toor", "https://plus.unsplash.com/premium_photo-1701064865162-db655bfb99c3" + qs},
		{"Organic Jaggery", "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9" + qs},
		{"Cold-Pressed", "https://plus.unsplash.com/premium_photo-1661454028102-e949326751a3" + qs},
		{"A2 Cow Ghee", "https://plus.unsplash.com/premium_photo-1673965777407-9ae155fde644" + qs},
		{"Organic Honey", "https://plus.unsplash.com/premium_photo-1664273586932-ab870d61f7e9" + qs},
		{"Organic Turmeric", "https://plus.unsplash.com/premium_photo-1726862790171-0d6208559224" + qs},
	}

	updated := 0
	for _, m := range mappings {
		result, err := db.Exec(`
			UPDATE products
			SET images = ARRAY[$1]
			WHERE LOWER(name) LIKE LOWER($2)
			AND (images IS NULL OR images = '{}' OR array_length(images, 1) IS NULL)
		`, m.imageURL, "%"+m.nameContains+"%")
		if err != nil {
			log.Printf("  ✗ %s: %v", m.nameContains, err)
			continue
		}
		rows, _ := result.RowsAffected()
		if rows > 0 {
			updated += int(rows)
			fmt.Printf("  ✓ %-30s  %d products updated\n", m.nameContains, rows)
		}
	}

	// Catch any remaining products without images — assign a generic grocery image
	genericImage := "https://plus.unsplash.com/premium_photo-1692776206795-60a58a4dc817" + qs
	result, err := db.Exec(`
		UPDATE products
		SET images = ARRAY[$1]
		WHERE images IS NULL OR images = '{}' OR array_length(images, 1) IS NULL
	`, genericImage)
	if err == nil {
		rows, _ := result.RowsAffected()
		if rows > 0 {
			updated += int(rows)
			fmt.Printf("  ✓ %-30s  %d products (generic)\n", "Remaining", rows)
		}
	}

	// Verify
	var total, withImages int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM products WHERE images IS NOT NULL AND array_length(images, 1) > 0").Scan(&withImages)

	fmt.Println("\n══════════════════════════════════════════════════════")
	fmt.Printf("  ✅  %d products updated with images\n", updated)
	fmt.Printf("  📊  %d/%d products now have images\n", withImages, total)
	_ = strings.ToLower // keep import used
	fmt.Println("══════════════════════════════════════════════════════")
}
