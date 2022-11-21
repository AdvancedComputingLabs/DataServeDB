package constraints

type HasId interface {
	GetId() int
}

type HasDbTypeDisplayName interface {
	GetDbTypeDisplayName() string
}
