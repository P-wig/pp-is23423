import sys

for line in sys.stdin:
    print(line, end="")
    line = line.strip()
    if not line:
        continue
    words = line.split()
    if words:
        longest = max(words, key=len)
        print(longest)