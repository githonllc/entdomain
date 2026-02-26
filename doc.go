// Package entdomain provides an [entgo.io/ent] extension that generates
// Domain-Driven Design (DDD) layer code from annotated Ent schemas.
//
// It serves two roles:
//
//   - At code-generation time (go generate), it produces domain models,
//     repository interfaces, and service interfaces for each Ent schema
//     annotated with EntDomain markers.
//
//   - At runtime, it provides the types and interfaces that the generated
//     code depends on: [ID], [DomainModel], [Repository], [GenericDomainService],
//     [ListRequest], [SearchRequest], and the generic [BaseGenericDomainService].
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
//	        entdomain.WithOutputDir("./ent"),
//	        entdomain.WithPackageName("ent"),
//	    )
//	    if err := entc.Generate("./schema", &gen.Config{}, entc.Extensions(ext)); err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// Run go generate to produce {entity}_domain_model.go, {entity}_domain_repository.go,
// and {entity}_domain_service.go for each annotated schema.
//
// See the README for the full annotation reference and generated code examples.
package entdomain
