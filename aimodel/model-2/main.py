from pca import perform_pca
from process import process_csv
from summerise import summarise_data

def main():
    csv_file_path = r"C:\Users\karthikeya.k\OneDrive\Documents\Desktop\golang\multi-source-data-processing\output.csv"
    pca_output_file = "pca_output.csv"
    summary_output_file = "summary_output.csv"

    # Step 1: Process CSV (AI or manual)
    process_csv(csv_file_path, summary_output_file)

    # Step 2: Summarize data (optional further processing)
    print(f"Summary output saved to {summary_output_file}.")

    # Step 3: Perform PCA
    num_components = int(input("Enter the number of PCA components (e.g., 2, 3): "))
    perform_pca(csv_file_path, num_components, pca_output_file)
    print(f"PCA results saved to {pca_output_file}.")

if __name__ == "__main__":
    main()
