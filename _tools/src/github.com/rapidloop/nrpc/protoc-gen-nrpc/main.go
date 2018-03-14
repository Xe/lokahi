package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/rapidloop/nrpc"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// baseName returns the last path element of the name, with the last dotted suffix removed.
func baseName(name string) string {
	// First, find the last element
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	// Now drop the suffix
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[0:i]
	}
	return name
}

// getGoPackage returns the file's go_package option.
// If it containts a semicolon, only the part before it is returned.
func getGoPackage(fd *descriptor.FileDescriptorProto) string {
	pkg := fd.GetOptions().GetGoPackage()
	if strings.Contains(pkg, ";") {
		parts := strings.Split(pkg, ";")
		if len(parts) > 2 {
			log.Fatalf(
				"protoc-gen-nrpc: go_package '%s' contains more than 1 ';'",
				pkg)
		}
		pkg = parts[0]
	}

	return pkg
}

// goPackageOption interprets the file's go_package option.
// If there is no go_package, it returns ("", "", false).
// If there's a simple name, it returns ("", pkg, true).
// If the option implies an import path, it returns (impPath, pkg, true).
func goPackageOption(d *descriptor.FileDescriptorProto) (impPath, pkg string, ok bool) {
	pkg = getGoPackage(d)
	if pkg == "" {
		return
	}
	ok = true
	// The presence of a slash implies there's an import path.
	slash := strings.LastIndex(pkg, "/")
	if slash < 0 {
		return
	}
	impPath, pkg = pkg, pkg[slash+1:]
	// A semicolon-delimited suffix overrides the package name.
	sc := strings.IndexByte(impPath, ';')
	if sc < 0 {
		return
	}
	impPath, pkg = impPath[:sc], impPath[sc+1:]
	return
}

// goPackageName returns the Go package name to use in the
// generated Go file.  The result explicit reports whether the name
// came from an option go_package statement.  If explicit is false,
// the name was derived from the protocol buffer's package statement
// or the input file name.
func goPackageName(d *descriptor.FileDescriptorProto) (name string, explicit bool) {
	// Does the file have a "go_package" option?
	if _, pkg, ok := goPackageOption(d); ok {
		return pkg, true
	}

	// Does the file have a package clause?
	if pkg := d.GetPackage(); pkg != "" {
		return pkg, false
	}
	// Use the file base name.
	return baseName(d.GetName()), false
}

// goFileName returns the output name for the generated Go file.
func goFileName(d *descriptor.FileDescriptorProto) string {
	name := *d.Name
	if ext := path.Ext(name); ext == ".proto" || ext == ".protodevel" {
		name = name[:len(name)-len(ext)]
	}
	name += ".nrpc.go"

	// Does the file have a "go_package" option?
	// If it does, it may override the filename.
	if impPath, _, ok := goPackageOption(d); ok && impPath != "" {
		// Replace the existing dirname with the declared import path.
		_, name = path.Split(name)
		name = path.Join(impPath, name)
		return name
	}

	return name
}

func fieldGoType(td *descriptor.FieldDescriptorProto) string {
	// Use protoc-gen-go generator to get the actual go type (for plain types
	// only!)
	t, _ := (*generator.Generator)(nil).GoType(nil, td)
	// We assume proto3, but pass nil to the generator, which will assume proto2.
	// The consequence is a leading star on the type that we need to trim
	return strings.TrimPrefix(t, "*")
}

// splitMessageTypeName split a message type into (package name, type name)
func splitMessageTypeName(name string) (string, string) {
	if len(name) == 0 {
		log.Fatal("Empty message type")
	}
	if name[0] != '.' {
		log.Fatalf("Expect message type name to start with '.', but it is '%s'", name)
	}
	lastDot := strings.LastIndex(name, ".")
	return name[1:lastDot], name[lastDot+1:]
}

func splitTypePath(name string) []string {
	if len(name) == 0 {
		log.Fatal("Empty message type")
	}
	if name[0] != '.' {
		log.Fatalf("Expect message type name to start with '.', but it is '%s'", name)
	}
	return strings.Split(name[1:], ".")
}

func lookupFileDescriptor(name string) *descriptor.FileDescriptorProto {
	for _, fd := range request.GetProtoFile() {
		if fd.GetPackage() == name {
			return fd
		}
	}
	return nil
}

