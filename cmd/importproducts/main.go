// Command importproducts loads or updates the product catalogue from a CSV file.
//
// There is no admin panel, so this is how the shop owner adds or edits products:
// edit catalog.csv in a spreadsheet, then run the importer. Products are matched by
// name (case-insensitive): existing ones are updated, new ones are created.
//
//	go run ./cmd/importproducts -file catalog.csv
package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var required = []string{"name", "category", "price", "skin_type", "image_url", "description", "volume"}

type product struct {
	name, category, skinType, badge, imageURL, description, volume string
	price, compareAt, reviewCount                                  int
	hasCompare                                                     bool
	rating                                                         float64
}

func main() {
	file := flag.String("file", "catalog.csv", "path to the CSV file")
	dryRun := flag.Bool("dry-run", false, "validate the file without touching the database")
	flag.Parse()
	_ = godotenv.Load()

	f, err := os.Open(*file)
	if err != nil {
		log.Fatalf("open %s: %v", *file, err)
	}
	defer f.Close()

	products, err := parse(f)
	if err != nil {
		log.Fatalf("invalid CSV: %v", err)
	}
	fmt.Printf("parsed %d products\n", len(products))
	if *dryRun {
		return
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(ctx)

	for _, p := range products {
		var compare any
		if p.hasCompare {
			compare = p.compareAt
		}
		var badge any
		if p.badge != "" {
			badge = p.badge
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO products (name, category, price, compare_at_price, rating, review_count, skin_type, badge, image_url, description, volume)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			ON CONFLICT ((lower(name))) DO UPDATE SET
				category = EXCLUDED.category, price = EXCLUDED.price, compare_at_price = EXCLUDED.compare_at_price,
				rating = EXCLUDED.rating, review_count = EXCLUDED.review_count, skin_type = EXCLUDED.skin_type,
				badge = EXCLUDED.badge, image_url = EXCLUDED.image_url, description = EXCLUDED.description,
				volume = EXCLUDED.volume
		`, p.name, p.category, p.price, compare, p.rating, p.reviewCount, p.skinType, badge, p.imageURL, p.description, p.volume)
		if err != nil {
			log.Fatalf("save %q: %v", p.name, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("imported %d products\n", len(products))
}

// parse validates the header and every row, reporting the first problem with its line number.
func parse(r io.Reader) ([]product, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true
	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	col := map[string]int{}
	for i, h := range header {
		col[strings.ToLower(strings.TrimSpace(h))] = i
	}
	for _, name := range required {
		if _, ok := col[name]; !ok {
			return nil, fmt.Errorf("missing required column %q", name)
		}
	}

	get := func(row []string, name string) string {
		i, ok := col[name]
		if !ok || i >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[i])
	}

	var out []product
	line := 1
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		line++
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}

		p := product{
			name: get(row, "name"), category: get(row, "category"), skinType: get(row, "skin_type"),
			badge: get(row, "badge"), imageURL: get(row, "image_url"), description: get(row, "description"),
			volume: get(row, "volume"),
		}
		for _, f := range required {
			if get(row, f) == "" {
				return nil, fmt.Errorf("line %d: %q is empty", line, f)
			}
		}
		if p.price, err = strconv.Atoi(get(row, "price")); err != nil || p.price <= 0 {
			return nil, fmt.Errorf("line %d: price must be a whole number of naira above zero", line)
		}
		if v := get(row, "compare_at_price"); v != "" {
			if p.compareAt, err = strconv.Atoi(v); err != nil || p.compareAt <= 0 {
				return nil, fmt.Errorf("line %d: compare_at_price must be a whole number above zero", line)
			}
			p.hasCompare = true
		}
		if v := get(row, "rating"); v != "" {
			if p.rating, err = strconv.ParseFloat(v, 64); err != nil || p.rating < 0 || p.rating > 5 {
				return nil, fmt.Errorf("line %d: rating must be between 0 and 5", line)
			}
		}
		if v := get(row, "review_count"); v != "" {
			if p.reviewCount, err = strconv.Atoi(v); err != nil || p.reviewCount < 0 {
				return nil, fmt.Errorf("line %d: review_count must be 0 or more", line)
			}
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("the file has no product rows")
	}
	return out, nil
}
