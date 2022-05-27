package repositories

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/jackc/pgx/v4"
)

func (s *repository) GetDashboard(ctx context.Context) (*response.GetDashboard, error) {
	usersQuery := `
	SELECT 
    	COUNT(id) AS count_total, 
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 day' THEN 1 ELSE 0 END) OVER() count_day,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 week' THEN 1 ELSE 0 END) OVER() count_week,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 month' THEN 1 ELSE 0 END) OVER() count_month,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 year' THEN 1 ELSE 0 END) OVER() count_year
	FROM 
    	core.users
	GROUP BY 
		users.created_at
	`
	postsQuery := `
	SELECT
		COUNT(id) AS count_total, 
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 day' THEN 1 ELSE 0 END) OVER() count_day,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 week' THEN 1 ELSE 0 END) OVER() count_week,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 month' THEN 1 ELSE 0 END) OVER() count_month,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 year' THEN 1 ELSE 0 END) OVER() count_year
	FROM
		feed.posts
	GROUP BY 
		posts.created_at
	`

	chatsQuery := `
	SELECT
		COUNT(id) AS count_total, 
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 day' THEN 1 ELSE 0 END) OVER() count_day,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 week' THEN 1 ELSE 0 END) OVER() count_week,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 month' THEN 1 ELSE 0 END) OVER() count_month,
		SUM(CASE WHEN created_at > NOW() - INTERVAL '1 year' THEN 1 ELSE 0 END) OVER() count_year
	FROM
		feed.goat_chat_messages
	GROUP BY 
		goat_chat_messages.created_at
	`

	dashboard := response.GetDashboard{
		Total: response.DashboardStats{},
		Day:   response.DashboardStats{},
		Week:  response.DashboardStats{},
		Month: response.DashboardStats{},
		Year:  response.DashboardStats{},
	}

	tx, txErr := s.MasterNode().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadUncommitted})
	if txErr != nil {
		return nil, txErr
	}

	defer tx.Commit(ctx)

	userRows, userQueryErr := tx.Query(ctx, usersQuery)
	if userQueryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(userQueryErr, usersQuery, nil))
		return nil, userQueryErr
	}

	defer userRows.Close()
	if !userRows.Next() {
		return nil, constants.ErrNotFound
	}

	userScanErr := userRows.Scan(
		&dashboard.Day.NewUsers,
		&dashboard.Week.NewUsers,
		&dashboard.Month.NewUsers,
		&dashboard.Year.NewUsers,
		&dashboard.Total.NewUsers,
	)
	if userScanErr != nil {
		return nil, userScanErr
	}

	userRows.Close()

	postRows, postQueryErr := tx.Query(ctx, postsQuery)
	if postQueryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(postQueryErr, postsQuery, nil))
		return nil, postQueryErr
	}

	defer postRows.Close()
	if !postRows.Next() {
		return nil, constants.ErrNotFound
	}

	postScanErr := postRows.Scan(
		&dashboard.Day.FeedPosts,
		&dashboard.Week.FeedPosts,
		&dashboard.Month.FeedPosts,
		&dashboard.Year.FeedPosts,
		&dashboard.Total.FeedPosts,
	)

	if postScanErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(postScanErr, postsQuery, nil))
		return nil, postScanErr
	}

	postRows.Close()

	chatRows, chatQueryErr := tx.Query(ctx, chatsQuery)
	if chatQueryErr != nil {
		vlog.Error(ctx, utils.SqlErrLogMsg(chatQueryErr, chatsQuery, nil))
		return nil, chatQueryErr
	}

	defer chatRows.Close()
	if !chatRows.Next() {
		return nil, constants.ErrNotFound
	}

	chatScanErr := chatRows.Scan(
		&dashboard.Day.GoatChats,
		&dashboard.Week.GoatChats,
		&dashboard.Month.GoatChats,
		&dashboard.Year.GoatChats,
		&dashboard.Total.GoatChats,
	)

	if chatScanErr != nil {
		return nil, chatScanErr
	}
	chatRows.Close()

	return &dashboard, nil
}
