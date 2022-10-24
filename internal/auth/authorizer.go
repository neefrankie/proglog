package auth

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authorizer wraps Casbin authorization.
type Authorizer struct {
	enforcer *casbin.Enforcer
}

// New creates a new instance of Authorizer.
// The model and policy arguments are paths to the files where you defined the model.
func New(model, policy string) *Authorizer {
	enforcer, err := casbin.NewEnforcer(model, policy, true)
	if err != nil {
		panic(err)
	}
	return &Authorizer{
		enforcer: enforcer,
	}
}

// Authorize defers to Casbin's Enforce function to return whether the given
// subject is permitted ron the given action on the given object based on the
// model and policy you configure Casbin with.
func (a *Authorizer) Authorize(subject, object, action string) error {
	ok, err := a.enforcer.Enforce(subject, object, action)
	if err != nil {
		return err
	}
	if !ok {
		msg := fmt.Sprintf("%s not permitted to %s to %s", subject, action, object)
		st := status.New(codes.PermissionDenied, msg)
		return st.Err()
	}
	return nil
}
