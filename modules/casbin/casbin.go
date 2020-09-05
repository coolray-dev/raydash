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
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj =~ p.obj && regexMatch(r.act, p.act) || g(r.sub, \"group::admin\")")

	// Init GORM Adapter using existing DB Connection
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		log.Log.WithError(err).Fatal("Error initializing Casbin")
	}

	// Auto load from DB
	Enforcer, err = casbin.NewEnforcer(m, adapter)

	Enforcer.ClearPolicy()
	// Setup Enforcer
	if len(Enforcer.GetPolicy()) == 0 {
		addPolicies()
	}
}

func addPolicies() {
	basicRules := [][]string{
		{"group::admin", "/*", ".*"},
		{"role::anonymous", "/*/swagger/.*", ".*"},
		{"role::anonymous", "/*/login", "POST"},
		{"role::anonymous", "/*/register", "POST"},
		{"role::anonymous", "/*/refresh", "POST"},
		{"role::anonymous", "/*/password/.*", "POST"},
	}
	Enforcer.AddPolicies(basicRules)

	// Add group policies
	var groups []models.Group
	if err := database.DB.Preload("Users").Find(&groups).Error; err != nil {
		log.Log.WithError(err).Error()
	}
	for _, g := range groups {
		Enforcer.AddPolicy("group::"+g.Name, "/*/groups/"+strconv.Itoa(int(g.ID))+".*", "*")
		for _, u := range g.Users {
			// Add user to group
			Enforcer.AddGroupingPolicy(u.Username, "group::"+g.Name)
		}
	}

	// Add user policies
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Log.WithError(err).Error()
	}
	for _, u := range users {
		AddDefaultUserPolicy(&u)
	}

	// Add node policies
	var nodes []models.Node
	if err := database.DB.Find(&nodes).Error; err != nil {
		log.Log.WithError(err).Error()
	}
	for _, n := range nodes {
		Enforcer.AddPolicy("node::"+strconv.Itoa(int(n.ID)),
			"/*/nodes/"+strconv.Itoa(int(n.ID))+".*",
			".*")
	}

	// Explicitly trigger save policies
	Enforcer.SavePolicy()
}

// AddDefaultUserPolicy add policies for a user
func AddDefaultUserPolicy(u *models.User) {
	Enforcer.AddPolicy(u.Username, "/*/announcements.*", "GET")
	Enforcer.AddPolicy(u.Username, "/*/logout", "DELETE")
	Enforcer.AddPolicy(u.Username, "/*/users/"+u.Username+"$", ".*")
	Enforcer.AddPolicy(u.Username, "/*/users/"+u.Username+"/(groups|services|nodes)$", "GET")
	Enforcer.AddPolicy(u.Username, "/*/nodes$", "GET")

	// Add policy for owned services
	var services []models.Service
	if err := database.DB.
		Where("uid = ?", u.ID).Find(&services).Error; err != nil {
		log.Log.WithError(err).Error()
	}

	for _, s := range services {
		Enforcer.AddPolicy(u.Username,
			"/*/users/"+u.Username+"/services/"+strconv.Itoa(int(s.ID)),
			".*")
	}
}
