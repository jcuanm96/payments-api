package repositories

import (
	"context"
	"fmt"

	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) React(ctx context.Context, reaction feed.Reaction, tx pgx.Tx) (*feed.ReactionMetadata, error) {
	params := reaction.GetInsertParams()
	columns, values := reaction.GetUpdateColumnsValues()

	// Query Composition
	updateQuery := fmt.Sprintf(`
		UPDATE 
			feed.posts 
		SET %s = %s -- ex: SET num_upvotes = num_upvotes-1; ex2: SET (num_upvotes, num_downvotes) = (num_upvotes-1, num_downvotes+1)
		WHERE id = $1
		RETURNING num_upvotes, num_downvotes;`,
		columns,
		values,
	)
	insertQuery := `
		INSERT INTO 
			feed.post_reactions (
				post_id, 
				user_id, 
				type
			) 
		VALUES (
			$1, 
			$2, 
			$3
		) 
		ON CONFLICT (
			post_id, 
			user_id
		) 
		DO UPDATE SET type = $4;
	`

	row, postUpdateErr := tx.Query(ctx, updateQuery, reaction.PostID)
	if postUpdateErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(postUpdateErr, updateQuery, []interface{}{reaction.PostID}))
		return nil, postUpdateErr
	}

	defer row.Close()
	if !row.Next() {
		return nil, fmt.Errorf("no rows returned, expected number of up/down votes. query: %s", updateQuery)
	}
	var metadata feed.ReactionMetadata
	scanErr := row.Scan(
		&metadata.UpVotes,
		&metadata.DownVotes,
	)

	row.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	_, insertErr := tx.Exec(ctx, insertQuery, params...)
	if insertErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(insertErr, insertQuery, params))
		return nil, insertErr
	}

	return &metadata, nil
}
