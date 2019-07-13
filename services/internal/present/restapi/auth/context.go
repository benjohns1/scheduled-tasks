package auth

import "net/http"

// ResponseContext wraps http.ResponseWriter in a context that provides authorization details to other handlers
type ResponseContext struct {
	http.ResponseWriter
	Auth Context
}

// Context contains relevant auth data from request
type Context struct {
	Issuer      string
	Subject     string
	Permissions []Permission
}

// HasPerm returns true if the request token has the specified permission
func (c *Context) HasPerm(permission Permission) bool {
	for _, check := range c.Permissions {
		if check == permission {
			return true
		}
	}
	return false
}
