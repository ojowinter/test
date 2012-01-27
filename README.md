GoScript
========

Compiles transforming Go into JavaScript so you can continue using a clean and
concise sintaxis.  
In fact, it is used a subset of Go since JavaScript has not native way to
represent some types neither Go's statements, although some of them could be
emulated (but that is not my goal).

Advantages:

+ Using one only language for all development. A great advantage for a company.

+ Allows many type errors to be caught early in the development cycle, due to
static typing. (ToDo: compile to checking errors at time of compiling)

+ The mathematical expressions in the constants are calculated at the
translation stage. (ToDo)

+ The lines numbers in the unminified generated JavaScript match up with the
lines numbers in the original source file.

+ Generates minimized JavaScript.

Go sintaxis not supported:

+ Complex numbers, integers of 64 bits.
+ Function type, interface type excepting the empty interface.
+ Channels, goroutines (could be transformed to [Web Workers][workers]).
+ Built-in function *recover()*.
+ Defer statement.
+ Goto, labels. (1) In JavaScript, the labels are restricted to "for" and
"while" loops when they are called from "continue" and "break" directives so
its use is very limited, and (2) it is advised to [avoid its use][label].


**Note:** JavaScript can not actually do meaningful integer arithmetic on anything
bigger than 2^53. Also bitwise logical operations only have defined results (per
the spec) up to 32 bits.  
By this reason, the integers of 64 bits are unsupported.

[workers]: http://www.html5rocks.com/en/tutorials/workers/basics/
[label]: https://developer.mozilla.org/en/JavaScript/Reference/Statements/label#Avoid_using_labels


## Transformation

#### Pointers

In JavaScript, the array is the only object that can be referenced. So:

`*x` is `x[0]` in javascript while `&x` would simply be `x`.

Then, for any value that is addressed, it is boxed in an array, i.e.:

`var x *bool` to `var x = [false]`

**Note:** the printing of an address in Go (`&x`) results into an hexadecimal
address. Instead, in JavaScript with this emulation, it prints the value.

#### Return of multiple values

When a Go function returns more than one value then those values are put into an
array. Then, to access to the different values it is created a variable
`_` assigned to the return of the function, and the variable's names defined in
Go are used to access to each value of that array.

By example, for a Go function like this:

	sum, product := sumAndProduct(x, y)

its transformation would be:

	var _ = SumAndProduct(x, y), sum = _[0], product = _[1];

#### Library

JavaScript has several built-in functions and constants which can be transformed
from Go. They are defined in the maps *Constant*, and *Function*.

Since the Go functions *print()* and *println()* are used to debug then they
are transformed to [*console.log()*][console], which only can be used if the
JavaScript code is run in Webkit (Chrome, Safari), of Mozilla Firefox with the
plugin FireBug.

*panic()* is transformed to *throw new Error()*

[console]: http://v0.joehewitt.com/software/firebug/docs.php

#### Modularity

JavaScript has not some kind of module system built in. To simulate it, all the
code for the package is written inside an anonymous function which is called
directly. Then, it is used a helper function to give to an object (named like
the package) the values that must be exported.

By example, for a package named *foo* with names exported *Add* and *Product*:

	var foo = {}; (function() {
	// Code of your package

	_export(foo, [Add, Product])
	})();


## Installation

	goinstall << DOWNLOAD URL >>


## Operating instructions

<< INSTRUCTIONS TO RUN THE PROGRAM >>


## Contributing

If you are going to change code related to the compiler then you should run
`go test` after of each change in your forked repository. It will transform the
Go files in the directory "test"; to see the differences use `git diff`,
checking whether the change is what you were expecting.  
It is also expected to get some errors and warnings in some of them, which are
validated using the test functions for examples. See file "gojs/gojs_test.go".

Ideas:

+ Implement the new JS API for HTML5, transforming it from Go functions. See
 both maps *Constant* and *Function* in file "gojs/library.go". But you must
 be sure that the API is already implemented in both browsers Firefox and Chrome.
+ The [Dart library](http://api.dartlang.org/) could be used like inspiration to
 write web libraries, especially "dom" and "html".
+ JavaScript library to handle integers of 64 bits. Build it in Go since it
 can be transformed to JS ;)


## Vision

The great problem of the web developing is that there is one specification and
multiple implementations with the result of that different browsers have
different implementations for both DOM and JavaScript.  
So the solution is obvious; one specification and one implementation. And here
comes Go.

Go it's a fast, statically typed, compiled language that feels like a
dynamically typed, interpreted language. Its concurrency mechanisms make it easy
to get the most out of multicore. The Go compilers support three instruction
sets (AMD64, 386, ARM) and can target the FreeBSD, Linux, OS X, and Windows
operating systems; which means that it can be run in whatever system where is
running your browser.  
Ago time I had the idea (by my past in Python) of use Go like if were a script
language; the result was [GoNow](https://github.com/kless/GoNow) which compiles,
caching the executable, to run it directly if the source code has not been
changed. Now, to know if a web program has changed would not be so difficult if
were used a convention in its name; i.e. with "foo-*12.21*.go" (12 for year, 21
for month) can be got its date of releasing.

I've done another little step with the building of
[GoScript](https://github.com/kless/GoScript/), a compiler from Go to
JavaScript, although it isn't finished but its development is very advanced.

Now, I hope that somebody motivated contributes to get that Go been the next
web language. It would be necessary:

+ A DOM library implemented in Go. References:

	https://developer.mozilla.org/en/Gecko_DOM_Reference  
	http://api.dartlang.org/dom.html  
	http://api.dartlang.org/html.html

Then, if we want to use web technology to build user interfaces in desktop
applications:

+ A parser for HTML5
+ Another one for CSS3
+ Building of visual elements, using SVG (or CSS3?)

Note that the great advantage of use HTML5/CSS3 is that (1) there are many
designers that know how to use it, and (2) our work to get that it works in
wathever platform (web, desktop, mobile) would be almost zero.


## Credits

Thanks

+ To Big Yuuta for licensing the examples of his awesome
 [book](http://go-book.appspot.com/) for novices under public domain.
+ To the community of
 [comp.lang.javascript news](http://www.rhinocerus.net/forum/lang-javascript/)
 for solving me some doubts in that language.
+ To the community of the [Go group](http://groups.google.com/group/golang-nuts).
+ And to the creators and contributors of [Go](http://golang.org/).


## Copyright and licensing

*Copyright 2011  The "GoScript" Authors*. See file AUTHORS and CONTRIBUTORS.  
Unless otherwise noted, the source files are distributed under the
*Apache License, version 2.0* found in the LICENSE file.


* * *
*Generated by [GoWizard](https://github.com/kless/GoWizard)*

