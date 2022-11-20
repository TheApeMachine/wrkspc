package errnie

type Logger interface {
	Error(...any)
	Warning(...any)
	Info(...any)
	Debug(...any)
}
