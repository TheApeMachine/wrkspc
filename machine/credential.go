package machine

/*
Credential is a common interface that any object can implement if they want to act
as an authenticator.
*/
type Credential interface {
	Identifier() string
	Secret() string
}

/*
NewCredential constructs a Credential of the type that is passed in.
*/
func NewCredential(credentialType Credential) Credential {
	return credentialType
}

/*
UsernamePassword is a very classic type of authentication mechanism.
*/
type UsernamePassword struct {
	Username string
	Password string
}

/*
Identifier is the username part of this Credential type.
*/
func (credential UsernamePassword) Identifier() string {
	return credential.Username
}

/*
Secret is the password part of this Credential type
*/
func (credential UsernamePassword) Secret() string {
	return credential.Password
}
