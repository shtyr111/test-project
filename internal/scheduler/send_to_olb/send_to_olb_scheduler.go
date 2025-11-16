package send_to_olb

import (
	"test-project/internal/service"
	"time"

	"github.com/go-co-op/gocron"
)

type SendToOlbScheduler struct {
	cron                string
	advisoryLockSection int
	userService         *service.UserService
}

func New(cron string, advisoryLockSection int, userService *service.UserService) *SendToOlbScheduler {
	return &SendToOlbScheduler{cron, advisoryLockSection, userService}
}

func (s SendToOlbScheduler) Start() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	scheduler := gocron.NewScheduler(loc)

	scheduler.Cron(s.cron).Do(executeFunc(s.advisoryLockSection, s.userService))

	scheduler.StartAsync()
}

func executeFunc(advisoryLockSection int, userService *service.UserService) func() {
	return func() {
		userService.SendUsersWithStatusNewToInternalSystem(advisoryLockSection)
	}
}
