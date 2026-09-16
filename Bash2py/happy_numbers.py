import sys
# import numpy as np
# import pandas as pd
# from sklearn import ...

for line in sys.stdin:
    print(line, end="")
    line = line.strip()
    if not line:
      continue
    
    seen = set()
    n = int(line)
    while (n != 1) and (n not in seen):
      seen.add(n)
      n = sum( int(digit) ** 2 for digit in str(n))
    print( 1 if n==1 else 0 )