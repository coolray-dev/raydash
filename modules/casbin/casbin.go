package casbin

import (
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
)

// Enforcer is a casbin enforcer instance
var Enforcer *casbin.Enforcer

func init() {

	// Init casbin Model
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj =~ p.obj && (r.act == p.act || p.act == \"*\")")

	// Init GORM Adapter using existing DB Connection
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		log.Log.WithError(err).Fatal("Error initializing Casbin")
	}

	Enforcer, err = casbin.NewEnforcer(m, adapter)

	// Setup Enforcer
	if len(Enforcer.GetPolicy()) == 0 {
		addPolicies()
	} else {
		Enforcer.LoadPolicy()
	}
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

// AddDefaultUserPolicy add policies for a new user
func AddDefaultUserPolicy(u *models.User) {
	Enforcer.AddPolicy(u.Username, "/*/announcements*", "GET")
	Enforcer.AddPolicy(u.Username, "/*/logout", "DELETE")
	Enforcer.AddPolicy(u.Username, "/*/users/"+strconv.Itoa(int(u.ID))+"*", "*")
}
