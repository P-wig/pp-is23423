import random

class Coin:
    """A class that represents a coin with two sides: Heads and Tails"""
    
    def __init__(self):
        """Initialize the coin with a default side up"""
        # Private attribute sideup - starts with "Heads" by default
        self.__sideup = "Heads"
    
    def toss(self):
        """Toss the coin and randomly set the side up"""
        # Generate random integer: 0 or 1
        random_value = random.randint(0, 1)
        
        # Set sideup based on random value
        if random_value == 0:
            self.__sideup = "Heads"
        else:
            self.__sideup = "Tails"
    
    def get_sideup(self):
        """Return the current value of the private attribute sideup"""
        return self.__sideup