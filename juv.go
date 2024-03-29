package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const scansText = `
// generated by juv; DO NOT EDIT
package {{.PackageName}}
import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
)

{{range .Tokens}}func(r *{{.Name}}) UnmarshalJSON(b []byte) error {
	type Alias {{.Name}} // avoid stack over flow error
	var a Alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	if err := validator.New().Struct(a); err != nil {
		return err
	}

	{{range .Fields}}r.{{.Name}} = a.{{.Name}}
	{{end}}
	
	return nil
}
{{end}}
`

type structToken struct {
	Name   string
	Fields []fieldToken
}

type fieldToken struct {
	Name string
	Type string
}

func main() {
	//log.SetFlags(0)

	outFilename := flag.String("o", "juv_gen.go", "-o is output file name")
	packName := flag.String("p", "current directory", "-p is package name")
	flag.StringVar(outFilename, "output", "juv_gen.go", "-output is output file name")
	flag.StringVar(packName, "package", "current directory", "-package is package name")
	flag.Parse()

	if *packName == "current directory" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal("couldn't get working directory:", err)
		}
		*packName = filepath.Base(wd)
	}

	files, err := findFiles(flag.Args())
	if err != nil {
		log.Fatal("couldn't find files:", err)
	}

	structTokens := make([]structToken, 0, 8)
	for _, f := range files {
		tokens, err := parseCode(f)
		if err != nil {
			log.Println(`"syntax error" - parser probably`)
			log.Fatal(err)
		}
		structTokens = append(structTokens, tokens...)
	}

	if err := generate(*outFilename, *packName, structTokens); err != nil {
		log.Fatal("couldn't generate f:", err)
	}

}

func findFiles(paths []string) ([]string, error) {
	if len(paths) < 1 {
		return nil, errors.New("no starting paths")
	}

	// using map to prevent duplicate file path entries
	// in case the user accidentally passes the same file path more than once probably because of autocomplete
	files := make(map[string]struct{})
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			files[path] = struct{}{}
			continue
		}

		_ = filepath.Walk(path, func(fp string, fi os.FileInfo, _ error) error {
			if fi.IsDir() {
				// will still enter directory
				return nil
			} else if fi.Name()[0] == '.' {
				// ignore dot files
				return nil
			}

			// add file path to files
			files[fp] = struct{}{}
			return nil
		})
	}

	deduped := make([]string, 0, len(files))
	for f := range files {
		deduped = append(deduped, f)
	}

	return deduped, nil
}

func parseCode(source string) ([]structToken, error) {
	structTokens := make([]structToken, 0, 8)

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, source, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, decl := range astf.Decls {
		genDecl, isGeneralDeclaration := decl.(*ast.GenDecl)
		if !isGeneralDeclaration {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, isTypeDeclaration := spec.(*ast.TypeSpec)
			if !isTypeDeclaration {
				continue
			}

			structType, isStructTypeDeclaration := typeSpec.Type.(*ast.StructType)
			if !isStructTypeDeclaration {
				continue
			}

			// !found a struct in the source code!
			structToken := structToken{
				Name:   typeSpec.Name.Name,
				Fields: make([]fieldToken, 0, len(structType.Fields.List)),
			}

			// iterate through struct fields (1 line at a time)
			for _, fieldLine := range structType.Fields.List {

				// get filed type
				var fieldType string
				switch typeToken := fieldLine.Type.(type) {
				case *ast.Ident:
					// simple types, e.g. bool, int
					fieldType = parseIdent(typeToken)
				case *ast.SelectorExpr:
					// struct fields, e.g. time.Time, sql.NullString
					fieldType = parseSelector(typeToken)
				case *ast.ArrayType:
					// arrays
					fieldType = parseArray(typeToken)
				case *ast.StarExpr:
					// pointers
					fieldType = parseStar(typeToken)
				}

				if fieldType == "" {
					continue
				}

				// multi fields can be declare in the same line
				fieldTokens := make([]fieldToken, 0, len(fieldLine.Names))
				// get field name (or names because multiple vars can be declared in 1 line)
				for _, fieldName := range fieldLine.Names {
					fieldTokens = append(fieldTokens, fieldToken{
						Name: parseIdent(fieldName),
						Type: fieldType,
					})
				}

				structToken.Fields = append(structToken.Fields, fieldTokens...)
			}
			structTokens = append(structTokens, structToken)
		}
	}

	return structTokens, nil
}

func parseIdent(fieldType *ast.Ident) string {
	// return like byte, string, int
	return fieldType.Name
}

func parseSelector(fieldType *ast.SelectorExpr) string {
	// return like time.Time, sql.NullString
	ident, isIdent := fieldType.X.(*ast.Ident)
	if !isIdent {
		return ""
	}
	return fmt.Sprintf("%s.%s", parseIdent(ident), fieldType.Sel.Name)
}

func parseArray(fieldType *ast.ArrayType) string {
	// return like []byte, []time.Time, []*byte, []*sql.NullString
	var arrayType string

	switch typeToken := fieldType.Elt.(type) {
	case *ast.Ident:
		arrayType = parseIdent(typeToken)
	case *ast.SelectorExpr:
		arrayType = parseSelector(typeToken)
	case *ast.StarExpr:
		arrayType = parseStar(typeToken)
	}

	if arrayType == "" {
		return ""
	}

	return fmt.Sprintf("[]%s", arrayType)
}

func parseStar(fieldType *ast.StarExpr) string {
	// return like *bool, *time.Time, *[]byte, and other array stuff
	var starType string

	switch typeToken := fieldType.X.(type) {
	case *ast.Ident:
		starType = parseIdent(typeToken)
	case *ast.SelectorExpr:
		starType = parseSelector(typeToken)
	case *ast.ArrayType:
		starType = parseArray(typeToken)
	}

	if starType == "" {
		return ""
	}

	return fmt.Sprintf("*%s", starType)
}

func generate(outFile, pkg string, tokens []structToken) error {
	if len(tokens) < 1 {
		return errors.New("no structs found")
	}

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	fnMap := template.FuncMap{"title": strings.Title}
	scansTmpl, err := template.New("juv").Funcs(fnMap).Parse(scansText)
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)
	bind := map[string]interface{}{"PackageName": pkg, "Tokens": tokens}
	if err := scansTmpl.Execute(buff, bind); err != nil {
		return err
	}

	source, err := format.Source(buff.Bytes())
	if err != nil {
		return nil
	}

	_, err = io.Copy(out, bytes.NewReader(source))
	return err
}
