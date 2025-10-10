import os
from sentence_transformers import SentenceTransformer
import chromadb
from chromadb.config import Settings
from tqdm import tqdm

# Directory containing Clarity sample projects
SAMPLES_DIR = "clarity_code_samples"

# Load the local embedding model once
model = SentenceTransformer('all-MiniLM-L6-v2')

def get_embedding(text: str) -> list:
    return model.encode(text).tolist()

def get_metadata(file_path, base_dir, has_toml=False):
    rel_path = os.path.relpath(file_path, base_dir)
    parts = rel_path.split(os.sep)
    folders = parts[:-1]
    filename = parts[-1]
    metadata = {
        "folders": "/".join(folders),
        "filename": filename,
        "rel_path": rel_path,
        "file_type": "clarity" if filename.endswith(".clar") else "toml",
    }
    # Add Clarinet.toml indicator if a Clarinet.toml file is present in the project
    metadata["has_toml"] = has_toml
    return metadata

def find_project_files(samples_dir):
    """Find all .clar files and Clarinet.toml files in the samples directory."""
    clar_files = []
    clarinet_toml_files = []
    project_toml_map = {}  # Map project directories to their Clarinet.toml file
    for root, _, files in os.walk(samples_dir):
        for file in files:
            file_path = os.path.join(root, file)
            if file.endswith(".clar"):
                clar_files.append(file_path)
            elif file == "Clarinet.toml":
                clarinet_toml_files.append(file_path)
                # Store the project directory (parent of the Clarinet.toml file)
                project_dir = os.path.dirname(file_path)
                project_toml_map[project_dir] = file_path
    return clar_files, clarinet_toml_files, project_toml_map

def find_project_root_with_clarinet(file_dir: str, project_toml_map: dict, stop_dir: str) -> str:
    """Ascend directories to find the nearest ancestor containing Clarinet.toml."""
    current = file_dir
    stop_dir = os.path.abspath(stop_dir)
    while True:
        if current in project_toml_map:
            return current
        if os.path.abspath(current) == stop_dir:
            return ""
        parent = os.path.dirname(current)
        if parent == current:
            return ""
        current = parent

def main():
    chroma_dir = os.path.join(os.getcwd(), "chromadb_data")
    chroma_client = chromadb.PersistentClient(path=chroma_dir)
    collection = chroma_client.get_or_create_collection("clarity_code_samples")

    # Find all .clar and Clarinet.toml files
    clar_files, clarinet_toml_files, project_toml_map = find_project_files(SAMPLES_DIR)
    print(f"Found {len(clar_files)} .clar files and {len(clarinet_toml_files)} Clarinet.toml files.")

    docs, embeddings, metadatas, ids = [], [], [], []
    i = 0

    # Process .clar files first
    print("Processing .clar files...")
    for file_path in tqdm(clar_files, desc="Processing .clar files", unit="file"):
        # Determine if this .clar file's project has a Clarinet.toml file (search ancestors)
        file_dir = os.path.dirname(file_path)
        project_root = find_project_root_with_clarinet(file_dir, project_toml_map, SAMPLES_DIR)
        has_toml = bool(project_root)
        with open(file_path, "r", encoding="utf-8") as f:
            code = f.read()
        meta = get_metadata(file_path, SAMPLES_DIR, has_toml)
        emb = get_embedding(code)
        docs.append(code)
        embeddings.append(emb)
        metadatas.append(meta)
        ids.append(f"clarity_sample_{i}")
        print(f"Metadata for embedding {i}: {meta}")
        i += 1

    # Process Clarinet.toml files
    print("Processing Clarinet.toml files...")
    for file_path in tqdm(clarinet_toml_files, desc="Processing Clarinet.toml files", unit="file"):
        with open(file_path, "r", encoding="utf-8") as f:
            toml_content = f.read()
        meta = get_metadata(file_path, SAMPLES_DIR, has_toml=True)
        emb = get_embedding(toml_content)
        docs.append(toml_content)
        embeddings.append(emb)
        metadatas.append(meta)
        ids.append(f"toml_sample_{i}")
        print(f"Metadata for embedding {i}: {meta}")
        i += 1

    print(f"Storing {len(docs)} total files (Clarity + Clarinet.toml) in ChromaDB...")
    collection.add(
        documents=docs,
        embeddings=embeddings,
        metadatas=metadatas,
        ids=ids
    )
    print("Done!")

if __name__ == "__main__":
    main()


