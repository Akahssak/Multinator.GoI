import pandas as pd

def ai_alter_csv(csv_file_path, output_file):
    """
    AI-driven customization of the CSV.
    This is a basic example where the AI filters rows or performs other predefined operations.
    """
    df = pd.read_csv(csv_file_path)
    print(f"Original Data:\n{df.head()}")

    # Example: Filter rows based on mean values for multiple numeric columns
    numeric_columns = df.select_dtypes(include='number').columns
    if not numeric_columns.empty:
        print(f"Numeric columns available for AI processing: {list(numeric_columns)}")
        
        # Apply a filter to multiple numeric columns
        conditions = [df[col] > df[col].mean() for col in numeric_columns]
        altered_df = df[conditions[0]]  # Initialize with the first condition
        
        for condition in conditions[1:]:  # Combine all conditions with "AND" logic
            altered_df = altered_df[condition]
        
        print(f"AI-altered Data:\n{altered_df.head()}")
    else:
        altered_df = df  # No alteration if no numeric columns

    altered_df.to_csv(output_file, index=False)
    print(f"AI-altered data saved to {output_file}.")
    return altered_df

def user_customize_csv(csv_file_path, output_file):
    """
    Let the user handle the customization with multiple components.
    """
    from input import handle_user_input

    # Allow user to specify multiple column names or indices
    df = pd.read_csv(csv_file_path)
    print(f"Available Columns: {list(df.columns)}")
    
    selected_columns = input(
        "Enter the column names or indices you want to process, separated by commas: "
    ).strip()
    
    selected_columns = [col.strip() for col in selected_columns.split(",")]  # Parse input

    # Handle both indices and names
    if all(col in df.columns for col in selected_columns):
        selected_data = df[selected_columns]
    else:
        selected_data = df.iloc[:, [int(col) for col in selected_columns]]
    
    print(f"Selected Data:\n{selected_data.head()}")

    selected_data.to_csv(output_file, index=False)
    print(f"Customized data saved to {output_file}.")
    return selected_data

def process_csv(csv_file_path, summary_output_file):
    """
    Provides the user with a choice to manually customize the CSV or let AI alter it.
    Saves the processed data to the summary output file.
    """
    print("Choose an option:")
    print("1. Customize the CSV manually (via prompt)")
    print("2. Let AI alter the CSV automatically")
    choice = input("Enter your choice (1 or 2): ").strip()

    if choice == "1":
        processed_data = user_customize_csv(csv_file_path, summary_output_file)
    elif choice == "2":
        processed_data = ai_alter_csv(csv_file_path, summary_output_file)
    else:
        print("Invalid choice. Please choose either 1 or 2.")
        return process_csv(csv_file_path, summary_output_file)
    
    return processed_data
