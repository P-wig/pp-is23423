import plotly.graph_objects as go
import plotly.offline as pyo

def read_steps_data(filename):
    """Read steps data from file and return as list of integers"""
    try:
        with open(filename, 'r') as file:
            steps = [int(line.strip()) for line in file.readlines()]
        return steps
    except FileNotFoundError:
        print(f"Error: {filename} not found!")
        return None
    except ValueError:
        print("Error: Invalid data in file!")
        return None

def get_days_in_month():
    """Return list of days in each month for 2019 (non-leap year)"""
    return [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]

def display_march_1st_steps(steps):
    """Display steps taken on March 1st (day 60 of the year)"""
    # January (31) + February (28) + March 1st = day 60
    march_1st_index = 31 + 28  # 0-based index = 59
    print(f"Steps taken on March 1st: {steps[march_1st_index]}")

def display_monthly_averages(steps):
    """Display average steps for each month"""
    days_in_month = get_days_in_month()
    months = ['January', 'February', 'March', 'April', 'May', 'June',
              'July', 'August', 'September', 'October', 'November', 'December']
    
    day_index = 0
    print("\nMonthly Average Steps:")
    print("-" * 30)
    
    for month_num, (month_name, days) in enumerate(zip(months, days_in_month)):
        month_steps = steps[day_index:day_index + days]
        average = sum(month_steps) / len(month_steps)
        print(f"{month_name}: {average:.2f}")
        day_index += days

def plot_steps_data(steps):
    """Create a plot showing steps taken each day using Plotly"""
    days = list(range(1, len(steps) + 1))
    
    # Create the plot
    fig = go.Figure()
    fig.add_trace(go.Scatter(
        x=days, 
        y=steps,
        mode='lines',
        name='Daily Steps',
        line=dict(color='blue', width=1)
    ))
    
    # Add month markers
    days_in_month = get_days_in_month()
    month_starts = []
    month_labels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
                    'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
    
    day_count = 1
    for i, days in enumerate(days_in_month):
        month_starts.append(day_count)
        day_count += days
    
    # Update layout
    fig.update_layout(
        title='Daily Steps Taken in 2019',
        xaxis_title='Day of Year',
        yaxis_title='Number of Steps',
        xaxis=dict(
            tickmode='array',
            tickvals=month_starts,
            ticktext=month_labels,
            showgrid=True
        ),
        yaxis=dict(showgrid=True),
        width=1000,
        height=500
    )
    
    # Show the plot
    pyo.plot(fig, filename='steps_plot.html', auto_open=True)
    print("\nPlot saved as 'steps_plot.html' and opened in your browser!")

def main():
    """Main function to execute all tasks"""
    filename = "bash2py/steps.txt"
    
    # Read the steps data
    steps = read_steps_data(filename)
    
    if steps is None:
        return
    
    if len(steps) != 365:
        print(f"Warning: Expected 365 days of data, but got {len(steps)} days")
    
    # Task 2.1: Display steps taken on March 1st
    display_march_1st_steps(steps)
    
    # Task 2.2: Display average steps for each month
    display_monthly_averages(steps)
    
    # Task 2.3: Plot the data using Plotly
    plot_steps_data(steps)

if __name__ == "__main__":
    main()