package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
)

const errUpvotingPost = "Something went wrong when trying to upvote post."

func (svc *usecase) UpvotePost(ctx context.Context, req request.UpvotePost) (interface{}, error) {
	user, getUserErr := svc.user.GetCurrentUser(ctx)
	if getUserErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpvotingPost,
			fmt.Sprintf("Could not find user in the current context. Err: %v", getUserErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errUpvotingPost,
			"user was nil in UpvotePost",
		)
	}

	tx, txErr := svc.repo.MasterNode().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if txErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to downvote post.",
			fmt.Sprintf("Could not create transaction: %v", txErr),
		)
	}

	commit := false
	defer svc.repo.FinishTx(ctx, tx, &commit)

	reactionStruct, getReactionErr := svc.GenerateReactionStruct(ctx, appconfig.Config.Vote.UpVote, user.ID, req.PostID, tx)
	if getReactionErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpvotingPost,
			fmt.Sprintf("Could not get reaction: %v", getReactionErr),
		)
	}

	votes, reactErr := svc.repo.React(ctx, *reactionStruct, tx)
	if reactErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errUpvotingPost,
			fmt.Sprintf("Could not save user %d reaction to database: %v", user.ID, reactErr),
		)
	}
	res := response.Reaction{
		NewState:     reactionStruct.NewState,
		NumUpvotes:   votes.UpVotes,
		NumDownvotes: votes.DownVotes,
	}
	commit = true
	return res, nil
}
