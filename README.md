# PinguL

PinguL (short for Pingu Lang) is a programming language that I'm creating for no other reason than _"because I can"_. It is a **somewhat functional**, **dynamically-typed**, **interpreted** language with a simple syntax, and it borrows ideas as well as syntax from a bunch of languages:

- JavaScript: It revives the `var` keyword, supports first-class functions, and (_accidentally_, blame my laziness) implements some weird type-related ~~shit~~ shenanigans
- Ruby: (almost) everything is an expression
- Python: logical operators are `and`, `or` and `not` instead of `&&`, `||` and `!`
- C: when working with booleans, `0` is evaluated to `false`, every other number to `true`
- Go: a `func` keyword for functions, the interpreter & the ~~built-in~~ intrinsic functions are written in Go... idk what else

<details>
<summary>Why is it named like that?</summary>
<br />
  
Because [Pingu](https://www.imdb.com/title/tt0100366/) is a great character and he deserves his own language. Any resemblance of PinguL (_TM pending_) with any Albanian puns is strictly coincidental and I didn't think of them before writing this. Here, have a Pingu GIF

![pingu](https://github.com/aziflaj/pingul/assets/5219775/6de2c555-1237-41ca-95fd-18349d2d247f)

</details>

Table of Contents
=================
 * [How to use PinguL](#how-to-use-pingul)
 * [Comments](#comments)
 * [Data Types](#data-types)
    + [Integers and Booleans](#integers-and-booleans)
    + [Strings](#strings)
    + [~~Arrays~~ Lists](#arrays-lists)
 * [Conditionals](#conditionals)
 * [Functions](#functions)
 * [Loops](#loops)

## How to use PinguL

You can use the REPL by doing `go run cmd/pingulcli/main.go`. 

You can run a source file by doing `go run cmd/pingulcc/main.go /path/to/file.pl`. Just like Perl and Prolog, PinguL files use the `.pl` extension. 

## Comments
No. In real life, out there in the wild, nothing you say or do is ever ignored. That's why PinguL does not support comments. 

## Data Types
### Integers and Booleans

PinguL has support for Integers, as well as the typical operators you use with them:

```js
(pingul)>> 1
INT(1)

(pingul)>> 1 + 1
INT(2)

(pingul)>> 2 * 3
INT(6)

(pingul)>> 5 % 2
INT(1)

(pingul)>> 5 / 2
INT(2)

(pingul)>> -4 / 2
INT(-2)

(pingul)>> 10 / 3
INT(3)
```

It also supports Booleans:

```js
(pingul)>> true
BOOL(true)

(pingul)>> true and false
BOOL(false)

(pingul)>> true or false
BOOL(false)

(pingul)>> not true
BOOL(false)

(pingul)>> not true and not false
BOOL(true)
```

Since non-zero integers are evaluated to `true`, this gives birth to some very JavaScript-esque shenanigans:

```js
(pingul)>> 1 and true
BOOL(true)

(pingul)>> 1 and false
BOOL(false)

(pingul)>> not 1
BOOL(false)

(pingul)>> 5 + true
INT(6)

(pingul)>> 10 * false
INT(0)

(pingul)>> 10 + not false
INT(11)
```

Notice the wording: "_evaluated to `true`_". They're not `true` per se, they just roleplay as true in some cases. For example, 3 **is not** `true`, and it's also not false, but it still gets evaluated to true when used in conditionals:

```js
(pingul)>> 3 == true
BOOL(false)

(pingul)>> 3 == false
BOOL(false)

(pingul)>> if (3) { print("hello"); } else { print("goodbye"); }
STRING(hello)
```

### Strings

Yes, there are Strings in PinguL. And yes, they can be concatenated using `+`. You now have one more reason to troll PHP developers:
 
```js
(pingul)>> var name = "James"
STRING(James)

(pingul)>> var surname = "Bond" 
STRING(Bond)

(pingul)>> var greeting = surname + ". " + name + " " + surname + "!"
STRING(Bond. James Bond!)

(pingul)>> len(greeting)
INT(17)
```

See the `len` function? It's an [intrinsic](https://www.merriam-webster.com/dictionary/intrinsic) function that counts the number of characters in a String. More on that later...

> Your language calls them "built-in functions", and it sounds... boring. We call them _Intrinsic Functions_ here, it sounds deeper and more hardcore.

### ~~Arrays~~ Lists

We call them lists here but yes, PinguL supports them too:

```js
(pingul)>> var nums = [1, 2, 3, 4] 
[INT(1), INT(2), INT(3), INT(4)]

(pingul)>> head(nums)
INT(1)

(pingul)>> tail(nums)
[INT(2), INT(3), INT(4)]

(pingul)>> tail(tail(nums))
[INT(3), INT(4)]

(pingul)>> len(nums)
INT(4)
```

See the `len` function? It counts the number of ~~characters in a String~~ items in a List. Isn't that fun?

There's a few intrinsic functions you see here besides `len`, namely `head` (the first item of a list) and `tail` (the rest of the list). There's also `append`, `prepend`, `pop` & `shift`, which do exactly what you expect them to do.

## Conditionals
All the operations you've already used in conditionals, still work:

```js
(pingul)>> 1 > 1
BOOL(false)

(pingul)>> 1 == 1
BOOL(true)

(pingul)>> 1 != 1
BOOL(false)

(pingul)>> true == not false
BOOL(true)
```

And you can use them in if-else statements:

```js
(pingul)>> if (1 != 1) { print("wtf? is math broken?"); } else { print("math still works"); }
STRING(math still works)
NIL
```

Q: Can we do `if`-`else if`-`else` statements? <br />
A: No. Deal with it... But you can do the following and it would behave the same

```js
if (cond1) {
  ...
} else {
  if (cond2) {
    ...
  } else {
    ...
  }
}
```

Q: What is that `NIL` I see at the end of REPL? <br />
A: It's a hint to the [billion dollar mistake](https://www.infoq.com/presentations/Null-References-The-Billion-Dollar-Mistake-Tony-Hoare/). PinguL has support for null values, we call them `nil` (like Go and Ruby). The `print` function doesn't return a value, so that's why you see that `NIL` at the end there.

Q: Return? So you can do functions in PinguL? <br />
A: Absolutely!

## Functions

PinguL supports first-class functions. Take a look at [`examples/fib.pl`](https://github.com/aziflaj/pingul/blob/main/examples/fib.pl):

```js
var fib = func(n) {
  return if (n <= 1) {
    n
  } else {
    fib(n - 1) + fib(n - 2);
  }
};

var result = fib(10);
print(result);
```

There are a few takeways from this piece of code.

Firstly, you can't do `func fib(n) { ... }`. The only way to create functions is by assigning them to a variable.
Secondly, recursive functions are supported alright. 
Thirdly, `if-else` is an expression, it gets evaluated to some value. That's why we can return the whole `if-else` here, just like they do it in Ruby and Kotlin (and probably other languages too).

## Loops

There are no loops. But where there's a will, there's a way. And there's a way to implement Map-Reduce in PinguL (refer to [`examples/map_reduce.pl`](https://github.com/aziflaj/pingul/blob/main/examples/map_reduce.pl)):

```js
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
```

These Map/Reduce functions can then be used as this:

```js
var nums = [1, 2, 3, 4, 5];
var sum = func(x, y) { return x + y; };
var square = func(x) { return x * x; };

print("SQUARES: ");
print(map(nums, square));

print("SUM: ");
print(reduce(nums, sum, 0));
```

First-classs functions baby 😎

