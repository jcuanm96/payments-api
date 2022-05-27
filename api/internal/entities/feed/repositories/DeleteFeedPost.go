package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	baserepo "github.com/VamaSingapore/vama-api/internal/entities/_shared/repositories"
)

func (s *repository) DeleteFeedPost(ctx context.Context, postID int, userID int) error {
	postsUpdate := `
		UPDATE
			feed.posts
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	reactionsUpdate := `
		UPDATE
			feed.post_reactions
		SET deleted_at = NOW()
		WHERE post_id = $1
	`
	imagesUpdate := `
		UPDATE
			feed.post_images
		SET deleted_at = NOW()
		WHERE post_id = $1
	`
	commentsUpdate := `
		UPDATE
			feed.post_comments
		SET deleted_at = NOW()
		WHERE post_id = $1
	`

	tx, err := s.MasterNode().Begin(ctx)
	if err != nil {
		return err
	}
	commit := false
	defer baserepo.FinishTx(ctx, tx, &commit)

	// Update feed.posts
	result, err := tx.Exec(ctx, postsUpdate, postID, userID)
	if err != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(err, postsUpdate, []interface{}{postID, userID}))
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected < 1 {
		return constants.ErrNoRowsAffected
	}

	// Update feed.post_reactions
	_, err = tx.Exec(ctx, reactionsUpdate, postID)
	if err != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(err, reactionsUpdate, []interface{}{postID}))
		return err
	}

	// Update feed.post_images
	_, err = tx.Exec(ctx, imagesUpdate, postID)
	if err != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(err, imagesUpdate, []interface{}{postID}))
		return err
	}

	// Update feed.post_comments
	_, err = tx.Exec(ctx, commentsUpdate, postID)
	if err != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(err, commentsUpdate, []interface{}{postID}))
		return err
	}

	commit = true
	return nil
}
