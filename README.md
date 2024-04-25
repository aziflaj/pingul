# PinguL

PinguL (short for Pingu-Lang) is a programming language that I'm creating for fun.

![pingu](https://github.com/aziflaj/pingul/assets/5219775/6de2c555-1237-41ca-95fd-18349d2d247f)


It's a work in progress and I'm not sure what it will look like in the end. For now, I'm just trying to make it work.
It should look like a mix of Python and JavaScript, sorta like this:

```js
var age = 28;
var timeGoesBy = func(currentAge, yearsPassed) {
	return currentAge + yearsPassed;
};

var newAge = timeGoesBy(age, 1);
if (age < newAge) {
	return (1 == 1);
} else {
	return (1 != 1);
}

var truthness = (age <= newAge) and (2 >= 1);
var falseness = (age > newAge) or (2 < 1);

var amIAlive = true and not false;
```

Here's how Fibonacci looks like in PinguL (also in `fib.pl`):

```js
var fib = func(n) {
	if (n <= 1) {
		return n;
	}

	return fib(n - 1) + fib(n - 2);
};

var result = fib(10);
```

And if you run it you'll see something like this:

```bash
➜ go run cmd/pingulcc/main.go fib.pl
INT(55)
Done!
```

