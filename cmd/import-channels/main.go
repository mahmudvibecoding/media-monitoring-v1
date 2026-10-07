package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"bufio"
	"bytes"
	"context"
	"time"
	"github.com/jackc/pgx/v5"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	file, err := os.Open("/data/youtube-channels.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	buffered := bufio.NewReader(file)

	prefix, _ := buffered.Peek(3)

	if bytes.Equal(prefix, []byte{0xEF, 0xBB, 0xBF}) {
		if _, err := buffered.Discard(3); err != nil {
			return fmt.Errorf("skip UTF-8 BOM: %w", err)
		}
	}

	reader := csv.NewReader(buffered)

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	idColumn := -1
	for i, name := range header {
		if strings.TrimSpace(name) == "Channel ID" {
			if idColumn != -1 {
				return fmt.Errorf("duplicate Channel ID column")
			}
			idColumn = i
		}
	}
	if idColumn == -1 {
		return fmt.Errorf("missing Channel ID column")
	}

	validID := regexp.MustCompile(`^UC[A-Za-z0-9_-]{22}$`)
	seen := make(map[string]struct{})
	rows, duplicates, invalid := 0, 0, 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read CSV after %d data rows: %w", rows, err)
		}
		rows++

		id := strings.TrimSpace(record[idColumn])
		if !validID.MatchString(id) {
			invalid++
			if invalid <= 5 {
				fmt.Printf("Invalid ID in data row %d: %q\n", rows, id)
			}
			continue
		}

		if _, exists := seen[id]; exists {
			duplicates++
			continue
		}
		seen[id] = struct{}{}
	}

	fmt.Printf(
		"Data rows: %d\nUnique valid IDs: %d\nDuplicate valid IDs: %d\nInvalid IDs: %d\n",
		rows, len(seen), duplicates, invalid,
	)

	if invalid > 0 {
		return fmt.Errorf("found %d invalid IDs; nothing imported", invalid)
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}

	inserted, err := saveChannels(ids)
	if err != nil {
		return err
	}

	fmt.Printf("Inserted: %d\nAlready stored: %d\n",
		inserted, int64(len(ids))-inserted)

	return nil
}

func saveChannels(ids []string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	conn, err := pgx.Connect(ctx, "")
	if err != nil {
		return 0, fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	defer conn.Close(context.Background())

	result, err := conn.Exec(ctx, `
		INSERT INTO public.youtube_channels (channel_id)
		SELECT unnest($1::text[])
		ON CONFLICT (channel_id) DO NOTHING
	`, ids)
	if err != nil {
		return 0, fmt.Errorf("insert channels: %w", err)
	}

	return result.RowsAffected(), nil
}
