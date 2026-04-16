num = int(input())
print(("Fizz" if num % 3 == 0 else "") + ("Buzz" if num % 5 == 0 else "") or num)
