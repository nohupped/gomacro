/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017-2018 Massimiliano Ghilardi
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Lesser General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Lesser General Public License for more details.
 *
 *     You should have received a copy of the GNU Lesser General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/lgpl>.
 *
 *
 * importer.go
 *
 *  Created on Feb 27, 2017
 *      Author Massimiliano Ghilardi
 */

package base

import (
	"bytes"
	"errors"
	"fmt"
	"go/importer"
	"go/types"
	"io/ioutil"
	"os"
	r "reflect"

	"github.com/cosmos72/gomacro/imports"
)

type ImportMode int

const (
	// ImBuiltin import mechanism is:
	// 1. write a file $GOPATH/src/github.com/cosmos72/gomacro/imports/$PKGPATH.go containing a single func init()
	//    i.e. *inside* gomacro sources
	// 2. tell the user to recompile gomacro
	ImBuiltin ImportMode = iota

	// ImThirdParty import mechanism is the same as ImBuiltin, except that files are created in a thirdparty/ subdirectory:
	// 1. write a file $GOPATH/src/github.com/cosmos72/gomacro/imports/thirdparty/$PKGPATH.go containing a single func init()
	//    i.e. *inside* gomacro sources
	// 2. tell the user to recompile gomacro
	ImThirdParty

	// ImPlugin import mechanism is:
	// 1. write a file $GOPATH/src/gomacro_imports/$PKGPATH/$PKGNAME.go containing a var Packages map[string]Package
	//    and a single func init() to populate it
	// 2. invoke "go build -buildmode=plugin" on the file to create a shared library
	// 3. load such shared library with plugin.Open().Lookup("Packages")
	ImPlugin

	// ImInception import mechanism is:
	// 1. write a file $GOPATH/src/$PKGPATH/x_package.go containing a single func init()
	//    i.e. *inside* the package to be imported
	// 2. tell the user to recompile $PKGPATH
	ImInception
)

type Importer struct {
	from       types.ImporterFrom
	compat     types.Importer
	srcDir     string
	mode       types.ImportMode
	PluginOpen r.Value // = reflect.ValueOf(plugin.Open)
}

func DefaultImporter() *Importer {
	imp := Importer{}
	compat := importer.Default()
	if from, ok := compat.(types.ImporterFrom); ok {
		imp.from = from
	} else {
		imp.compat = compat
	}
	return &imp
}

func (imp *Importer) setPluginOpen() bool {
	if imp.PluginOpen == Nil {
		imp.PluginOpen = imports.Packages["plugin"].Binds["Open"]
		if imp.PluginOpen == Nil {
			imp.PluginOpen = None // cache the failure
		}
	}
	return imp.PluginOpen != None
}

func (imp *Importer) Import(path string) (*types.Package, error) {
	return imp.ImportFrom(path, imp.srcDir, imp.mode)
}

func (imp *Importer) ImportFrom(path string, srcDir string, mode types.ImportMode) (*types.Package, error) {
	if imp.from != nil {
		return imp.from.ImportFrom(path, srcDir, mode)
	} else if imp.compat != nil {
		return imp.compat.Import(path)
	} else {
		return nil, errors.New(fmt.Sprintf("importer.Default() returned nil, cannot import %q", path))
	}
}

// LookupPackage returns a package if already present in cache
func (g *Globals) LookupPackage(name, path string) *PackageRef {
	pkg, found := imports.Packages[path]
	if !found {
		return nil
	}
	return &PackageRef{Package: pkg, Name: name, Path: path}
}

func (g *Globals) ImportPackage(name, path string) *PackageRef {
	ref := g.LookupPackage(name, path)
	if ref != nil {
		return ref
	}
	gpkg, err := g.Importer.Import(path) // loads names and types, not the values!
	if err != nil {
		g.Errorf("error loading package %q metadata, maybe you need to download (go get), compile (go build) and install (go install) it? %v", path, err)
	}
	var mode ImportMode
	switch name {
	case "_b":
		mode = ImBuiltin
	case "_i":
		mode = ImInception
	case "_3":
		mode = ImThirdParty
	default:
		havePluginOpen := g.Importer.setPluginOpen()
		if havePluginOpen {
			mode = ImPlugin
		} else {
			mode = ImThirdParty
		}
	}
	file := g.createImportFile(path, gpkg, mode)
	ref = &PackageRef{Name: name, Path: path}
	if len(file) == 0 || mode != ImPlugin {
		// either the package exports nothing, or user must rebuild gomacro.
		// in both cases, still cache it to avoid recreating the file.
		imports.Packages[path] = ref.Package
		return ref
	}
	soname := g.compilePlugin(file, g.Stdout, g.Stderr)
	ipkgs := g.loadPluginSymbol(soname, "Packages")
	pkgs := *ipkgs.(*map[string]imports.PackageUnderlying)

	// cache *all* found packages for future use
	imports.Packages.Merge(pkgs)

	// but return only requested one
	pkg, found := imports.Packages[path]
	if !found {
		g.Errorf("error loading package %q: the compiled plugin %q does not contain it! internal error? %v", path, soname)
	}
	ref.Package = pkg
	return ref
}

func (g *Globals) createImportFile(path string, pkg *types.Package, mode ImportMode) string {
	file := g.computeImportFilename(path, mode)

	buf := bytes.Buffer{}
	isEmpty := g.writeImportFile(&buf, path, pkg, mode)
	if isEmpty {
		g.Warnf("package %q exports zero constants, functions, types and variables", path)
		return ""
	}

	err := ioutil.WriteFile(file, buf.Bytes(), os.FileMode(0666))
	if err != nil {
		g.Errorf("error writing file %q: %v", file, err)
	}
	if mode == ImPlugin {
		g.Debugf("created file %q...", file)
	} else {
		g.Warnf("created file %q, recompile gomacro to use it", file)
	}
	return file
}

func sanitizeIdentifier(str string) string {
	return sanitizeIdentifier2(str, '_')
}

func sanitizeIdentifier2(str string, replacement rune) string {
	runes := []rune(str)
	for i, ch := range runes {
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' ||
			(i != 0 && ch >= '0' && ch <= '9') {
			continue
		}
		runes[i] = replacement
	}
	return string(runes)
}

func (g *Globals) computeImportFilename(path string, mode ImportMode) string {
	gosrc_dir := GoSrcPath()

	switch mode {
	case ImBuiltin:
		// user will need to recompile gomacro
		return Subdir(gosrc_dir, GomacroDir, "imports", sanitizeIdentifier(path)+".go")
	case ImInception:
		// user will need to recompile gosrc_dir / path
		return Subdir(gosrc_dir, path, "x_package.go")
	case ImThirdParty:
		// either plugin.Open is not available, or user explicitly requested import _3 "package".
		// In both cases, user will need to recompile gomacro
		return Subdir(gosrc_dir, GomacroDir, "imports", "thirdparty", sanitizeIdentifier(path)+".go")
	}

	file := FileName(path) + ".go"
	file = Subdir(gosrc_dir, "gomacro_imports", path, file)
	dir := DirName(file)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		Errorf("error creating directory %q: %v", dir, err)
	}
	return file
}
