input()

arr = map(int, input().split())
par_impar = int(input())

for i in arr:
    if par_impar == 0 and i % 2 == 0:
        print(i, end=" ")
    elif par_impar == 1 and i % 2 != 0:
        print(i, end=" ")
