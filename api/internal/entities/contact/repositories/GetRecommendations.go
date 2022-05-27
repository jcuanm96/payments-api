package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
)

func (s *repository) GetRecommendations(ctx context.Context, userID int) ([]response.User, error) {
	// This query selects all users who have userID as a contact but
	// who userID has not added as a contact, as well as all users
	// who have userID in their phone contacts and who userID has not added as a contact.
	query := `
	WITH my_contacts AS (
		SELECT 
		 contact_id
		FROM core.users_contacts
		WHERE user_id = $1
	   )
	   SELECT
	 	users.id,
		users.first_name,
		users.last_name,
		users.phone_number,
		users.country_code,
		users.email,
		users.username,
		users.user_type,
		users.profile_avatar
	   FROM core.users users
	   LEFT JOIN core.users_contacts contacts ON contacts.user_id = users.id
	   LEFT JOIN core.pending_contacts pending ON pending.user_id = users.id
	   WHERE (
		contacts.contact_id = $1 AND
		contacts.user_id NOT IN ( 
		 SELECT contact_id 
		 FROM my_contacts 
		)
	   ) OR (
		pending.signed_up_user_id = $1 AND
		pending.user_id NOT IN ( 
		 SELECT contact_id 
		 FROM my_contacts 
		)
	   )
	   ORDER BY random()
	   LIMIT 100
	`
	args := []interface{}{userID}
	rows, queryErr := s.MasterNode().Query(ctx, query, args...)
	if queryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(queryErr, query, args))
		return nil, queryErr
	}

	recommendedContacts := []response.User{}
	addedIDsSet := map[int]struct{}{}
	defer rows.Close()
	for rows.Next() {
		var user response.User
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
		)
		if scanErr != nil {
			return nil, scanErr
		}
		if _, ok := addedIDsSet[user.ID]; ok {
			continue
		}
		recommendedContacts = append(recommendedContacts, user)
		addedIDsSet[user.ID] = struct{}{}
	}

	return recommendedContacts, nil
}
