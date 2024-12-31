package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Migration struct {
	Id        int32     `db:"id"`
	Name      string    `db:"name"`
	AppliedOn time.Time `db:"applied_on"`
}

func (this *Migration) MigrationFileName() string {
	return this.Name + ".sql"
}

const MigrationDir = "migrations"

const InitialMigrationCheckQuery = `
SELECT EXISTS(
    SELECT * FROM pg_tables WHERE schemaname='public' AND tablename='migrations'
);
`

const LatestMigrationQuery = `
SELECT id, name, applied_on
FROM migrations
ORDER BY id DESC
LIMIT 1;
`

const InsertMigration = `
INSERT INTO migrations(id, name) VALUES ($1, $2) returning *;
`

func getLastMigration(ctx context.Context, dbpool *pgxpool.Pool) (*Migration, error) {
	rows, _ := dbpool.Query(ctx, InitialMigrationCheckQuery)
	isDbInitiated, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return nil, err
	}

	if !isDbInitiated {
		return nil, nil
	}

	rows, _ = dbpool.Query(ctx, LatestMigrationQuery)
	migration, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Migration])
	return migration, nil
}

func stripFileExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func applyMigration(ctx context.Context, dbpool *pgxpool.Pool, fileName string, migrationNum int) *Migration {
	filePath := filepath.Join(".", MigrationDir, fileName)
	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error while reading initial migration file: %v", err)
	}
	// Convert byte array to string
	query := string(b)
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(ctx)
	// Execute the DDL command for creating initial migration

	_, err = tx.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Error while executing initial migration query: %v", err)
	}

	rows, _ := tx.Query(ctx, InsertMigration, migrationNum, stripFileExtension(fileName))
	migration, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Migration])
	if err != nil {
		log.Fatalf("Error while executing create query for initial migration: %v", err)
	}
	err = tx.Commit(ctx)

	log.Printf("Successfully applied migration from %s", fileName)

	return migration
}

func applyInitialMigration(ctx context.Context, dbpool *pgxpool.Pool) *Migration {
	initialMigrationFileName := "migration_0000.sql"
	return applyMigration(ctx, dbpool, initialMigrationFileName, 0)
}

func applyPendingMigrations(ctx context.Context, dbpool *pgxpool.Pool, lastMigration *Migration) {

	files, err := os.ReadDir(filepath.Join(".", MigrationDir))
	if err != nil {
		log.Fatalf("Error while fetching migration files: %v\n", err)
	}
	log.Printf("Found %d migration files in total\n", len(files))

	pattern := regexp.MustCompile(`^.*_(\d{4})\.sql$`)

	pendingMigrationNumToFileName := make(map[int]string)
	var pendingMigrationNums []int

	for _, file := range files {
		fileName := file.Name()
		matches := pattern.FindStringSubmatch(fileName)
		if len(matches) != 2 {
			log.Fatalf("Invalid migration file name: %s", fileName)
		}
		// No need to check for err here because the number is fetched from regex, so must be valid
		migrationNum, _ := strconv.Atoi(matches[1])
		if migrationNum > int(lastMigration.Id) {
			pendingMigrationNumToFileName[migrationNum] = fileName
			pendingMigrationNums = append(pendingMigrationNums, migrationNum)
		}
	}

	if len(pendingMigrationNums) == 0 {
		log.Println("No new migratiosn to apply, exiting...")
		return
	}
	sort.Ints(pendingMigrationNums)

	for _, num := range pendingMigrationNums {
		pendingMigrationFileName := pendingMigrationNumToFileName[num]
		log.Printf("Pending migration file: %s\n", pendingMigrationFileName)
		applyMigration(ctx, dbpool, pendingMigrationFileName, num)
	}

}

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URI"))
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
	defer dbpool.Close()
	lastMigration, err := getLastMigration(ctx, dbpool)
	if err != nil {
		log.Fatalf("Error initiating migrations: %v\n", err)
	}

	if lastMigration == nil {
		log.Println("No previously applied migration. Initial migration will be applied first")
		lastMigration = applyInitialMigration(ctx, dbpool)

	} else {
		log.Printf("Last migration | id: %d | name: %v\n", lastMigration.Id, lastMigration.Name)
	}
	applyPendingMigrations(ctx, dbpool, lastMigration)

}