func lookupMessageType(name string) (*descriptor.FileDescriptorProto, *descriptor.DescriptorProto) {
	path := splitTypePath(name)

	pkgpath := path[:len(path)-1]

	var fd *descriptor.FileDescriptorProto
	for {
		pkgname := strings.Join(pkgpath, ".")
		fd = lookupFileDescriptor(pkgname)
		if fd != nil {
			break
		}
		if len(pkgpath) == 1 {
			log.Fatalf("Could not find the .proto file for package '%s' (from message %s)", pkgname, name)
		}
		pkgpath = pkgpath[:len(pkgpath)-1]
	}

	path = path[len(pkgpath):]

	var d *descriptor.DescriptorProto
	for _, mt := range fd.GetMessageType() {
		if mt.GetName() == path[0] {
			d = mt
			break
		}
	}
	if d == nil {
		log.Fatalf("No such type '%s' in package '%s'", path[0], strings.Join(pkgpath, "."))
	}
	for i, token := range path[1:] {
		var found bool
		for _, nd := range d.GetNestedType() {
			if nd.GetName() == token {
				d = nd
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("No such nested type '%s' in '%s.%s'",
				token, strings.Join(pkgpath, "."), strings.Join(path[:i+1], "."))
		}
	}
	return fd, d
}

func getField(d *descriptor.DescriptorProto, name string) *descriptor.FieldDescriptorProto {
	for _, f := range d.GetField() {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func getOneofDecl(d *descriptor.DescriptorProto, name string) *descriptor.OneofDescriptorProto {
	for _, of := range d.GetOneofDecl() {
		if of.GetName() == name {
			return of
		}
	}
	return nil
}

func pkgSubject(fd *descriptor.FileDescriptorProto) string {
	if options := fd.GetOptions(); options != nil {
		e, err := proto.GetExtension(options, nrpc.E_PackageSubject)
		if err == nil {
			value := e.(*string)
			return *value
		}
	}
	return fd.GetPackage()
}

func getResultType(
	md *descriptor.MethodDescriptorProto,
) string {
	return md.GetOutputType()
}

func getGoType(pbType string) (string, string) {
	if !strings.Contains(pbType, ".") {
		return "", pbType
	}
	fd, _ := lookupMessageType(pbType)
	name := strings.TrimPrefix(pbType, "."+fd.GetPackage()+".")
	name = strings.Replace(name, ".", "_", -1)
	return getGoPackage(fd), name
}

func getPkgImportName(goPkg string) string {
	if goPkg == getGoPackage(currentFile) {
		return ""
	}
	replacer := strings.NewReplacer(".", "_", "/", "_", "-", "_")
	return replacer.Replace(goPkg)
}

var pluginPrometheus bool

var funcMap = template.FuncMap{
	"GoPackageName": func(fd *descriptor.FileDescriptorProto) string {
		p, _ := goPackageName(fd)
		return p
	},
	"GetPkg": func(pkg, s string) string {
		s = strings.TrimPrefix(s, ".")
		s = strings.TrimPrefix(s, pkg)
		s = strings.TrimPrefix(s, ".")
		return s
	},
	"GetExtraImports": func(fd *descriptor.FileDescriptorProto) []string {
		// check all the types used and imports packages from where they come
		var imports = make(map[string]string)
		for _, sd := range fd.GetService() {
			for _, md := range sd.GetMethod() {
				goPkg, _ := getGoType(md.GetInputType())
				pkgImportName := getPkgImportName(goPkg)
				if pkgImportName != "" {
					imports[pkgImportName] = goPkg
				}
				goPkg, _ = getGoType(getResultType(md))
				pkgImportName = getPkgImportName(goPkg)
				if pkgImportName != "" {
					imports[pkgImportName] = goPkg
				}
			}
		}
		var result []string
		for importName, goPkg := range imports {
			result = append(result,
				fmt.Sprintf("%s \"%s\"",
					importName,
					goPkg,
				),
			)
		}
		return result
	},
	"GetPkgSubjectPrefix": func(fd *descriptor.FileDescriptorProto) string {
		if s := pkgSubject(fd); s != "" {
			return s + "."
		}
		return ""
	},
	"GetPkgSubject": pkgSubject,
	"GetPkgSubjectParams": func(fd *descriptor.FileDescriptorProto) []string {
		if opts := fd.GetOptions(); opts != nil {
			e, err := proto.GetExtension(opts, nrpc.E_PackageSubjectParams)
			if err == nil {
				value := e.([]string)
				return value
			}
		}
		return nil
	},
	"GetServiceSubject": func(sd *descriptor.ServiceDescriptorProto) string {
		if opts := sd.GetOptions(); opts != nil {
			s, err := proto.GetExtension(opts, nrpc.E_ServiceSubject)
			if err == nil {
				value := s.(*string)
				return *value
			}
		}
		if opts := currentFile.GetOptions(); opts != nil {
			s, err := proto.GetExtension(opts, nrpc.E_ServiceSubjectRule)
			if err == nil {
				switch *(s.(*nrpc.SubjectRule)) {
				case nrpc.SubjectRule_COPY:
					return sd.GetName()
				case nrpc.SubjectRule_TOLOWER:
					return strings.ToLower(sd.GetName())
				}
			}
		}
		return sd.GetName()
	},
	"ServiceNeedsHandler": func(sd *descriptor.ServiceDescriptorProto) bool {
		for _, md := range sd.GetMethod() {
			if md.GetInputType() != ".nrpc.NoRequest" {
				return true
			}
		}
		return false
	},
	"GetMethodSubject": func(md *descriptor.MethodDescriptorProto) string {
		if opts := md.GetOptions(); opts != nil {
			s, err := proto.GetExtension(opts, nrpc.E_MethodSubject)
			if err == nil {
				value := s.(*string)
				return *value
			}
		}
		if opts := currentFile.GetOptions(); opts != nil {
			s, err := proto.GetExtension(opts, nrpc.E_MethodSubjectRule)
			if err == nil {
				switch *(s.(*nrpc.SubjectRule)) {
				case nrpc.SubjectRule_COPY:
					return md.GetName()
				case nrpc.SubjectRule_TOLOWER:
					return strings.ToLower(md.GetName())
				}
			}
		}
		return md.GetName()
	},
	"GetMethodSubjectParams": func(md *descriptor.MethodDescriptorProto) []string {
		if opts := md.GetOptions(); opts != nil {
			if e, err := proto.GetExtension(opts, nrpc.E_MethodSubjectParams); err == nil {
				return e.([]string)
			}
		}
		return []string{}
	},
	"GetServiceSubjectParams": func(sd *descriptor.ServiceDescriptorProto) []string {
		if opts := sd.GetOptions(); opts != nil {
			if e, err := proto.GetExtension(opts, nrpc.E_ServiceSubjectParams); err == nil {
				return e.([]string)
			}
		}
		return []string{}
	},
	"HasStreamedReply": func(md *descriptor.MethodDescriptorProto) bool {
		if opts := md.GetOptions(); opts != nil {
			if e, err := proto.GetExtension(opts, nrpc.E_StreamedReply); err == nil {
				return *(e.(*bool))
			}
		}
		return false
	},
	"Prometheus": func() bool {
		return pluginPrometheus
	},
	"GetResultType": getResultType,
	"GoType": func(pbType string) string {
		goPkg, goType := getGoType(pbType)
		if goPkg != "" {
			importName := getPkgImportName(goPkg)
			if importName != "" {
				goType = importName + "." + goType
			}
		}
		return goType
	},
}

var request plugin.CodeGeneratorRequest
var currentFile *descriptor.FileDescriptorProto

func main() {

	log.SetPrefix("protoc-gen-nrpc: ")
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error: reading input: %v", err)
	}

	var response plugin.CodeGeneratorResponse
	if err := proto.Unmarshal(data, &request); err != nil {
		log.Fatalf("error: parsing input proto: %v", err)
	}

	if len(request.GetFileToGenerate()) == 0 {
		log.Fatal("error: no files to generate")
	}

	pluginPrometheus = request.GetParameter() == "plugins=prometheus"

	tmpl, err := template.New(".").Funcs(funcMap).Parse(tFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range request.GetFileToGenerate() {
		var fd *descriptor.FileDescriptorProto
		for _, fd = range request.GetProtoFile() {
			if name == fd.GetName() {
				break
			}
		}
		if fd == nil {
			log.Fatalf("could not find the .proto file for %s", name)
		}

		currentFile = fd

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, fd); err != nil {
			log.Fatal(err)
		}

		currentFile = nil

		response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(goFileName(fd)),
			Content: proto.String(buf.String()),
		})
	}

	if data, err = proto.Marshal(&response); err != nil {
		log.Fatalf("error: failed to marshal output proto: %v", err)
	}
	if _, err := os.Stdout.Write(data); err != nil {
		log.Fatalf("error: failed to write output proto: %v", err)
	}
}
