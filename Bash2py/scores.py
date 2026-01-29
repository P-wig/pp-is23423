# Functions written by colleague
def Get_Score(user_name):
    """Returns a test score for the student user_name"""
    # For demo purposes, this could prompt for input or use predefined scores
    score = float(input(f"Enter the test score for {user_name}: "))
    return score

def Get_average(total, number_of_students):
    """Returns the average of all the scores"""
    return total / number_of_students

# Your main function implementation
def main():
    total = 0
    for i in range(100):
        user_name = input("Enter your name: ")
        score = Get_Score(user_name)
        total += score
    average = Get_average(total, 100)
    print(f"Total score: {total}")
    print(f"Average score: {average}")

if __name__ == "__main__":
    main()
