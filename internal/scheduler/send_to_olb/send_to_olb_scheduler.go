package send_to_olb

import (
	"test-project/internal/service/user"
	"time"

	"github.com/go-co-op/gocron"
)

type SendToOlbScheduler struct {
	cron                 string
	advisoryLockSection  int
	parallelCurrencySend int
	userService          *user.UserService
}

func New(cron string, advisoryLockSection int, parallelCurrencySend int, userService *user.UserService) *SendToOlbScheduler {
	return &SendToOlbScheduler{cron, advisoryLockSection, parallelCurrencySend, userService}
}

func (s SendToOlbScheduler) Start() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	scheduler := gocron.NewScheduler(loc)

	scheduler.Cron(s.cron).Do(executeFunc(s.advisoryLockSection, s.parallelCurrencySend, s.userService))

	scheduler.StartAsync()
}

func executeFunc(advisoryLockSection int, parallelCurrencySend int, userService *user.UserService) func() {
	return func() {
		userService.FindAndSendUsersWithStatusNewToInternalSystem(advisoryLockSection, parallelCurrencySend)
	}
}
