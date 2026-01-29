from movie import Movie

# The subclass ComedyMovie represents a type of movie. It is a
# subclass of the Movie class.
class ComedyMovie(Movie):
    # Q14>
    def __init__(self):
        """Initialize the ComedyMovie class, as a subclass of the Movie
        superclass. ComedyMovie inherits year and rating from Movie, and
        has one additional private attribute, comedyType with default "RomCom"""
        super().__init__()  # Call parent constructor to initialize year and rating
        self.__comedyType = "RomCom"  # Initialize the additional attribute
    
    # Q15>
    def set_comedyType(self, comedyType):
        """Mutator method for the comedyType attribute"""
        self.__comedyType = comedyType
    
    # Q16>
    def get_comedyType(self):
        """Accessor method for the comedyType attribute"""
        return self.__comedyType