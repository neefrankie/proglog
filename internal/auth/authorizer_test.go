package auth

import (
	"com.github/neefrankie/proglog/internal/config"
	"testing"
)

func TestAuthorizer_Authorize(t *testing.T) {
	authorizer := New(config.ACLModelFile, config.ACLPolicyFile)

	err := authorizer.Authorize("root", "*", "produce")

	if err != nil {
		t.Error(err)
	}
}
