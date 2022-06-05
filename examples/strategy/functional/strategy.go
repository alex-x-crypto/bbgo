package functional

import (
	"github.com/sirupsen/logrus"

	"github.com/c9s/bbgo/pkg/bbgo"
	"github.com/c9s/bbgo/pkg/types"
)

var log = logrus.WithField("strategy", ID)

const ID = "functional"

type Config struct {
	Symbol   string         `json:"symbol"`
	Interval types.Interval `json:"interval"`
}

func init() {
	bbgo.RegisterStrategy(ID, func(config Config, session *bbgo.ExchangeSession) func(orderExecutor bbgo.OrderExecutor) error {
		session.Subscribe(types.KLineChannel, "BTCUSDT", types.SubscribeOptions{Interval: types.Interval1h})

		// return Run() function
		return func(orderExecutor bbgo.OrderExecutor) error {
			return nil
		}
	})
}
