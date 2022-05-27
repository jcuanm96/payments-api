package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	sendbird "github.com/VamaSingapore/vama-api/internal/sendbirdClient"
)

func (svc *usecase) GetMessages(ctx context.Context, query string, userID int, limit int) ([]response.SearchMessage, error) {
	searchMessagesParams := sendbird.SearchMessagesParams{
		UserID: fmt.Sprint(userID),
		Query:  query,
		Limit:  limit,
	}

	searchMessagesRes, searchMessagesErr := svc.sendbirdClient.SearchMessages(searchMessagesParams)

	if searchMessagesErr != nil {
		return nil, searchMessagesErr
	} else if searchMessagesRes == nil {
		return []response.SearchMessage{}, nil
	}

	var wg sync.WaitGroup

	results := []response.SearchMessage{}
	queryLen := len(query)
	wg.Add(searchMessagesRes.TotalCount)
	for _, message := range searchMessagesRes.Results {
		if message.Message == nil {
			continue
		}
		go func(m sendbird.Message) {
			defer wg.Done()
			startRange := strings.Index(strings.ToLower(*m.Message), query)
			endRange := startRange + queryLen - 1

			resultMessage := response.SearchMessage{
				StartRange: startRange,
				EndRange:   endRange,
				Message:    m,
			}

			results = append(results, resultMessage)
		}(message)
	}

	wg.Wait()

	return results, nil
}
