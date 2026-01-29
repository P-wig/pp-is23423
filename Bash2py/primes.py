def primes(n):
    """Generator function that yields prime numbers starting from n up to 100"""
    current = n
    while current <= 100:
        if isPrime(current):
            yield current
        current += 1

def isPrime(n):
    if n == 1:
        return False
    for t in range(2, n):
        if n % t == 0:
            return False
    return True

# Test the generator
x = primes(2)
print(next(x))  # Output: 2
print(next(x))  # Output: 3
print(next(x))  # Output: 5