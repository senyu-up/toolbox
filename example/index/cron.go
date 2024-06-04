package index

import (
	"github.com/senyu-up/toolbox/example/internal/cron"
	"github.com/senyu-up/toolbox/tool/cronv2"
)

func CronRegister(c *cronv2.Client) {
	c.Register("*/30 * * * *", "RefreshOrganize", cron.DailyCountJob)
}
