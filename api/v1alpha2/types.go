// +groupName=foo.example.com
package v1alpha2

type Foo struct {
	Spec   FooSpec   `json:"spec,omitempty"`
	Status FooStatus `json:"status,omitempty"`
}

type FooStatus struct {
	Size int32 `json:"size"`
}
