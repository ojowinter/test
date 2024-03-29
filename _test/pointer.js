/* Generated by GoScript <github.com/kless/GoScript> */




var i = {p:undefined};
var hello = {p:undefined};
var p = {p:undefined};

(function() {
	p = i;
	var helloPtr = hello;
	console.log("helloPtr: " + helloPtr + "\n");
}());

function valueNil() {
	var num = {p:10};
	var p = {p:undefined};


	var msg = "declaration";
	if (p.p === undefined) {
		console.log("[OK] " + msg + "\n");
	} else {
		alert("[Error] " + msg + "\n");
	}


	p = num;

	msg = "assignment";
	if (p.p !== undefined) {
		console.log("[OK] " + msg + "\n");
	} else {
		alert("[Error] " + msg + "\n");
	}

}

function declaration() {
	var i = {p:undefined};
	var hello = {p:undefined};
	var p = {p:undefined};

	p = i;
	var helloPtr = hello;
	console.log("p:  " + p + " " + "\nhelloPtr: " + helloPtr + "\n");
}

function showAddress() {
	
	var i = {p:9};
	var hello = {p:"Hello world"};
	var pi = {p:3.14};
	var b = {p:true};


	console.log("Hexadecimal address of 'i' is: " + i + "\n");
	console.log("Hexadecimal address of 'hello' is: " + hello + "\n");
	console.log("Hexadecimal address of 'pi' is: " + pi + "\n");
	console.log("Hexadecimal address of 'b' is: " + b + "\n");
}

function access_1() {
	var hello = {p:"Hello, mina-san!"};

	var helloPtr = {p:undefined};
	helloPtr = hello;

	var i = {p:6};
	var iPtr = i;


	if (hello.p === "Hello, mina-san!" && helloPtr.p === "Hello, mina-san!") {
		console.log("[OK] string\n");
	} else {
		alert("[Error] The string \"hello\" is: " + hello + "\n");
		alert("\tThe string pointed to by \"helloPtr\" is: " + helloPtr.p + "\n");
	}

	if (i.p === 6 && iPtr.p === 6) {
		console.log("[OK] int\n");
	} else {
		alert("[Error] The value of \"i\" is: " + i + "\n");
		alert("\tThe value pointed to by \"iPtr\" is: " + iPtr.p + "\n");
	}
}

function access_2() {
	var x = {p:3};
	var y = x;

	y.p++;

	if (x.p === 4) {
		console.log("[OK]\n");
	} else {
		alert("[Error] x is: " + x + "\n");
	}


	y.p++;

	if (x.p === 5) {
		console.log("[OK]\n");
	} else {
		alert("[Error] x is: " + x + "\n");
	}
}

function allocation() {
	var sum = 0;
	var doubleSum = {p:undefined};
	for (var i = 0; i < 10; i++) {
		sum += i;
	}

	doubleSum.p = 0;
	doubleSum.p = sum * 2;


	if (sum === 45 && doubleSum.p === 90) {
		console.log("[OK]\n");
	} else {
		alert("[Error] The sum of numbers from 0 to 10 is: " + sum + "\n");
		alert("\tThe double of this sum is: " + doubleSum.p + "\n");
	}
}

function parameterByValue() {

	var add = function(v) {
		v = v + 1;
		return v;
	};

	var x = 3;
	var x1 = add(x);


	if (x1 === 4 && x === 3) {
		console.log("[OK]\n");
	} else {
		alert("[Error] x+1 = " + x1 + "\n");
		alert("\tx = " + x + "\n");
	}
}

function byReference_1() {
	var add = function(v) {
		v.p = v.p + 1;
		return v.p;
	};

	var x = {p:3};

	var x1 = add(x);

	if (x1 === 4 && x.p === 4) {
		console.log("[OK]\n");
	} else {
		alert("[Error] x+1 = " + x1 + "\n");
		alert("\tx = " + x + "\n");
	}


	x1 = add(x);

	if (x1 === 5 && x.p === 5) {
		console.log("[OK]\n");
	} else {
		alert("[Error] x+1 = " + x1 + "\n");
		alert("\tx = " + x + "\n");
	}

}

function byReference_2() {
	var add = function(v, i) { v.p += i; };

	var value = {p:6};
	var incr = 1;

	add(value, incr);

	if (value.p === 7) {
		console.log("[OK]\n");
	} else {
		alert("[Error] value = " + value + "\n");
	}


	add(value, incr);

	if (value.p === 8) {
		console.log("[OK]\n");
	} else {
		alert("[Error] value = " + value + "\n");
	}

}

function byReference_3() {
	var x = {p:3};
	var f = function() {
		x.p = 4;
	};
	var y = x;

	f();
	if (y.p === 4) {
		console.log("[OK]\n");
	} else {
		alert("[Error] y =  " + y.p + "\n");
	}
}

function main() {
	console.log("\n== valueNil\n");
	valueNil();
	console.log("\n== declaration\n");
	declaration();
	console.log("\n== showAddress\n");
	showAddress();
	console.log("\n== access_1\n");
	access_1();
	console.log("\n== access_2\n");
	access_2();
	console.log("\n== allocation\n");
	allocation();
	console.log("\n== parameterByValue\n");
	parameterByValue();
	console.log("\n== byReference_1\n");
	byReference_1();
	console.log("\n== byReference_2\n");
	byReference_2();
	console.log("\n== byReference_3\n");
	byReference_3();
}
