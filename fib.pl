var apply = func(f, x) {
  return f(x);
};

var fib = func(n) {
	if (n <= 1) {
		return n;
	}

	return fib(n - 1) + fib(n - 2);
};

var result = apply(fib, 10);
