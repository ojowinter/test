Transforming to JavaScript
==========================

## TODO

1. Do the type inference. There is some (incomplete) typechecker next to Go
parser. Type inference seems not to be complex, but it is quite tricky,
especially with << (see shift2.go in test directory).

2. Give types to the polymorphic functions (infix +, etc). Note that some
of arithmetic functions (+ for strings, <<, >>) may works not like in C or
your target lang, thus it would be better to convert them to functions
(__shl_32_32, __plus_str, ...) at the AST level.

3. Resolve consts. That implies infinite-precision arithmetic with complex
numbers. There are also tricky issues with type casting inside the const
declarations.

4. Resove the syntax which cannot be emited as-is to C or javascript, for
example { x := 1; { x := x + 1; /* two different x'es in a expression */ } }

5. Remove syntax sugar from AST to make further processing simpler (It may
be worth to define "Kernel Go" (same way as Kernel Mozart/Oz), which is
still compilable Go but has minimal number of construct).

 - get rid of constructions with duplicated semantic (&([100]int{}) and
new([100]int), make([]int, 50, 100) and new([100]int)[:50])

 - get rid of :=, "x := 1", "var x int = 1", "var x = 1" to be the same

 - desugar multielement assigment ( *a(),*b()=c(),d() ), see Go 1 spec for
the evaluation order.

 - desugar swap(swap(a,b))

 - implement 'switch' and 'for' via 'goto' (and get rid of
break/continue/fallthrough)

 - convert struct comparision to elementwise comparision

 - desugar named returns to local vars with uniq names

 - make zero-initialization of vars explicit

+ Substitute "var" by "let" for local variables, when browsers use ECMAScript 6

http://kishorelive.com/2011/11/22/ecmascript-6-looks-promising/


## Struct

> A struct in an object-oriented environment is a "public class". Classes
> represent objects (they become objects when instantiated) composed of
> properties and methods. Properties of classes probably have a direct
> correspondence to members of a struct type.

> While there are no classes in JS, you effectively create them when you
> create the constructor for an object in JS, in which you can assign names
> of properties and methods within the constructor, and then populate them
> with their values through arguments or constants/other assignments within
> the function block of the constructor. You would instantiate 'structs' or
> 'classes' in this manner.

[Reference](http://bytes.com/topic/javascript/answers/441203-structs-javascript)


## Map

It is defined using an object.

	`var m1 = new Object();` or `var m1 = {};`

	m1['one']='first';
	m1['two']='second';
	m1['three']='third';

To create a map with values initialized:

	var m2 = {
		1: 'Joe',
		3: 'Sam',
		8: 'Eve',
	};

+ Loop: `for(var i in m1)`


## Testing

The JavaScript output of the files in directory "test" have been checked using
[JavaScript Lint](http://javascriptlint.com/download.htm):

	wget http://javascriptlint.com/download/jsl-0.3.0-src.tar.gz
	tar xzf jsl-*.tar.gz && cd jsl-*/src && make -f Makefile.ref BUILD_OPT=1 &&
	sudo ln -s $(pwd)/Linux_All_*.OBJ/jsl /usr/local/bin/jsl && cd - &&
	sudo cp Linux_All_*.OBJ/jsl /usr/local/bin/ && cd - && rm -rf jsl-*

	jsl -process test/var.js
