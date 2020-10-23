package config

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/c9s/bbgo/pkg/bbgo"
)

type SingleExchangeStrategyConfig struct {
	Mounts   []string
	Strategy bbgo.SingleExchangeStrategy
}

type StringSlice []string

func (s *StringSlice) decode(a interface{}) error {
	switch d := a.(type) {
	case string:
		*s = append(*s, d)

	case []string:
		*s = append(*s, d...)

	case []interface{}:
		for _, de := range d {
			if err := s.decode(de); err != nil {
				return err
			}
		}

	default:
		return errors.Errorf("unexpected type %T for StringSlice: %+v", d, d)
	}

	return nil
}

func (s *StringSlice) UnmarshalJSON(b []byte) error {
	var a interface{}
	var err = json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return s.decode(a)
}

type PnLReporter struct {
	AverageCostBySymbols StringSlice `json:"averageCostBySymbols"`
	Of                   StringSlice `json:"of" yaml:"of"`
	When                 StringSlice `json:"when" yaml:"when"`
}

type Config struct {
	Imports []string `json:"imports" yaml:"imports"`

	ExchangeStrategies      []SingleExchangeStrategyConfig
	CrossExchangeStrategies []bbgo.CrossExchangeStrategy

	PnLReporters []PnLReporter `json:"reportPnL" yaml:"reportPnL"`
}

type Stash map[string]interface{}

func loadStash(configFile string) (Stash, error) {
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	stash := make(Stash)
	if err := yaml.Unmarshal(config, stash); err != nil {
		return nil, err
	}

	return stash, err
}

func Load(configFile string) (*Config, error) {
	var config Config

	stash, err := loadStash(configFile)
	if err != nil {
		return nil, err
	}

	if err := loadExchangeStrategies(&config, stash); err != nil {
		return nil, err
	}

	if err := loadCrossExchangeStrategies(&config, stash); err != nil {
		return nil, err
	}

	if err := loadReportPnL(&config, stash); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadReportPnL(config *Config, stash Stash) error {
	reporterStash, ok := stash["reportPnL"]
	if !ok {
		return nil
	}

	reporters, err := reUnmarshal(reporterStash, &config.PnLReporters)
	if err != nil {
		return err
	}

	config.PnLReporters = *(reporters.(*[]PnLReporter))
	return nil
}

func loadCrossExchangeStrategies(config *Config, stash Stash) (err error) {
	if len(bbgo.LoadedCrossExchangeStrategies) == 0 {
		return errors.New("no cross exchange strategy is registered")
	}


	exchangeStrategiesConf, ok := stash["crossExchangeStrategies"]
	if !ok {
		return nil
	}

	configList, ok := exchangeStrategiesConf.([]interface{})
	if !ok {
		return errors.New("expecting list in crossExchangeStrategies")
	}

	for _, entry := range configList {
		configStash, ok := entry.(Stash)
		if !ok {
			return errors.Errorf("strategy config should be a map, given: %T %+v", entry, entry)
		}

		for id, conf := range configStash {
			// look up the real struct type
			if st, ok := bbgo.LoadedExchangeStrategies[id]; ok {
				val, err := reUnmarshal(conf, st)
				if err != nil {
					return err
				}

				config.CrossExchangeStrategies = append(config.CrossExchangeStrategies, val.(bbgo.CrossExchangeStrategy))
			}
		}
	}

	return nil
}

func loadExchangeStrategies(config *Config, stash Stash) (err error) {
	if len(bbgo.LoadedExchangeStrategies) == 0 {
		return errors.New("no exchange strategy is registered")
	}

	exchangeStrategiesConf, ok := stash["exchangeStrategies"]
	if !ok {
		return nil
	}

	configList, ok := exchangeStrategiesConf.([]interface{})
	if !ok {
		return errors.New("expecting list in exchangeStrategies")
	}


	for _, entry := range configList {
		configStash, ok := entry.(Stash)
		if !ok {
			return errors.Errorf("strategy config should be a map, given: %T %+v", entry, entry)
		}

		var mounts []string
		if val, ok := configStash["on"]; ok {
			if values, ok := val.([]string); ok {
				mounts = append(mounts, values...)
			} else if str, ok := val.(string); ok {
				mounts = append(mounts, str)
			}
		}


		for id, conf := range configStash {
			// look up the real struct type
			if st, ok := bbgo.LoadedExchangeStrategies[id]; ok {
				val, err := reUnmarshal(conf, st)
				if err != nil {
					return err
				}

				config.ExchangeStrategies = append(config.ExchangeStrategies, SingleExchangeStrategyConfig{
					Mounts:   mounts,
					Strategy: val.(bbgo.SingleExchangeStrategy),
				})
			}
		}
	}

	return nil
}

func reUnmarshal(conf interface{}, tpe interface{}) (interface{}, error) {
	// get the type "*Strategy"
	rt := reflect.TypeOf(tpe)

	// allocate new object from the given type
	val := reflect.New(rt)

	// now we have &(*Strategy) -> **Strategy
	valRef := val.Interface()

	plain, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(plain, valRef); err != nil {
		return nil, errors.Wrapf(err, "json parsing error, given payload: %s", plain)
	}

	return val.Elem().Interface(), nil
}
