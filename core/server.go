package core

import (
	context2 "github.com/hamster-shared/hamster-provider/core/context"
	"github.com/hamster-shared/hamster-provider/core/corehttp"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"os"
)

type Server struct {
	ctx          context2.CoreContext
	eventService event.IEventService
}

func NewServer(ctx context2.CoreContext) *Server {
	return &Server{
		ctx:          ctx,
		eventService: ctx.EventService,
	}
}

func (s *Server) Run() {

	err := corehttp.StartApi(&s.ctx)

	if err != nil {
		os.Exit(1)
	}
}
