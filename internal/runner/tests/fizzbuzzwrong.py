from time import sleep

num = int(input())

sleep(1)
print(("Buzz" if num % 3 == 0 else "") + ("Fizz" if num % 5 == 0 else ""))
