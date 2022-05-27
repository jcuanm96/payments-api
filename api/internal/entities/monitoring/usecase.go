package monitoring

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
)

type Usecase interface {
	PingMainDatabase(c context.Context) error
	GetMinRequiredVersionByDevice(c context.Context, req request.GetMinRequiredVersionByDevice) (*response.GetMinRequiredVersionByDevice, error)
	GetDashboard(c context.Context) (*response.GetDashboard, error)
}
