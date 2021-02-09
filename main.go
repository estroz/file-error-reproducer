package main

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/tools/go/packages"
	"sigs.k8s.io/controller-tools/pkg/crd"
	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	roots, err := loader.LoadRoots(filepath.Join(wd, "api") + "/...")
	if err != nil {
		log.Fatal(err)
	}
	registry := &markers.Registry{}
	if err := crdmarkers.Register(registry); err != nil {
		log.Fatal(err)
	}
	ctx := &genall.GenerationContext{
		Collector: &markers.Collector{
			Registry: registry,
		},
		Roots:     roots,
		InputRule: genall.InputFromFileSystem,
		Checker:   &loader.TypeChecker{},
	}

	types := make(map[crd.TypeIdent]*markers.TypeInfo)
	for _, root := range ctx.Roots {
		fmt.Println("root", root)

		cb := func(info *markers.TypeInfo) {
			ident := crd.TypeIdent{
				Package: root,
				Name:    info.Name,
			}
			fmt.Println("\tadding ident", ident)
			types[ident] = info
		}
		if err := markers.EachType(ctx.Collector, root, cb); err != nil {
			root.AddError(err)
		}
	}

	if loader.PrintErrors(ctx.Roots, packages.TypeError) {
		log.Fatal("one or more API packages had type errors")
	}

	fmt.Println()

	var fooIdents []crd.TypeIdent
	for ident := range types {
		if ident.Name == "Foo" {
			fooIdents = append(fooIdents, ident)
		}
	}
	sort.Slice(fooIdents, func(i, j int) bool {
		return fooIdents[i].Package.String() < fooIdents[j].Package.String()
	})

	for _, fooIdent := range fooIdents {
		foo := types[fooIdent]
		fmt.Printf("%s.%s\n", fooIdent.Package, fooIdent.Name)
		nextFields := foo.Fields
		for len(nextFields) > 0 {
			fields := []markers.FieldInfo{}
			for _, field := range nextFields {
				fmt.Println("\tfield:", field.Name)
				ast.Inspect(field.RawField, func(n ast.Node) bool {
					if n == nil {
						return true
					}

					var info *markers.TypeInfo
					var ident crd.TypeIdent
					var hasInfo bool
					switch nt := n.(type) {
					case *ast.Ident:
						if nt.Obj != nil && nt.Obj.Kind == ast.Typ {
							ident = crd.TypeIdent{Package: fooIdent.Package, Name: nt.Name}
							info, hasInfo = types[ident]
						}
					}
					if !hasInfo {
						return true
					}

					fmt.Println("\t\tfound type info:", ident.Name)
					for i := range info.Fields {
						fmt.Println("\t\t\tadding next field:", info.Fields[i].Name)
						fields = append(fields, info.Fields[i])
					}

					return true
				})
			}
			nextFields = fields
		}
		fmt.Println()
	}

}
