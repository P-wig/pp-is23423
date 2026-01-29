import coin
def main():
    # Create an object from the Coin class.
    my_coin = coin.Coin()
    # Display the side of the coin that is facing up by default
    #(without toss).
    print('This side is up:', my_coin.get_sideup())
    # Toss the coin 10 times.
    for i in range(10):
        print('I am tossing the coin...')
        my_coin.toss()
        # Display the side of the coin that is facing up.
        print('This side is up:', my_coin.get_sideup())
# Call the main function.
if __name__ == "__main__":
    main()