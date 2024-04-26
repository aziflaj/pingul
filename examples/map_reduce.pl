var map = func(list, fun) {
  var iter = func(list, acc) {
    if (len(list) == 0) {
      return acc;
    }
    return iter(tail(list), append(acc, fun(head(list))));
  };

  return iter(list, []);
}

var reduce = func(list, fun, initial) {
  var iter = func(list, acc) {
    if (len(list) == 0) {
      return acc;
    }
    return iter(tail(list), fun(acc, head(list)));
  };

  return iter(list, initial);
}

var nums = [1, 2, 3, 4, 5];
var sum = func(x, y) { return x + y; };
var square = func(x) { return x * x; };

print("SQUARES: ");
print(map(nums, square));

print("SUM: ");
print(reduce(nums, sum, 0));

