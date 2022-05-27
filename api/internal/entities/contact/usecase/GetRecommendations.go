package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const contactRecsErr = "Something went wrong when getting contact recommendations."

func (svc *usecase) GetRecommendations(ctx context.Context, limit int) (*response.ContactRecommendations, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			contactRecsErr,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			contactRecsErr,
			"user came back nil in GetRecommendations",
		)
	}

	recommendations, getRecommendationsErr := svc.repo.GetRecommendations(ctx, user.ID)
	if getRecommendationsErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			contactRecsErr,
			fmt.Sprintf("Error getting contact recommendations for %d: %v", user.ID, getRecommendationsErr),
		)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(recommendations), func(i, j int) {
		recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
	})

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	res := &response.ContactRecommendations{
		Recommendations: recommendations,
	}
	return res, nil
}
