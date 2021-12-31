package spdg

/*
ContextError is a custom error type we can use to get strong typed errors instead of Go's string approach.
*/
type ContextError uint

const (
	// CONTEXTOK represents no error, or error == nil.
	CONTEXTOK ContextError = iota
	// VALIDATIONERROR is used when the Validate method discovers invalid data.
	VALIDATIONERROR
)

/*
ContextChecker gives us a custom type we can define methods on. There would be no other way to define
a method on an empty interface (I think...).
*/
type ContextChecker struct {
	value     interface{}
	innerType string
}

/*
Check wraps our value into a ContextChecker so we can call a validation method on it.
*/
func Check(value interface{}, innerType string) *ContextChecker {
	return &ContextChecker{
		value:     value,
		innerType: innerType,
	}
}

/*
IsNot validates if the wrapped value inside the ContextChecker is not the same as the value we pass in.
*/
func (checker *ContextChecker) IsNot(proposal interface{}) ContextError {
	// This will reflect on the underlying type of the custom empty interface type.
	switch checker.innerType {
	case "string":
		if checker.value.(string) == proposal.(string) {
			return VALIDATIONERROR
		}
	case "ContextRole":
		if checker.value.(ContextRole) == proposal.(ContextRole) {
			return VALIDATIONERROR
		}
	}

	return CONTEXTOK
}
