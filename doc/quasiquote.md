Quasiquote
==========

implementing quasiquote, unquote and unquote_splice in Go
--------------------------------------------------------

One of the main motivations behind the creation of Go interpreter `gomacro`
was to add Lisp-like macros to Go.

This includes implementing Common Lisp `quote`, `quasiquote` and, more crucially,
`unquote` and `unquote_splice` i.e. Common Lisp macro characters `'` `` ` `` `,` and `,@`

Since Go language is not homoiconic, i.e. (source) code and (program) data
are not represented identically, this is a challenge.

### Parser ###

The first (moderate) difficulty is adding support for `'` `` ` `` `,` and `,@` to Go parser.
It was solved by forking Go standard packages https://golang.org/pkg/go/scanner/
and https://golang.org/pkg/go/parser/ and patching them.

Characters `'` `` ` `` and `,` are already reserved in Go,
so the author decided to replace them as follows:
* quote          `'`  must be written `~'`
* quasiquote `` ` ``  must be written `~"` (not ``~` `` because the latter messes up syntax hilighting in Go-aware editors and IDEs - starts a multiline raw string)
* unquote        `,`  must be written `~,`
* unquote_splice `,@` must be written `~,@`

the prefix `~` is configurabile when manually instantiating the modified parser.

Go parser produces as output an abstract syntax tree (AST) represented as a tree of `ast.Node`,
from the standard package https://golang.org/pkg/go/ast/

Defining new node types is deliberately impossible (`ast.Node` is an interface with unexported methods),
luckily the existing types are flexible enough to accommodate the new syntax.

The chosen representation is somewhat cumbersome but fully general: newly created constants `token.QUOTE`, `token.QUASIQUOTE`,
`token.UNQUOTE` and `token.UNQUOTE_SPLICE` are used as unary operators on a fictitious closure containing the quoted code.
Examples:
* `'x` must be written `~'x` and is parsed as if written `~' func() { x }`
* `` `{x = y}`` must be written `~"{x = y}` and is parsed as if written `~" func() { x = y }`
* `,{1 + 2}`  must be written `~,{1 + 2}` and is parsed as if written `~, func() { 1 + 2 }`
* `,@{foo()}` must be written `~,@{foo()}` and is parsed as if written `~,@ func() { foo() }`

The fictitious closures are necessary because `ast.UnaryExpr` only allows an expression as its operand - not arbitrary
statements or declarations.
In Go, the only expression that can contain arbitrary statements and declarations is a closure (in Go terms, a "function literal")

### Classic interpreter ###

`gomacro` contains two interpreters: "classic" and "fast".

The classic interpreter is compact (about 5k LOC) and directly executes the AST, producing `reflect.Value` objects as output.
It is also quite slow (1000-2000 times slower than compiled Go), due to the overhead of continuously dispatching on the type
and contents of `ast.Node` and working with `reflect.Value` instead of native Go types.

One significant advantage of directly executing the AST is the simplicity of quasiquote implementation:
it visits depth-first the whole AST, looking for `ast.UnaryExpr` whose operator is `token.QUOTE`, `token.QUASIQUOTE`,
`token.UNQUOTE` or `token.UNQUOTE_SPLICE`, and performs the corresponding operation (either return the quoted code literally or evaluate it)
while keeping track of the current quasiquotation depth (the number of entered `~"` minus the number of entered `~,` and `~,@`)

### Fast interpreter ###

The second, "fast" interpreter included in `gomacro` is more sophisticated.

Instead of directly executing the AST, it splits the execution in two phases:
1. visits the AST depth-first and create a tree of closures - one for each expression to be executed.
   For example, `a + b` is transformed into something equivalent to:
   ```
   var a = resolve("a").(func(env *Env) int)
   var b = resolve("b").(func(env *Env) int)
   var sum_ab = func(env *Env) int {
	   return a(env) + b(env)
   }
   ```
   The fast interpreter also performs type checking and type inference while creating this tree of closures.

   Statements (including declarations) are transformed a bit differently: each becomes a closure
   executing the statement in the interpreter, and returning the next statement to be executed.
   For example, `if x { foo() } else { bar() }` is transformed into something equivalent to:
   ```
   var x = resolve("x").(func(env *Env) bool)
   var foo = resolve("foo").(func(env *Env) func())
   var bar = resolve("bar").(func(env *Env) func())
   var ip_then, ip_else, ip_finish int // will be set below
   Code.Append(func(env *Env) (Stmt, *Env) {
	  var ip int
	  if x(env) {
		 ip = ip_then // goto ip_then
	  } else {
		 ip = ip_else // goto ip_else
	  }
	  env.Code.IP = ip
	  return env.Code[ip], env
   })
   ip_then = Code.Len()
   Code.Append(func(env *Env) (Stmt, *Env) {
	   foo(env)()
	   env.Code.IP = ip_finish // goto ip_finish i.e. skip else branch
	   return env.Code[ip_finish], env
   })
   ip_else = Code.Len()
   Code.Append(func(env *Env) (Stmt, *Env) {
	   bar(env)()
	   env.Code.IP = ip_finish // can also be written env.Code.IP++
	   return env.Code[ip_finish], env
   })
   ip_finish = Code.Len()
   ```
   Note the extensive use of closures, i.e. anonymous functions that access **mutable** variables
   of the surrounding scope: `x` `foo` `bar` `ip_then` `ip_else` and `ip_finish`.

2) executes the created closures

"fast" interpreter also uses native Go types where possible, to further speed up execution
and reduce the reliance on `reflect.Value`.

The result is a much larger interpreter:
* 20k LOC written manually
* plus further 80k LOC, generated from 8k LOC of macros, by using the "classic" interpreter and its quote,
  quasiquote and macros as a code generation tool

It is also significantly faster than the "classic" interpreter:
on most microbenchmarks, it's 10-100 times slower than compiled code, instead of 1000-2000 times slower.