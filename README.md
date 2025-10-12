# Clarity Builder

Clarity Builder is an MCP-enabled RAG system that enhances Clarity smart contract coding for Cursor/VS Code. It ingests official Clarity docs and sample projects into ChromaDB, retrieves the most relevant context, and uses Gemini to produce accurate answers and code.

## RAG Pipeline

<img width="3978" height="3612" alt="Clarity_RAG" src="https://github.com/user-attachments/assets/717b3721-49cc-4042-b2e5-cba625f5f84f" />


## Architecture

- Ingestion
  - `clone_clarity_repos.py`: clones example Clarity projects (with `.clar` and `Clarinet.toml`).
  - `clone_clarity_docs.py`: clones official Clarity documentation into `clarity_official_docs/`.
  - `ingest/clarity_samples_ingester.py`: embeds `.clar` files and `Clarinet.toml` into `clarity_code_samples`.
  - `ingest/clarity_docs_ingester.py`: chunks and embeds docs into `clarity_docs`.
- RAG Core
  - `rag/inference_base.py`: retrieves from `clarity_code_samples` and `clarity_docs` and builds prompts.
  - `rag/inference_gemini.py`: Gemini strategy for generation using the retrieved context.
- Tools (MCP)
  - `tool/get_clarity_context.py`: returns top doc chunks for a query.
  - `tool/generate_clarity_code.py`: generates responses/code with context.
- Server
  - MCP HTTP server: `MCP_Server/server.py` (tools exposed over StreamableHTTP at `/mcp/`).

## Requirements

- Python 3.11+
- pip packages (install once):

```bash
pip install -r requirements.txt
```

Create a `.env` in the project root (only this is required):

```env
GEMINI_API_KEY=your-gemini-api-key
```

## Manual Pipeline (Recommended)

1) Clone sources

```bash
python clone_clarity_repos.py
python clone_clarity_docs.py
```

2) Ingest into ChromaDB

```bash
python ingest/clarity_samples_ingester.py   # adds .clar + Clarinet.toml → clarity_code_samples
python ingest/clarity_docs_ingester.py      # adds markdown docs → clarity_docs
```

3) Inspect the database (optional)

```bash
python inspect_chromadb.py
```

## MCP Server for Cursor (HTTP)

1) Start the MCP server (JSON responses enabled):

```bash
python MCP_Server/server.py --port 3001 --json-response
```

2) Configure Cursor → Settings → MCP Tools to point to the HTTP server:

```json
{
  "mcpServers": {
    "clarity-builder": {
      "transport": "http",
      "url": "http://127.0.0.1:3001/mcp/"
    }
  }
}
```

Notes:
- The trailing slash in `/mcp/` avoids HTTP redirects.
- You should see tools `get_clarity_context` and `generate_clarity_code` available.

## Project Layout (key paths)

- `clarity_official_docs/`: cloned documentation
- `temp_clarity_clone/projects/`: cloned example projects (.clar, Clarinet.toml)
- `chromadb_data/`: ChromaDB persistent data
- `ingest/`: ingestion scripts
- `rag/`: retrieval and generation strategies
- `tool/`: MCP tools
- `MCP_Server/`: MCP server and config
- `automated_ingestion_job/`: full reinit pipeline

## Troubleshooting

- Port in use: change the `--port` value when starting the server.
- Missing data: re-run ingestion (manual or automated).
- 404s on `/`: expected; only `/mcp/` is served.
- Fresh DB: delete `chromadb_data/` and re-run the pipeline.

## Notes

- Embedding model: `all-MiniLM-L6-v2` (local SentenceTransformer) for ingestion and retrieval.
- Collections: `clarity_code_samples` and `clarity_docs`.
- Gemini model (default): `models/gemini-2.5-flash`.
