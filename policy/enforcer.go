package policy

import (
	"diamond/models"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func init() {
	// 支持资源角色的RBAC
	text :=
		`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _
		g2 = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		log.Fatalln("load policy model failed: ", err)
	}

	a, err := gormadapter.NewAdapterByDB(models.DB)
	if err != nil {
		log.Fatalln("load policy adapter failed: ", err)
	}

	// 创建执行者
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalln("init policy enforcer failed: ", err)
	}
	// add match func
	e.AddNamedMatchingFunc("g", "KeyMatch2", util.KeyMatch2)
	// Load the policy from DB.
	e.LoadPolicy()

	Enforcer = e
}

/* policy
// route policy
p, role::1, GET /api/users, route
p, role::2, POST /api/users, route

// menu policy
p, role::1, system, menu
p, role::2, system-user, menu

// server policy
p, role::1, group::1, server
p, role::2, group::2, server

// user role assign
g, user::1, role::1

// server group assign
g2, server::1, group::1
g2, server::2, group::2
*/
