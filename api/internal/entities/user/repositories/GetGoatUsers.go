package repositories

import (
	"context"
	"fmt"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"

	"github.com/Masterminds/squirrel"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (s *repository) GetGoatUsers(ctx context.Context, currUserID int, req request.GetGoatUsers) ([]response.GetGoatUsersListItem, error) {
	offset := req.PageSize * (req.Page - 1)
	queryBuilder := squirrel.Select(
		"users.id",
		"users.first_name",
		"users.last_name",
		"users.phone_number",
		"users.country_code",
		"users.email",
		"users.username",
		"users.user_type",
		"users.profile_avatar",
		"-1 AS num_goat_chats", // Deprecated field on the frontend
	).
		From("core.users users").
		Where("users.user_type = 'GOAT'")

	if req.ExcludeSelf {
		queryBuilder = queryBuilder.Where("users.id != ?", currUserID)
	}

	priority := constants.CreatorDiscoverPriority[appconfig.Config.Gcloud.Project]

	var orderByClause string

	// Assign scores to each item in the ORDER BY clause and order based on 1) the priority array 2) the number of follows a creator has.
	// Creators in the priority array will have values <= 0 and the rest of the scores will be 1/n, where n is the number of followers.
	// This system ensures that creators in the priority array come first and then the rest of the creators appear in descending order based
	// on number of followers
	//
	// For more clarity, feel free to print out the query variable below and run it in pgAdmin to see how the data is structured
	if len(priority) < 1 {
		orderByClause = `
			(SELECT COUNT(*) FROM feed.follows follows WHERE follows.goat_user_id = users.id) DESC	
		`
	} else {
		orderbyf := "CASE %s END ASC"
		buffer := ""

		for i, userID := range priority {
			whenClause := fmt.Sprintf(" WHEN users.id=%d THEN %d ", userID, -1*i)
			buffer = buffer + whenClause
		}

		elseClause := " ELSE 1.0/(SELECT NULLIF(COUNT(*),0) FROM feed.follows follows WHERE follows.goat_user_id = users.id) "
		buffer = buffer + elseClause

		orderByClause = fmt.Sprintf(orderbyf, buffer)
	}

	query, args, squirrelErr :=
		queryBuilder.
			OrderBy(orderByClause).
			Limit(uint64(req.Limit)).
			Offset(uint64(offset)).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

	if squirrelErr != nil {
		return nil, squirrelErr
	}

	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	defer rows.Close()
	users := []response.GetGoatUsersListItem{}
	var numGoatChats int
	user := &response.User{}
	for rows.Next() {
		scanErr := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Phonenumber,
			&user.CountryCode,
			&user.Email,
			&user.Username,
			&user.Type,
			&user.ProfileAvatar,
			&numGoatChats,
		)

		if scanErr != nil {
			return nil, scanErr
		}

		listItem := response.GetGoatUsersListItem{
			User:         *user,
			NumGoatChats: numGoatChats,
		}

		users = append(users, listItem)
	}
	return users, nil
}
