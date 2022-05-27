package feed

import (
	"fmt"
)

type Reaction struct {
	NewState        string
	ReactionType    string
	IncrementColumn *string
	DecrementColumn *string
	PreviousState   string
	UserID          int
	PostID          int
}

type ReactionMetadata struct {
	UpVotes   int64
	DownVotes int64
}

type GetFeedPostsParams struct {
	IsTextContentFullLength bool
	// This only non-nil for GetFeedPostByLinkSuffix
	LinkSuffix *string

	// This only non-nil for GetFeedPostByID
	PostID *int

	// This only non-nil for GetGoatFeedPosts
	GoatUserID *int

	// These non-nil for GetGoatFeedPosts and GetUserFeedPosts
	CursorPostID *int
	Limit        *int
}

func (r Reaction) GetUpdateColumnsValues() (string, string) {
	// Declare columns and values
	var column1, column2 string
	var value1, value2 string

	// If there is a column to decrement, assign that column, and it's corresponding decrement statement, to column1, value1
	if r.DecrementColumn != nil {
		column1 = *r.DecrementColumn
		value1 = *r.DecrementColumn + "-1"
	}

	// If there is a column to Increment, assign that column, and it's corresponding decrement statement, to column1, value1
	if r.IncrementColumn != nil {
		column2 = *r.IncrementColumn
		value2 = *r.IncrementColumn + "+1"
		// If there is both a column to increment and decrement, we can safely return a tuple formatted for a SET () = () statement
		if r.DecrementColumn != nil {
			return fmt.Sprintf("(%s, %s)", column1, column2), fmt.Sprintf("(%s, %s)", value1, value2)
		}
		// If decrement column is nil, we know that it is the increment columns/values we want, so return col2/val2
		return column2, value2
	}
	// We didn't return in the increment block, meaning we can safely return the decrement col/val
	return column1, value1
}

func (r Reaction) GetInsertParams() []interface{} {
	params := make([]interface{}, 0)
	params = append(
		params,
		r.PostID,
		r.UserID,
		r.NewState,
		r.NewState,
	)
	return params
}
