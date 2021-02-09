# file-error-reproducer

## Purpose and setup

`api/` contains two sub-packages, `v1alpha1` and `v1alpha2`, containing the `Foo` struct type and its children.
Type-wise these packages are the same except for package name. Organization-wise `v1alpha2`'s `FooSpec`
and child types are in a separate file from the main `types.go` file.

```sh
api/
├── v1alpha1
│   └── types.go # Contains all types
└── v1alpha2
    ├── spec.go  # Contains FooSpec and SomeType
    └── types.go # Contains Foo and FooStatus
```

## Results

Environment:

```console
$ go version
go version go1.15.2 linux/amd64
```

Running `go run main.go` will output the following:

```
root github.com/example/file-error-reproducer/api/v1alpha1
	adding ident "github.com/example/file-error-reproducer/api/v1alpha1".Foo
	adding ident "github.com/example/file-error-reproducer/api/v1alpha1".FooStatus
	adding ident "github.com/example/file-error-reproducer/api/v1alpha1".FooSpec
	adding ident "github.com/example/file-error-reproducer/api/v1alpha1".SomeType
root github.com/example/file-error-reproducer/api/v1alpha2
	adding ident "github.com/example/file-error-reproducer/api/v1alpha2".FooSpec
	adding ident "github.com/example/file-error-reproducer/api/v1alpha2".SomeType
	adding ident "github.com/example/file-error-reproducer/api/v1alpha2".Foo
	adding ident "github.com/example/file-error-reproducer/api/v1alpha2".FooStatus

github.com/example/file-error-reproducer/api/v1alpha1.Foo
	field: Spec
		found type info: FooSpec
			adding next field: Type
	field: Status
		found type info: FooStatus
			adding next field: Size
	field: Type
		found type info: SomeType
			adding next field: SomeField
	field: Size
	field: SomeField

github.com/example/file-error-reproducer/api/v1alpha2.Foo
	field: Spec
	field: Status
		found type info: FooStatus
			adding next field: Size
	field: Size
```

The above says that the `v1alpha2` package was parsed differently than `v1alpha1`, since `ast.Inspect`
on each `Foo` type cannot find the types for fields of `FooSpec`.
