package types

type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

func SortDirectionFromString(value string) SortDirection {
	direction := SortDirection(value)

	switch direction {
	case Asc, Desc:
		return direction
	default:
		return Asc
	}
}
