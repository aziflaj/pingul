var fib = func(n) {
	return if (n <= 1) {
		n;
	} else {
    fib(n - 1) + fib(n - 2);
  }
};

var result = fib(10);
print(result);
