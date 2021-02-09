package v1alpha2

type FooSpec struct {
	Type SomeType `json:"typ"`
}

type SomeType struct {
	SomeField string `json:"someField"`
}
