package berrt

/*
ObjectType enumerates a list of valid types that can be defined in
a PlantUML script.
*/
type ObjectType string

const (
	Participant ObjectType = "participant"
	Actor       ObjectType = "actor"
	Boundary    ObjectType = "boudary"
	Entity      ObjectType = "entity"
	Database    ObjectType = "database"
	Collections ObjectType = "collections"
	Queue       ObjectType = "queue"
)

/*
Object defines the properties that can be used for objects in a
PlantUML script.
*/
type Object struct {
	ID   string
	Type ObjectType
}

/*
NewObject constructs an Object type and returns a pointer reference
to the instance.
*/
func NewObject(ID string, Type ObjectType) *Object {
	return &Object{ID, Type}
}
