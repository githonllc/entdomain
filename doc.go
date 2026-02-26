// Package entdomain provides an [entgo.io/ent] extension that generates
// request/response DTOs and service/handler scaffolding from annotated Ent schemas.
//
// It serves two roles:
//
//   - At code-generation time (go generate), it produces request/response DTOs,
//     BaseService structs, and BaseHandler structs for each Ent schema
//     annotated with EntDomain markers.
//
//   - At runtime, it provides the types and helpers that the generated
//     code depends on: [PageInfo], error sentinel values, and pointer utilities.
//
// # Quick Start
//
// Annotate fields in your Ent schema:
//
//	field.String("name").
//	    Annotations(entdomain.DefaultField().
//	        WithRequired(entdomain.ScopeCreate))
//
// Wire the extension in your entc.go:
//
//	func main() {
//	    ext := entdomain.NewExtensionWithOptions(
//	        entdomain.WithEntDomainPackage("github.com/githonllc/entdomain"),
//	        entdomain.WithBaseService(true),
//	        entdomain.WithBaseHandler(true),
//	    )
//	    if err := entc.Generate("./schema", &gen.Config{}, entc.Extensions(ext)); err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// Run go generate to produce {entity}_dto.go, {entity}_base_service.go,
// and {entity}_base_handler.go for each annotated schema.
//
// See the README for the full annotation reference and generated code examples.
package entdomain
