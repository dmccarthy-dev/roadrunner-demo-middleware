package roadrunner_demo_middleware

import (
	"context"
	"github.com/goddtriffin/helmet"
	"github.com/google/uuid"
	"github.com/roadrunner-server/errors"
	"github.com/roadrunner-server/http/v4/attributes"
	"github.com/roadrunner-server/http/v4/common"
	"go.uber.org/zap"
	"net/http"
)

// PluginName contains default service name.
const (
	PluginName     = "roadrunner_demo_middleware"
	RootPluginName = "http"
)

type Configurer interface {
	UnmarshalKey(name string, out any) error
	Has(name string) bool
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

type Plugin struct {
	log *zap.Logger
}

func (p *Plugin) Init(cfg common.Configurer, log common.Logger) error {
	const op = errors.Op("roadrunner_demo_middleware")

	if !cfg.Has(RootPluginName) {
		return errors.E(op, errors.Disabled)
	}

	p.log = new(zap.Logger)
	p.log = log.NamedLogger(PluginName)

	return nil
}

func (p *Plugin) Name() string {

	return PluginName
}

func (p *Plugin) Middleware(next http.Handler) http.Handler {

	h := helmet.Default()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r = attributes.Init(r)
		_ = attributes.Set(r, "attr_request_id", uuid.New().String())

		//next.ServeHTTP(w, r)
		h.Secure(next).ServeHTTP(w, r)
	})

}

func (p *Plugin) Stop(_ context.Context) error {

	p.log.Info("Closing open resources")

	return nil
}
