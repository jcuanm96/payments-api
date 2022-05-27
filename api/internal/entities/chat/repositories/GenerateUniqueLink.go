package repositories

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	userrrepo "github.com/VamaSingapore/vama-api/internal/entities/user/repositories"
	"github.com/google/uuid"
)

func (s *repository) GenerateUniqueLink(ctx context.Context, groupName string) (*string, error) {
	preparedName := strings.Replace(groupName, " ", "-", -1)
	preparedName = user.StripIllegalChars(preparedName, constants.LINKSUFFIX_REGEX)
	RANGES := []int{0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}
	attempts := 1
	for _, rangeSize := range RANGES {
		for j := 0; j < attempts; j++ {
			link := generateLink(preparedName, rangeSize)
			taken, isTakenErr := userrrepo.IsLinkSuffixTaken(ctx, s.MasterNode(), link)
			if isTakenErr != nil {
				return nil, isTakenErr
			}
			if !taken.IsTaken {
				return &link, nil
			}
		}
		attempts = 8
	}

	// try UUID as last resort
	link := uuid.New().String()
	taken, isTakenErr := userrrepo.IsLinkSuffixTaken(ctx, s.MasterNode(), link)
	if isTakenErr != nil {
		return nil, isTakenErr
	}
	if !taken.IsTaken {
		return &link, nil
	}

	return nil, fmt.Errorf("could not generate unique link for group name: %s", groupName)
}

func generateLink(name string, rangeSize int) string {
	randomSuffix := ""
	if rangeSize != 0 {
		randomSuffix = "-" + strconv.Itoa(rand.Intn(rangeSize)+1)
	}
	link := fmt.Sprintf("%s%s", name, randomSuffix)

	if len(link) <= constants.MAX_LINK_SUFFIX_LENGTH {
		return link
	} else {
		maxCharsLeft := constants.MAX_LINK_SUFFIX_LENGTH - len(randomSuffix)
		nameTrimmed := name[:maxCharsLeft]
		return nameTrimmed + randomSuffix
	}
}
