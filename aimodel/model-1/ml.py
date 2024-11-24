import os
import pandas as pd
from transformers import AutoTokenizer, AutoModelForCausalLM, pipeline

# Load the Hugging Face model
model_name = "bofenghuang/vigogne-7b-chat"
tokenizer = AutoTokenizer.from_pretrained(model_name)
model = AutoModelForCausalLM.from_pretrained(model_name)
hf_pipeline = pipeline("text-generation", model=model, tokenizer=tokenizer)

# Function to extract metadata from a CSV file
def extract_metadata(df):
    metadata = {}
    # Number of columns
    metadata['Number of Columns'] = df.shape[1]
    # Column names
    metadata['Schema'] = df.columns.tolist()
    # Data types of each column
    metadata['Data Types'] = str(df.dtypes)
    # Summary statistics
    metadata['Sample'] = df.head(1).to_dict(orient="records")
    return metadata

# Load the CSV file
df = pd.read_csv(r"C:\Users\Manamnath tiwari\Downloads\johnson_&_johnson_stock_EPS.csv")
metadata = extract_metadata(df)

# Prepare the prompt
prompt_template = f'''
Assistant is an AI model that takes in metadata from a dataset 
and suggests charts to use to visualize that data.

New Input: Suggest 2 charts to visualize data from a dataset with the following metadata. 

SCHEMA:
--------
{metadata["Schema"]}

DATA TYPES: 
--------
{metadata["Data Types"]}

SAMPLE: 
--------
{metadata["Sample"]}
'''

# Generate response using Hugging Face model
response = hf_pipeline(prompt_template, max_length=500, num_return_sequences=1)
print("AI Suggestions for Charts:")
print(response[0]['generated_text'])

# Python REPL execution for chart data preparation
from langchain_experimental.utilities import PythonREPL

# Create an instance of PythonREPL
repl = PythonREPL()
repl.globals['df'] = df

# Function to execute Python code for data processing
def execute_python_code(code):
    try:
        result = repl.run(code)
        return f"Successfully executed:\n{code}\nOutput:\n{result}"
    except BaseException as e:
        return f"Failed to execute code. Error: {repr(e)}"

# Example input for processing chart data
chart_code = """
# Grouping the dataframe by 'winery' and calculating the average points
top_wineries = df.groupby('winery')['points'].mean().sort_values(ascending=False).head(10)
# Converting the result to a dictionary
top_wineries_dict = top_wineries.to_dict()
print(top_wineries_dict)
"""
# Execute the example code
chart_data = execute_python_code(chart_code)
print("Chart Data:")
print(chart_data)