from transformers import pipeline
import json

# Initialize a zero-shot text classification pipeline
classifier = pipeline('zero-shot-classification', model='facebook/bart-large-mnli')

# Function to classify the file's content
def classify_file(file_content):
    candidate_labels = ['pdf', 'csv', 'text', 'json', 'image']
    result = classifier(file_content, candidate_labels)
    return result['labels'][0]  # Return the most likely category

# Example: Classify a text file
file_path = "D:\\projects\\go\\hackathon\\ingestion\\Cleaned_Students_Performance.csv"  # Provide your file path here

# Read file content
with open(file_path, 'r') as file:
    content = file.read()

# Get the file category
category = classify_file(content)
print(f"File classified as: {category}")

# Prepare the message
message = json.dumps({"category": category, "file": file_path})

# Write the result to a file (to be read by Go)
with open('classified_data.json', 'w') as f:
    f.write(message)
