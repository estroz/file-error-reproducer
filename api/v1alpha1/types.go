// +groupName=foo.example.com
package v1alpha1

type Foo struct {
	Spec   FooSpec   `json:"spec,omitempty"`
	Status FooStatus `json:"status,omitempty"`
}

type FooStatus struct {
	Size int32 `json:"size"`
}

type FooSpec struct {
	Type SomeType `json:"typ"`
}

type SomeType struct {
	SomeField string `json:"someField"`
}
