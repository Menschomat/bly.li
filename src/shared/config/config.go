package config

import (
	"sync"

	m "github.com/Menschomat/bly.li/shared/model"
	cfgUtils "github.com/Menschomat/bly.li/shared/utils/config"
)

// ConfigManager handles thread-safe initialization for any config type
type ConfigManager[T any] struct {
	once   sync.Once
	config *T
	err    error
}

func NewConfigManager[T any]() *ConfigManager[T] {
	return &ConfigManager[T]{}
}

func (cm *ConfigManager[T]) Get() (*T, error) {
	cm.once.Do(func() {
		cm.config = new(T)
		if err := cfgUtils.FillEnvStruct(cm.config); err != nil {
			cm.err = err
		}
	})
	return cm.config, cm.err
}

func (cm *ConfigManager[T]) MustGet() *T {
	cfg, err := cm.Get()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Configuration accessors
func ShortnConfig() *m.ShortnConfig {
	return managers.Shortn.MustGet()
}

func OidcConfig() *m.OidcConfig {
	return managers.Oidc.MustGet()
}

func MongoConfig() *m.MongoDdConfig {
	return managers.Mongo.MustGet()
}

func LoggingConfig() *m.LoggingConfig {
	return managers.Logging.MustGet()
}

func BlowupConfig() *m.BlowupConfig {
	return managers.Blowup.MustGet()
}

func DasherConfig() *m.DasherConfig {
	return managers.Dasher.MustGet()
}

func PersoConfig() *m.PersoConfig {
	return managers.Perso.MustGet()
}

// Manager registry
var managers = struct {
	Shortn  *ConfigManager[m.ShortnConfig]
	Oidc    *ConfigManager[m.OidcConfig]
	Mongo   *ConfigManager[m.MongoDdConfig]
	Logging *ConfigManager[m.LoggingConfig]
	Blowup  *ConfigManager[m.BlowupConfig]
	Dasher  *ConfigManager[m.DasherConfig]
	Perso   *ConfigManager[m.PersoConfig]
}{
	Shortn:  NewConfigManager[m.ShortnConfig](),
	Oidc:    NewConfigManager[m.OidcConfig](),
	Mongo:   NewConfigManager[m.MongoDdConfig](),
	Logging: NewConfigManager[m.LoggingConfig](),
	Blowup:  NewConfigManager[m.BlowupConfig](),
	Dasher:  NewConfigManager[m.DasherConfig](),
	Perso:   NewConfigManager[m.PersoConfig](),
}
