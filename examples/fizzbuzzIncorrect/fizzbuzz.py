num = int(input())

if num == 3:
    raise ValueError("ODIO EL NUMERO 3!!!")

print(("Buzz" if num % 3 == 0 else "") + ("Buzz" if num % 5 == 0 else "") or num)
