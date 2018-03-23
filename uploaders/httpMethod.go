package uploaders

import "fmt"

type httpMethod int

//
const (
	Post httpMethod = 1
	Get  httpMethod = 2
	Put  httpMethod = 3
)

func (e httpMethod) String() string {
	switch e {
	case Post:
		return "POST"
	case Get:
		return "GET"
	case Put:
		return "PUT"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}
