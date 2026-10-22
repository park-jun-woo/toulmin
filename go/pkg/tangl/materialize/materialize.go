// Package materialize resolves syntax-directed Definition values before evaluation or generation.
package materialize

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// DefinitionLoadRequest describes one file-backed definition load.
type DefinitionLoadRequest struct {
	DefinitionName string
	TypeName       string
	SourcePath     string
	ResolvedPath   string
}

// DefinitionSnapshot is the portable result shared by eval and gen.
type DefinitionSnapshot struct {
	Strings []string
}

// DefinitionLoader reads and parses one file-backed definition.
type DefinitionLoader interface {
	LoadDefinition(DefinitionLoadRequest) (DefinitionSnapshot, error)
}

// DefinitionLoaderRegistry resolves a surface loader type name.
type DefinitionLoaderRegistry interface {
	LookupDefinitionLoader(typeName string) (DefinitionLoader, bool)
}

// Lists eagerly materializes every syntax-directed list Definition exactly once.
func Lists(doc *ast.Document, registry DefinitionLoaderRegistry) (map[string]toulmin.StringListSpec, error) {
	out := make(map[string]toulmin.StringListSpec)
	for _, def := range doc.Defs {
		switch def.Kind {
		case ast.StringListDef:
			out[def.Name] = toulmin.NewStringListSpec(def.Name, def.Strings)
		case ast.FileSpecDef:
			if registry == nil {
				return nil, defError(doc, def, "loader registry is required")
			}
			loader, ok := registry.LookupDefinitionLoader(def.LoadType)
			if !ok || loader == nil {
				return nil, defError(doc, def, "loader %q is not registered", def.LoadType)
			}
			resolved, err := resolvePath(doc.Path, def.FilePath)
			if err != nil {
				return nil, defWrap(doc, def, err)
			}
			snapshot, err := loader.LoadDefinition(DefinitionLoadRequest{
				DefinitionName: def.Name, TypeName: def.LoadType,
				SourcePath: def.FilePath, ResolvedPath: resolved,
			})
			if err != nil {
				return nil, fmt.Errorf("%s:%d: definition %q (loader %q, path %q): %w", doc.Path, def.Line, def.Name, def.LoadType, def.FilePath, err)
			}
			out[def.Name] = toulmin.NewStringListSpec(def.Name, snapshot.Strings)
		}
	}
	return out, nil
}

func resolvePath(documentPath, sourcePath string) (string, error) {
	if filepath.IsAbs(sourcePath) {
		return filepath.Clean(sourcePath), nil
	}
	if documentPath == "" {
		return "", fmt.Errorf("relative file path requires Document.Path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(documentPath), sourcePath)), nil
}

func defError(doc *ast.Document, def ast.Definition, format string, args ...any) error {
	return fmt.Errorf("%s:%d: definition %q: %s", doc.Path, def.Line, def.Name, fmt.Sprintf(format, args...))
}

func defWrap(doc *ast.Document, def ast.Definition, cause error) error {
	return fmt.Errorf("%s:%d: definition %q: %w", doc.Path, def.Line, def.Name, cause)
}
