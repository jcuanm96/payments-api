package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const (
	errSearchMention = "Something went wrong getting users to mention."
	limit            = 100
)

func (svc *usecase) SearchMention(ctx context.Context, req request.SearchMention) (*response.SearchMention, error) {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			errSomethingWentWrongSearch,
			fmt.Sprintf("Could not find user in the current context: %v", userErr),
		)
	}
	if user == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			errSomethingWentWrongSearch,
			"Current user was nil in SearchMention",
		)
	}
	currentUserIDstr := fmt.Sprint(user.ID)
	if req.Query == "" {
		listMembersParams := &sendbird.ListGroupChannelMembersParams{
			Limit: limit,
		}
		memberList, listMembersErr := svc.sendbirdClient.ListGroupChannelMembers(req.ChannelID, *listMembersParams)
		if listMembersErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errSearchMention,
				fmt.Sprintf("Error getting empty search results: %v", listMembersErr),
			)
		}

		// don't include the requesting user in the mention results
		filteredMembers := make([]sendbird.User, 0, len(memberList.Members))
		for _, member := range memberList.Members {
			if member.UserID != currentUserIDstr {
				filteredMembers = append(filteredMembers, member)
			}
		}

		res := &response.SearchMention{
			Users: filteredMembers,
		}
		return res, nil
	}

	type getMemberBatch struct {
		batch []sendbird.User
		err   error
	}
	// Buffer this channel so we can get more members while we process.
	memberBatches := make(chan getMemberBatch, 20)
	go func() {
		defer close(memberBatches)
		token := ""
		for {
			listMembersParams := &sendbird.ListGroupChannelMembersParams{
				Token: token,
				Limit: limit, // max limit is 100: https://sendbird.com/docs/chat/v3/platform-api/guides/group-channel#2-list-members

			}
			memberList, listMembersErr := svc.sendbirdClient.ListGroupChannelMembers(req.ChannelID, *listMembersParams)

			memberBatches <- getMemberBatch{memberList.Members, listMembersErr}
			if listMembersErr != nil {
				return
			}

			if memberList.Next == "" {
				return
			}
			token = memberList.Next
		}
	}()

	prefixMatches := []sendbird.User{}
	substringMatches := []sendbird.User{}
	for memberBatch := range memberBatches {
		if memberBatch.err != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				errSearchMention,
				fmt.Sprintf("Error getting member batch from sendbird: %v", memberBatch.err),
			)
		}

		for _, member := range memberBatch.batch {
			// don't include the requesting user in the mention results
			if member.UserID == currentUserIDstr {
				continue
			}
			if isMatchErr := checkMentionMatch(&member, req.Query, &prefixMatches, &substringMatches); isMatchErr != nil {
				return nil, httperr.NewCtx(
					ctx,
					500,
					http.StatusInternalServerError,
					errSearchMention,
					fmt.Sprintf("Error checking search match: %v", isMatchErr),
				)
			}
		}
	}

	res := &response.SearchMention{
		Users: append(prefixMatches, substringMatches...),
	}
	return res, nil
}

func checkMentionMatch(
	member *sendbird.User,
	query string,
	prefixMatches *[]sendbird.User,
	substringMatches *[]sendbird.User,
) error {
	metadataBytes, marshalErr := json.Marshal(member.Metadata)
	if marshalErr != nil {
		return marshalErr
	}
	metadata := &response.SendbirdUserMetadata{}
	unmarshalErr := json.Unmarshal(metadataBytes, metadata)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	username := strings.ToLower(metadata.Username)
	nickname := strings.ToLower(member.NickName)
	trimmedQuery := strings.TrimSpace(query)

	isPrefixMatch := strings.HasPrefix(username, trimmedQuery) ||
		strings.HasPrefix(nickname, trimmedQuery) ||
		strings.HasPrefix(nickname, query)

	if isPrefixMatch {
		*prefixMatches = append(*prefixMatches, *member)
		return nil
	}

	isSubstringMatch := strings.Contains(username, trimmedQuery) ||
		strings.Contains(nickname, trimmedQuery) ||
		strings.Contains(nickname, query)

	if isSubstringMatch {
		*substringMatches = append(*substringMatches, *member)
		return nil
	}
	return nil
}
