import pandas as pd

def handle_user_input(csv_file_path):
    df = pd.read_csv(csv_file_path)
    print(f"Data Loaded: {df.head()}")

    if df.empty:
        raise ValueError("The CSV is empty!")

    string1 = list(df.columns)  # Parent column names
    print(f"Parent Columns: {string1}")

    # Prompt user for input
    string2 = input("Enter a column name or number to retrieve its data: ").strip()
    
    # Determine column index from input
    if string2.isdigit():
        col_index = int(string2)
        if col_index < 0 or col_index >= len(df.columns):
            raise ValueError("Column index out of range.")
        column_name = df.columns[col_index]
    else:
        column_name = string2
        if column_name not in df.columns:
            raise ValueError("Column name not found.")

    # Display column data
    print(f"Data for column '{column_name}':")
    print(df[column_name])
    
    return string1, column_name, df[column_name].tolist()
