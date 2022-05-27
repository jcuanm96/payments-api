package service

import (
	"context"
	"errors"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/feed"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (svc *usecase) GenerateReactionStruct(ctx context.Context, reactionType string, userID int, postID int, runnable utils.Runnable) (*feed.Reaction, error) {
	upvotesColumn := "num_upvotes"
	downvotesColumn := "num_downvotes"

	previousState, getPreviousStateErr := svc.repo.GetReaction(ctx, userID, postID, runnable)
	if getPreviousStateErr != nil {
		vlog.Errorf(ctx, "Error getting previous reaction state: %v", getPreviousStateErr)
		return nil, getPreviousStateErr
	} else if previousState == nil {
		getPreviousStateErr = errors.New("previous state is nil while previous state error is also nil")
		vlog.Errorf(ctx, getPreviousStateErr.Error())
		return nil, getPreviousStateErr
	}

	reactionStruct := feed.Reaction{
		PreviousState:   *previousState,
		ReactionType:    reactionType,
		UserID:          userID,
		PostID:          postID,
		DecrementColumn: nil,
		IncrementColumn: nil,
		NewState:        appconfig.Config.Vote.Nil,
	}

	if *previousState == appconfig.Config.Vote.UpVote {
		reactionStruct.DecrementColumn = &upvotesColumn
	} else if *previousState == appconfig.Config.Vote.DownVote {
		reactionStruct.DecrementColumn = &downvotesColumn
	}

	// Decide whether or not a column needs to be incremented
	if *previousState != reactionType {
		reactionStruct.IncrementColumn = &downvotesColumn
		if reactionType == appconfig.Config.Vote.UpVote {
			reactionStruct.IncrementColumn = &upvotesColumn
		}
	}

	if reactionType == appconfig.Config.Vote.UpVote && *previousState != appconfig.Config.Vote.UpVote {
		reactionStruct.NewState = appconfig.Config.Vote.UpVote
	} else if reactionType == appconfig.Config.Vote.DownVote && *previousState != appconfig.Config.Vote.DownVote {
		reactionStruct.NewState = appconfig.Config.Vote.DownVote
	}
	return &reactionStruct, nil
}
