package casbin

import (
	"strconv"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
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
	addPolicies()
}

func addPolicies() {
	basicRules := [][]string{
		[]string{"group::admin", "/*", "*"},
		[]string{"role::anonymous", "/*/swagger/*", "*"},
		[]string{"role::anonymous", "/*/login", "POST"},
		[]string{"role::anonymous", "/*/register", "POST"},
		[]string{"role::anonymous", "/*/refresh", "POST"},
		[]string{"role::anonymous", "/*/password/*", "POST"},
	}
	Enforcer.AddPolicies(basicRules)
	var groups []models.Group
	if err := database.DB.Preload("Users").Find(&groups).Error; err != nil {
		log.Log.WithError(err).Error()
	}
	for _, g := range groups {
		Enforcer.AddPolicy("group::"+g.Name, "/*/groups/"+strconv.Itoa(int(g.ID))+"*", "*")
		for _, u := range g.Users {
			Enforcer.AddGroupingPolicy(u.Username, "group::"+g.Name)
			Enforcer.AddPolicy(u.Username, "/*/announcements*", "GET")
			Enforcer.AddPolicy(u.Username, "/*/logout", "DELETE")
		}
	}
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Log.WithError(err).Error()
	}
	for _, u := range users {
		Enforcer.AddPolicy(u.Username, "/*/users/"+strconv.Itoa(int(u.ID))+"*", "*")
	}
}

func AddDefaultUserPolicy(u *models.User) {
	Enforcer.AddPolicy(u.Username, "/*/announcements*", "GET")
	Enforcer.AddPolicy(u.Username, "/*/logout", "DELETE")
	Enforcer.AddPolicy(u.Username, "/*/users/"+strconv.Itoa(int(u.ID))+"*", "*")
}
