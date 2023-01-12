package policy

import (
	"diamond/misc"
	"diamond/models"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func init() {
	// 具有超级用户的ACL
	text :=
		`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act || r.sub == "admin"
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		misc.Logger.Fatal().Err(err).Str("from", "policy").Msg("load policy model failed")
	}

	a, err := gormadapter.NewAdapterByDB(models.DB)
	if err != nil {
		misc.Logger.Fatal().Err(err).Str("from", "policy").Msg("load policy adapter failed")
	}

	// 创建执行者
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		misc.Logger.Fatal().Err(err).Str("from", "policy").Msg("init policy enforcer failed")
	}
	// add match func
	// e.AddNamedMatchingFunc("g2", "KeyMatch2", util.KeyMatch2)
	// Load the policy from DB.
	e.LoadPolicy()

	Enforcer = e
}

/* policy
// route policy
p, alice, /api/users, GET
p, bob, /api/users/:id, PUT
*/
