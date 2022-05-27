package test_util

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
	"github.com/VamaSingapore/vama-api/internal/token"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/scrypt"
)

func InitializeDbSchemas(db *pgxpool.Pool) {
	// The order of these files matters as some schemas depend on
	// other schemas already being defined.
	schemaFiles := []string{
		"updated_at_trigger.sql",
		"core_schema.sql",
		"feed_schema.sql",
		"subscription_schema.sql",
		"cache_schema.sql",
		"sharing_schema.sql",
		"push_schema.sql",
	}

	ctx := context.Background()

	for _, schemaFile := range schemaFiles {
		// assumes we are in project root dir
		path := filepath.Join("migrations/", schemaFile)

		c, ioErr := ioutil.ReadFile(path)
		if ioErr != nil {
			vlog.Fatalf(ctx, "io error for schema file %s: %s", schemaFile, ioErr)
		}
		sql := string(c)
		_, sqlErr := db.Exec(context.Background(), sql)
		if sqlErr != nil {
			vlog.Fatalf(ctx, "sql error for schema file %s: %s", schemaFile, sqlErr)
		}
	}
}

func ClearDB(db *pgxpool.Pool) {
	if appconfig.Config.Gcloud.Project != "vama-test" {
		panic("Running ClearDB function (DANGEROUS) in unauthorized environment.")
	}
	ctx := context.Background()
	res, showErr := db.Query(ctx,
		`
	SELECT table_schema, table_name
	FROM information_schema.tables
	WHERE 
	   (
		table_schema = 'public' OR
	    table_schema = 'core' OR
	    table_schema = 'feed' OR
	    table_schema = 'subscription' OR
	    table_schema = 'wallet' OR
	    table_schema = 'product'
	   )
	   AND table_type = 'BASE TABLE'
	ORDER BY table_schema;
	`)
	if showErr != nil {
		vlog.Fatalf(ctx, "show tables error: %s", showErr)
	}
	defer res.Close()

	// Foreign key relationships mean we have to drop certain tables
	// in a specific order. This list contains the ones that need
	// to be dropped first.

	tables := []string{
		"subscription.paid_group_chat_subscriptions",
		"product.paid_group_chats",
		"core.goat_invite_codes",
		"product.goat_chats",
		"subscription.tiers",
		"core.users_contacts",
		"feed.post_comments",
		"feed.post_reactions",
		"feed.posts",
		"wallet.balances",
		"wallet.ledger",
		"feed.follows",
		"core.user_blocks",
	}
	for _, table := range tables {
		deleteTable(db, table)
	}

	for res.Next() {
		var table string
		var schema string
		res.Scan(&schema, &table)
		fullTableName := fmt.Sprintf("%s.%s", schema, table)
		deleteTable(db, fullTableName)
	}
	resetSequences(db)
}

// fullTableName has format {schema}.{table name}
func deleteTable(db *pgxpool.Pool, fullTableName string) {
	query := fmt.Sprintf("DELETE from %s;", fullTableName)
	ctx := context.Background()
	_, deleteErr := db.Exec(ctx, query)
	if deleteErr != nil {
		vlog.Errorf(ctx, "sql error deleting from table %s: %s", fullTableName, deleteErr)
		vlog.Fatalf(ctx, "erroring query: %s", query)
	}
}

// Deleting a database does not reset the value of SERIAL fields
// such as ID. This function enumerates all sequences in db and resets
// them to their initial values.
func resetSequences(db *pgxpool.Pool) {
	sequencesQuery := `
		WITH seqs AS (	
			SELECT
			pg_get_serial_sequence(t.schemaname || '.' || t.tablename, c.column_name) as seq
		FROM pg_tables t
		JOIN information_schema.columns c ON
			c.table_schema = t.schemaname AND
			c.table_name = t.tablename
		WHERE
			t.schemaname <> 'pg_catalog' AND
			t.schemaname <> 'information_schema' AND
			pg_get_serial_sequence(t.schemaname || '.' || t.tablename, c.column_name) IS NOT NULL
		)
		SELECT
			seq,
			SETVAL(seq, 1, false)
		FROM seqs
	`
	ctx := context.Background()
	_, sequencesErr := db.Exec(ctx, sequencesQuery)
	if sequencesErr != nil {
		vlog.Fatalf(ctx, "sql error on retrieving sequences: %v \nerroring query: %s", sequencesErr, sequencesQuery)
	}
}

func CheckIsValidToken(tokenSvc token.Service, tokenClaims *token.Claims, tokenID string, tokenType string) (bool, error) {
	authTkn := strings.TrimPrefix(tokenID, "Bearer ")

	costParameter := 16384
	r := 8
	p := 1
	keyLen := 32
	tokenHash, hashErr := scrypt.Key([]byte(authTkn), []byte(tokenClaims.Data.UUID), costParameter, r, p, keyLen)
	if hashErr != nil {
		return false, hashErr
	}

	isTokenValid := new(bool)
	baserepo.GetFromCache(context.Background(), tokenSvc.MasterNode(), tokenSvc.GetRedisClient(), tokenType, fmt.Sprintf("%x", tokenHash), &isTokenValid)

	if isTokenValid == nil || !*isTokenValid {
		return false, nil
	}

	return *isTokenValid, nil
}
