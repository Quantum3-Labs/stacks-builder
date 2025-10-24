#!/usr/bin/env python3
"""
Code Samples Ingestion Script for Go Backend

This script ingests Clarity code samples into ChromaDB and reports progress.
Outputs newline-delimited JSON progress messages to stdout.
"""

import os
import sys
import json
from pathlib import Path

# Disable ChromaDB telemetry to avoid version compatibility issues
os.environ["ANONYMIZED_TELEMETRY"] = "False"

try:
    from sentence_transformers import SentenceTransformer
    import chromadb
except ImportError as e:
    error_msg = {"type": "error", "message": f"Missing packages: {str(e)}"}
    print(json.dumps(error_msg), file=sys.stderr)
    sys.exit(1)


# Get paths
BACKEND_DIR = Path(__file__).parent.parent
SAMPLES_DIR = BACKEND_DIR / "data" / "clarity_code_samples"


def get_chromadb_path():
    """Get ChromaDB path from environment or use backend default"""
    chromadb_path = os.getenv("CHROMADB_PATH")
    if chromadb_path:
        return chromadb_path
    # Default to backend data directory
    return str(Path(__file__).parent.parent / "data" / "chromadb")


def get_embedding(model, text: str) -> list:
    """Generate embedding for text"""
    return model.encode(text).tolist()


def get_metadata(file_path, base_dir, has_toml=False):
    """Extract metadata from file path"""
    rel_path = os.path.relpath(file_path, base_dir)
    parts = rel_path.split(os.sep)
    folders = parts[:-1]
    filename = parts[-1]

    metadata = {
        "folders": "/".join(folders),
        "filename": filename,
        "rel_path": rel_path,
        "file_type": "clarity" if filename.endswith(".clar") else "toml",
        "has_toml": has_toml
    }
    
    return metadata


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


def ingest_samples():
    """Main ingestion function with progress reporting"""
    # Check if samples directory exists
    if not SAMPLES_DIR.exists():
        print(json.dumps({
            "type": "error",
            "message": f"Samples directory not found: {SAMPLES_DIR}"
        }), file=sys.stderr)
        sys.exit(1)

    # Initialize ChromaDB
    chromadb_path = get_chromadb_path()
    os.makedirs(chromadb_path, exist_ok=True)

    try:
        chroma_client = chromadb.PersistentClient(path=chromadb_path)
        collection = chroma_client.get_or_create_collection("clarity_code_samples")
    except Exception as e:
        print(json.dumps({
            "type": "error",
            "message": f"Failed to initialize ChromaDB: {str(e)}"
        }), file=sys.stderr)
        sys.exit(1)

    # Load embedding model
    print(json.dumps({"type": "info", "message": "Loading embedding model..."}), flush=True)
    model = SentenceTransformer('all-MiniLM-L6-v2')

    # Find files
    clar_files, clarinet_toml_files, project_toml_map = find_project_files(SAMPLES_DIR)
    total_files = len(clar_files) + len(clarinet_toml_files)

    if total_files == 0:
        print(json.dumps({
            "type": "error",
            "message": "No files found to ingest"
        }), file=sys.stderr)
        sys.exit(1)

    # Report start
    print(json.dumps({"type": "start", "total": total_files}), flush=True)

    docs, embeddings, metadatas, ids = [], [], [], []
    current = 0

    # Process .clar files first
    print(json.dumps({"type": "info", "message": "Processing .clar files..."}), flush=True)
    for file_path in clar_files:
        current += 1

        try:
            # Determine if this .clar file's project has a Clarinet.toml file (search ancestors)
            file_dir = os.path.dirname(file_path)
            project_root = find_project_root_with_clarinet(file_dir, project_toml_map, str(SAMPLES_DIR))
            has_toml = bool(project_root)

            with open(file_path, "r", encoding="utf-8") as f:
                code = f.read()

            if not code.strip():
                continue

            meta = get_metadata(file_path, SAMPLES_DIR, has_toml)
            emb = get_embedding(model, code)

            docs.append(code)
            embeddings.append(emb)
            metadatas.append(meta)
            ids.append(f"clarity_sample_{current}")

            # Report progress every 10 files
            if current % 10 == 0 or current == 1:
                print(json.dumps({
                    "type": "progress",
                    "current": current,
                    "total": total_files,
                    "message": f"Processing {os.path.basename(file_path)}"
                }), flush=True)

        except Exception as e:
            print(json.dumps({
                "type": "warning",
                "message": f"Error processing {file_path}: {str(e)}"
            }), flush=True)

    # Process Clarinet.toml files
    print(json.dumps({"type": "info", "message": "Processing Clarinet.toml files..."}), flush=True)
    for file_path in clarinet_toml_files:
        current += 1

        try:
            with open(file_path, "r", encoding="utf-8") as f:
                toml_content = f.read()

            if not toml_content.strip():
                continue

            meta = get_metadata(file_path, SAMPLES_DIR, has_toml=True)
            emb = get_embedding(model, toml_content)

            docs.append(toml_content)
            embeddings.append(emb)
            metadatas.append(meta)
            ids.append(f"toml_sample_{current}")

            if current % 10 == 0:
                print(json.dumps({
                    "type": "progress",
                    "current": current,
                    "total": total_files,
                    "message": f"Processing {os.path.basename(file_path)}"
                }), flush=True)

        except Exception as e:
            print(json.dumps({
                "type": "warning",
                "message": f"Error processing {file_path}: {str(e)}"
            }), flush=True)

    # Store in ChromaDB
    if docs:
        print(json.dumps({
            "type": "info",
            "message": f"Storing {len(docs)} documents in ChromaDB..."
        }), flush=True)

        try:
            collection.add(
                documents=docs,
                embeddings=embeddings,
                metadatas=metadatas,
                ids=ids
            )
        except Exception as e:
            print(json.dumps({
                "type": "error",
                "message": f"Failed to store in ChromaDB: {str(e)}"
            }), file=sys.stderr)
            sys.exit(1)

    # Report completion
    print(json.dumps({
        "type": "complete",
        "total_processed": len(docs)
    }), flush=True)


if __name__ == "__main__":
    try:
        ingest_samples()
    except Exception as e:
        print(json.dumps({"type": "error", "message": str(e)}), file=sys.stderr)
        sys.exit(1)
