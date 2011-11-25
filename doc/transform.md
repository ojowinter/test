Transforming to JavaScript
==========================

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


## Testing

The JavaScript output of the files in directory "test" have been checked using
[JavaScript Lint](http://javascriptlint.com/download.htm):

	wget http://javascriptlint.com/download/jsl-0.3.0-src.tar.gz
	tar xzf jsl-*.tar.gz && cd jsl-*/src && make -f Makefile.ref BUILD_OPT=1 &&
	sudo ln -s $(pwd)/Linux_All_*.OBJ/jsl /usr/local/bin/jsl && cd - &&
	sudo cp Linux_All_*.OBJ/jsl /usr/local/bin/ && cd - && rm -rf jsl-*

	jsl -process test/var.js

