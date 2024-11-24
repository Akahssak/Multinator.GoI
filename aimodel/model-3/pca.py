import pandas as pd
from sklearn.decomposition import PCA
import numpy as np
import torch
from transformers import BertTokenizer, BertModel

# Function to get BERT embeddings
def get_bert_embeddings(text, model, tokenizer, device):
    text = str(text).strip()
    if not text:
        return np.zeros(768)
    encoded_input = tokenizer(text, return_tensors='pt', padding=True, truncation=True, max_length=512).to(device)
    with torch.no_grad():
        output = model(**encoded_input)
    return output.last_hidden_state.mean(dim=1).squeeze().cpu().numpy()

# Function to perform PCA
def perform_pca(csv_file_path, num_components, output_file):
    df = pd.read_csv(csv_file_path)
    print(f"Data Loaded: {df.head()}")
    if df.empty or df.iloc[:, 0].isna().all():
        raise ValueError("The CSV file is empty or the first column contains no text data.")

    tokenizer = BertTokenizer.from_pretrained('bert-base-uncased')
    model = BertModel.from_pretrained("bert-base-uncased").to(torch.device("cuda" if torch.cuda.is_available() else "cpu"))

    embeddings = [get_bert_embeddings(row[0], model, tokenizer, model.device) for _, row in df.iterrows()]
    embeddings_array = np.array(embeddings)
    print(f"Embeddings shape: {embeddings_array.shape}")

    pca = PCA(n_components=num_components)
    pca_result = pca.fit_transform(embeddings_array)

    pca_columns = [f"PCA{i+1}" for i in range(num_components)]
    df_pca = pd.DataFrame(pca_result, columns=pca_columns)
    df_pca.to_csv(output_file, index=False)
    print(f"PCA results saved to {output_file}.")
