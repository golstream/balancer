package configuration

import (
	"balancer/internal/constants"
	"github.com/caarlos0/env/v11"
	"reflect"
	"strings"
)

type config struct {
	Host string `env:"HOST" envDefault:"127.0.0.1"`
	Port int64  `env:"PORT,required"`

	Method  method  `env:"METHOD,required"`
	Servers servers `env:"SERVERS,required"`
	Weights weights `env:"WEIGHTS, required"`

	HealthCheckInterval int `env:"HEALTHCHECK_INTERVAL" envDefault:"60"`
	HealthCheckTimeout  int `env:"HEALTHCHECK_TIMEOUT"  envDafault:"15"`

	WithLog bool `env:"WITH_LOG"`
}

type (
	servers []string
	weights []string
	method  = constants.BalanceMethod
)

func Init() (cfg config, err error) {
	err = env.ParseWithOptions(
		&cfg,
		env.Options{FuncMap: options},
	)

	return
}

var (
	options = map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(servers(nil)): func(v string) (interface{}, error) {
			return strings.Split(v, ","), nil
		},

		reflect.TypeOf(weights(nil)): func(v string) (interface{}, error) {
			return strings.Split(v, ","), nil
		},
	}
)
