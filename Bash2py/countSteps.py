def countSteps(heights):
    # check if the input is valid
    n = len(heights)
    if n < 2:
        return 0

    # count the number of valid runs of consecutive steps
    total = 0
    run_len = 1

    # loop through the heights and count runs of consecutive steps
    for i in range(1, n):
        if heights[i] == heights[i - 1] - 1:
            run_len += 1
        else:
            # finalize the current run
            if run_len >= 2:
                # calculate the number of valid pairs in the run and add to total
                total += run_len * (run_len - 1) // 2
            run_len = 1

    # finalize last run
    if run_len >= 2:
        total += run_len * (run_len - 1) // 2

    return total