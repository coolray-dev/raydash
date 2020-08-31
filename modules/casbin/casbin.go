package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/coolray-dev/raydash/modules/utils"
)

var Enforcer *casbin.Enforcer

func init() {
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		log.Log.WithError(err).Fatal("Error initializing Casbin")
	}
	Enforcer, err = casbin.NewEnforcer(utils.AbsPath("modules/casbin/rbac.conf"), adapter)
}
