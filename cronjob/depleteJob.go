package cronjob

import (
	"github.com/wanan9999/s-ui/database"
	"github.com/wanan9999/s-ui/logger"
	"github.com/wanan9999/s-ui/service"
)

type DepleteJob struct {
	service.ClientService
	service.InboundService
}

func NewDepleteJob() *DepleteJob {
	return new(DepleteJob)
}

func (s *DepleteJob) Run() {
	inboundIds, err := s.ClientService.DepleteClients()
	if err != nil {
		logger.Warning("Disable depleted users failed: ", err)
		return
	}
	if len(inboundIds) > 0 {
		err := s.InboundService.UpdateInboundsUsers(database.GetDB(), inboundIds)
		if err != nil {
			logger.Error("unable to update inbound users: ", err)
		}
	}
}
