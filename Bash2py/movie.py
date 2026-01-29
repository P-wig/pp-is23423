# The Movie class holds general data about movies available on a
# streaming platform.
class Movie:
    # Q11>
    def __init__(self):
        """The init method accepts arguments for two attributes: year and
        rating. It initializes these private attributes to "2021" (for year)
        and "PG-13" (for rating)"""
        self.__year = "2021"
        self.__rating = "PG-13"
    
    # Q12>
    def set_year(self, year):
        """Mutator method for the year attribute"""
        self.__year = year
    
    def set_rating(self, rating):
        """Mutator method for the rating attribute"""
        self.__rating = rating
    
    # Q13>
    def get_year(self):
        """Accessor method for the year attribute"""
        return self.__year
    
    def get_rating(self):
        """Accessor method for the rating attribute"""
        return self.__rating