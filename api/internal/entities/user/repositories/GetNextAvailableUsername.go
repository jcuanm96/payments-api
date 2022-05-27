package repositories

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/google/uuid"
)

func (s *repository) GetNextAvailableUsername(ctx context.Context, firstName string, lastName string) (*string, error) {
	RANGES := []int{0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}
	attempts := 1
	for _, rangeSize := range RANGES {
		for j := 0; j < attempts; j++ {
			userName := generateUsername(firstName, lastName, rangeSize)
			taken, isTakenErr := s.IsLinkSuffixTaken(ctx, s.MasterNode(), userName)
			if isTakenErr != nil {
				return nil, isTakenErr
			}
			if !taken.IsTaken {
				return &userName, nil
			}
		}
		attempts = 8
	}

	// try UUID as last resort
	link := uuid.New().String()
	taken, isTakenErr := s.IsLinkSuffixTaken(ctx, s.MasterNode(), link)
	if isTakenErr != nil {
		return nil, isTakenErr
	}
	if !taken.IsTaken {
		return &link, nil
	}

	return nil, fmt.Errorf("could not generate unique username for name: %s %s", firstName, lastName)
}

var emptyFirstNameWords = []string{"Amazing", "Bright", "Concentrated", "Energetic", "Enchanted",
	"Dynamic", "Fearless", "Funny", "Focused", "Generous", "Helpful", "Inquisitive", "Knowledgeable", "Magnificent",
	"Omniscient", "Optimal", "Resourceful", "Shining", "Tactical", "Quick"}

var emptyLastNameWords = []string{"Ant", "Albatross", "Beagle", "Bat", "Bison", "Camel", "Capybara", "Centaur", "Crab",
	"Doge", "Dove", "Dragonfly", "Eagle", "Falcon", "Fox", "Giraffe", "Hyena", "Jellyfish", "Jaguar", "Kraken", "Lion",
	"Otter", "Owl", "Panther", "Parrot", "Pegasus", "Penguin", "Rhinoceros", "Seal", "Snake", "Squid", "Turtle", "Unicorn",
	"Zebra"}

func generateUsername(first string, last string, rangeSize int) string {
	first = user.StripIllegalChars(first, constants.LINKSUFFIX_REGEX)
	last = user.StripIllegalChars(last, constants.LINKSUFFIX_REGEX)

	if first == "" {
		// Choose random element in slice
		first = emptyFirstNameWords[rand.Intn(len(emptyFirstNameWords))]
	}

	if last == "" {
		last = emptyLastNameWords[rand.Intn(len(emptyLastNameWords))]
	}

	randomUserNameSuffix := ""
	if rangeSize != 0 {
		randomUserNameSuffix = "-" + strconv.Itoa(rand.Intn(rangeSize)+1)
	}
	username := fmt.Sprintf("%s-%s%s", first, last, randomUserNameSuffix)

	if len(username) <= constants.MAX_LINK_SUFFIX_LENGTH {
		return username
	} else {
		suffix := "-" + fmt.Sprint(randomUserNameSuffix)
		maxCharsLeft := constants.MAX_LINK_SUFFIX_LENGTH - len(suffix)
		lastTrimmed := last[:min(len(last), maxCharsLeft)]
		return lastTrimmed + suffix
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
