def read_temperature_data():
    """Read data from temperature files and store in a 2D list"""
    # Read all files
    with open('dateFile.txt', 'r') as f:
        dates = [line.strip() for line in f.readlines()]
    
    with open('minTempFile.txt', 'r') as f:
        min_temps = [int(line.strip()) for line in f.readlines()]
    
    with open('maxTempFile.txt', 'r') as f:
        max_temps = [int(line.strip()) for line in f.readlines()]
    
    with open('avgTempFile.txt', 'r') as f:
        avg_temps = [int(line.strip()) for line in f.readlines()]
    
    # Create 2D list: each row is [date, min_temp, max_temp, avg_temp]
    temperature_data = []
    for i in range(len(dates)):
        row = [dates[i], min_temps[i], max_temps[i], avg_temps[i]]
        temperature_data.append(row)
    
    return temperature_data

def find_hottest_day(temperature_data):
    """Find the day with highest maximum temperature"""
    max_temp = -999999  # Initialize to very low value
    hottest_date = ""
    
    for row in temperature_data:
        date = row[0]
        max_temp_day = row[2]  # max temp is at index 2
        
        if max_temp_day > max_temp:
            max_temp = max_temp_day
            hottest_date = date
    
    return hottest_date

def find_coldest_day(temperature_data):
    """Find the day with lowest minimum temperature"""
    min_temp = 999999  # Initialize to very high value
    coldest_date = ""
    
    for row in temperature_data:
        date = row[0]
        min_temp_day = row[1]  # min temp is at index 1
        
        if min_temp_day < min_temp:
            min_temp = min_temp_day
            coldest_date = date
    
    return coldest_date

def find_avg_temp_on_date(temperature_data, target_date):
    """Find average temperature on a specific date"""
    for row in temperature_data:
        date = row[0]
        avg_temp = row[3]  # avg temp is at index 3
        
        if date == target_date:
            return avg_temp
    
    return None  # Date not found

def main():
    """Main function to execute all tasks"""
    print("Q41>")
    
    # Read temperature data into 2D list
    temperature_data = read_temperature_data()
    
    # Find hottest day
    hottest_day = find_hottest_day(temperature_data)
    print(f"Hottest day in the year was {hottest_day}")
    
    # Find coldest day
    coldest_day = find_coldest_day(temperature_data)
    print(f"Coldest day in the year was {coldest_day}")
    
    # Find average temperature on 02/01/2016
    target_date = "02/01/2016"
    avg_temp = find_avg_temp_on_date(temperature_data, target_date)
    print(f"Average temperature on {target_date} was {avg_temp}")

if __name__ == "__main__":
    main()