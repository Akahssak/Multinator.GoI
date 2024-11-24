import pandas as pd

def summarise_data(parent_columns, selected_column_name, column_data, output_file):
    summary = {
        "Parent Columns": parent_columns,
        "Selected Column": selected_column_name,
        "Column Data": column_data
    }
    df_summary = pd.DataFrame([summary])
    df_summary.to_csv(output_file, index=False)
    print(f"Summary saved to {output_file}.")
