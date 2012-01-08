GoJscript
=========

Compiles transforming Go into JavaScript so you can continue using a clean and
concise sintaxis.  
In fact, it is used a subset of Go since JavaScript has not native way to
represent some types neither Go's statements, although some of them could be
emulated (but that is not my goal).

Advantages:

+ Using one only language for all development. A great advantage for a company.

+ Allows many type errors to be caught early in the development cycle, due to
static typing. (ToDo: compile to checking errors at time of compiling)

+ The mathematical expressions are calculated at the translation stage. (ToDo)

+ The lines numbers in the unminified generated JavaScript match up with the
lines numbers in the original source file.

+ Generates minimized JavaScript.

Go sintaxis not supported:

+ Complex numbers, integers of 64 bits.
+ Function type, interface type excepting the empty interface.
+ Channels, goroutines (could be transformed to [Web Workers][workers]).
+ Built-in functions panic, recover.
+ Defer statement.
+ Goto, labels. (1) In JavaScript, the labels are restricted to "for" and
"while" loops when they are called from "continue" and "break" directives so
its use is very limited, and (2) it is advised to [avoid its use][label].

Status:

	const					[OK]
	itoa					[OK]
	blank identifier		[OK]
	var						[OK]
	array					[OK]
	slice					[OK]
	ellipsis				[OK]
	map						[OK]
	empty interface			[OK]
	check channel			[OK]
	struct					[OK]
	pointer					[OK]
	imports					[OK]
	functions				[OK]
	assignments in func		[OK]
	return					[OK]
	if						[OK]
	error for goroutine		[OK]
	switch					[OK]
	Comparison operators	[OK]
	Assignment operators	[OK]
	for						[OK]
	range					[OK]
	break, continue			[OK]
	fallthrough				[OK]
	JS functions			[OK]
	goto, label				[OK]
	anonymous function		[OK]
	JS constants			[OK]
	Return multiple values	[OK]
	Modularity				[OK]
	Functions init			[OK]

**Note:** JavaScript can not actually do meaningful integer arithmetic on anything
bigger than 2^53. Also bitwise logical operations only have defined results (per
the spec) up to 32 bits.  
By this reason, the integers of 64 bits are unsupported.

[workers]: http://www.html5rocks.com/en/tutorials/workers/basics/
[label]: https://developer.mozilla.org/en/JavaScript/Reference/Statements/label#Avoid_using_labels

#### Pointers

In JavaScript, the array is the only object that can be referenced. So it is
transformed:

`var x *bool` to `var x = [false];`, and `*x` to `x[0]`.

Then, for variables that are not defined at the beginning like pointers but
they are referenced to be used i.e. into a function whose parameter is a pointer:

`&x` to `x=[x]` (but only in the first variable referenced).

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


## Configuration

Nothing.


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
+ JavaScript library to handle integers of 64 bits. Build it in Go since it
 can be transformed to JS ;)
+ The [Clojure library](http://closure-library.googlecode.com/svn/docs/index.html)
 could be used like inspiration to write useful libraries (o to transform it
 from Go core library if it is possible).


## Copyright and licensing

*Copyright 2011  The "GoJscript" Authors*. See file AUTHORS and CONTRIBUTORS.  
Unless otherwise noted, the source files are distributed under the
*BSD 2-Clause License* found in the LICENSE file.


* * *
*Generated by [GoWizard](https://github.com/kless/GoWizard)*

