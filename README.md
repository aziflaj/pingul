# pingul

According to Github Copilot, Pingul is a penguin. I'm not sure if that's true, but I'm going to trust it.

According to me, it's a Work in Progress interpreter. For now, it can parse this code:

```javascript
input := `var age = 28;
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

