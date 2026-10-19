---
name: rust-api-guidelines
description: >
  Enforces the official Rust API Guidelines (https://rust-lang.github.io/api-guidelines/)
  when writing, reviewing, or refactoring Rust library code. Use this skill whenever you are
  writing new Rust structs, enums, traits, functions, or modules that form a public API —
  including when designing error types, implementing conversions, naming methods, adding
  documentation, creating macros, or structuring crate features. Also use it when reviewing
  existing Rust code for API quality, idiomatic patterns, or crate publishing readiness.
  Trigger on any Rust library design task, even if the user doesn't explicitly mention
  "API guidelines."
---

# Rust API Guidelines

The official Rust API Guidelines capture the Rust library team's consensus on idiomatic, interoperable crate design. They cover naming, trait implementations, documentation, type safety, macros, and more. Use the checklist below when writing, reviewing, or publishing Rust library code.

## Checklist

### Naming ([details](references/naming.md))
- [ ] Casing conforms to RFC 430 ([C-CASE])
- [ ] Ad-hoc conversions follow `as_`, `to_`, `into_` conventions ([C-CONV])
- [ ] Getter names follow Rust convention ([C-GETTER])
- [ ] Methods on collections that produce iterators follow `iter`, `iter_mut`, `into_iter` ([C-ITER])
- [ ] Iterator type names match the methods that produce them ([C-ITER-TY])
- [ ] Feature names are free of placeholder words ([C-FEATURE])
- [ ] Names use a consistent word order ([C-WORD-ORDER])

### Interoperability ([details](references/interoperability.md))
- [ ] Types eagerly implement common traits ([C-COMMON-TRAITS])
  - `Copy`, `Clone`, `Eq`, `PartialEq`, `Ord`, `PartialOrd`, `Hash`, `Debug`,
    `Display`, `Default`
- [ ] Conversions use the standard traits `From`, `AsRef`, `AsMut` ([C-CONV-TRAITS])
- [ ] Collections implement `FromIterator` and `Extend` ([C-COLLECT])
- [ ] Data structures implement Serde's `Serialize`, `Deserialize` ([C-SERDE])
- [ ] Types are `Send` and `Sync` where possible ([C-SEND-SYNC])
- [ ] Error types are meaningful and well-behaved ([C-GOOD-ERR])
- [ ] Binary number types provide `Hex`, `Octal`, `Binary` formatting ([C-NUM-FMT])
- [ ] Generic reader/writer functions take `R: Read` and `W: Write` by value ([C-RW-VALUE])

### Macros ([details](references/macros.md))
- [ ] Input syntax is evocative of the output ([C-EVOCATIVE])
- [ ] Macros compose well with attributes ([C-MACRO-ATTR])
- [ ] Item macros work anywhere that items are allowed ([C-ANYWHERE])
- [ ] Item macros support visibility specifiers ([C-MACRO-VIS])
- [ ] Type fragments are flexible ([C-MACRO-TY])

### Documentation ([details](references/documentation.md))
- [ ] Crate level docs are thorough and include examples ([C-CRATE-DOC])
- [ ] All items have a rustdoc example ([C-EXAMPLE])
- [ ] Examples use `?`, not `try!`, not `unwrap` ([C-QUESTION-MARK])
- [ ] Function docs include error, panic, and safety considerations ([C-FAILURE])
- [ ] Prose contains hyperlinks to relevant things ([C-LINK])
- [ ] Cargo.toml includes all common metadata ([C-METADATA])
  - authors, description, license, homepage, documentation, repository,
    keywords, categories
- [ ] Release notes document all significant changes ([C-RELNOTES])
- [ ] Rustdoc does not show unhelpful implementation details ([C-HIDDEN])

### Predictability ([details](references/predictability.md))
- [ ] Smart pointers do not add inherent methods ([C-SMART-PTR])
- [ ] Conversions live on the most specific type involved ([C-CONV-SPECIFIC])
- [ ] Functions with a clear receiver are methods ([C-METHOD])
- [ ] Functions do not take out-parameters ([C-NO-OUT])
- [ ] Operator overloads are unsurprising ([C-OVERLOAD])
- [ ] Only smart pointers implement `Deref` and `DerefMut` ([C-DEREF])
- [ ] Constructors are static, inherent methods ([C-CTOR])

### Flexibility ([details](references/flexibility.md))
- [ ] Functions expose intermediate results to avoid duplicate work ([C-INTERMEDIATE])
- [ ] Caller decides where to copy and place data ([C-CALLER-CONTROL])
- [ ] Functions minimize assumptions about parameters by using generics ([C-GENERIC])
- [ ] Traits are object-safe if they may be useful as a trait object ([C-OBJECT])

### Type Safety ([details](references/type-safety.md))
- [ ] Newtypes provide static distinctions ([C-NEWTYPE])
- [ ] Arguments convey meaning through types, not `bool` or `Option` ([C-CUSTOM-TYPE])
- [ ] Types for a set of flags are `bitflags`, not enums ([C-BITFLAG])
- [ ] Builders enable construction of complex values ([C-BUILDER])

### Dependability ([details](references/dependability.md))
- [ ] Functions validate their arguments ([C-VALIDATE])
- [ ] Destructors never fail ([C-DTOR-FAIL])
- [ ] Destructors that may block have alternatives ([C-DTOR-BLOCK])

### Debuggability ([details](references/debuggability.md))
- [ ] All public types implement `Debug` ([C-DEBUG])
- [ ] `Debug` representation is never empty ([C-DEBUG-NONEMPTY])

### Future Proofing ([details](references/future-proofing.md))
- [ ] Sealed traits protect against downstream implementations ([C-SEALED])
- [ ] Structs have private fields ([C-STRUCT-PRIVATE])
- [ ] Newtypes encapsulate implementation details ([C-NEWTYPE-HIDE])
- [ ] Data structures do not duplicate derived trait bounds ([C-STRUCT-BOUNDS])

### Necessities ([details](references/necessities.md))
- [ ] Public dependencies of a stable crate are stable ([C-STABLE])
- [ ] Crate and its dependencies have a permissive license ([C-PERMISSIVE])

## How to Use

When writing or reviewing Rust library code, walk through the checklist above. Each section
heading links to a reference file with full rationale, examples, and edge cases -- read those
when you need deeper context. The most commonly missed items are: missing `Debug` impls, missing
common trait impls (`Clone`, `Default`), `get_` prefixes on getters, `unwrap()` in doc
examples, missing `# Errors`/`# Panics` sections, and implementing `Into` instead of `From`.
When preparing to publish, verify `Cargo.toml` metadata, licensing, and stable dependencies.

## Reviewing Existing Code

When reviewing existing code for conformance to these guidelines:

1. **Only report violations.** Do not mention guidelines that the code already follows. The user
   only needs to know what to fix.
2. **Format each violation** as follows. Use blank lines around each heading for readability:

   ### Violation: <short guideline summary>

   A brief one-or-two-sentence explanation of what the guideline requires and why the code
   violates it. Do not repeat the full guideline text.

   - `path/to/file.rs:42`: `TypeOrFunction` — what it is now → what it should be
   - `path/to/file.rs:87`: `AnotherType` — what it is now → what it should be

   Include line numbers so users can click to jump to the location in their IDE.
3. **After listing all violations**, show a summary table like:

   | Guideline | Severity | Count |
   |-----------|----------|-------|
   | Missing Debug impl | 🔴 High | 3 |
   | Getter naming | 🟡 Medium | 2 |

   Rank rows from most severe to least severe.

   Severity levels: **🔴 High** (breaks interoperability, safety, or correctness), **🟡 Medium**
   (hurts usability or idiomatic style), **🟢 Low** (cosmetic or minor convention).
