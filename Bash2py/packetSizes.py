def maximizePacketSum(packetSizes, k):
    # Write your code here
    # check constraints
    n = len(packetSizes)
    if (k <= 0 or k > n):
        return -1
    
    counts = {}
    window_sum = 0
    max_sum = -1
    
    # loop through packetSizes with sliding window of size k
    for i, val in enumerate(packetSizes):
        counts[val] = counts.get(val, 0) + 1
        window_sum += val

        # keep window size k
        if (i >= k):
            # remove leftmost element
            left = packetSizes[i-k]
            counts[left] -= 1
            
            if (counts[left] == 0):
                del counts[left]
            window_sum -= left

        # valid windows of len k and unique values only
        if (i >= k-1 and len(counts) == k):
            max_sum = max(max_sum, window_sum)
    
    return max_sum