# Clarity Builder

Clarity Builder is an MCP-enabled RAG system that enhances Clarity smart contract coding for Cursor/VS Code. It ingests official Clarity docs and sample projects into ChromaDB, retrieves the most relevant context, and uses Gemini to produce accurate answers and code.

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
- Servers
  - MCP: `MCP_Server/server.py` (tools exposed over StreamableHTTP).
  - API: `API/auth_server.py` (auth, API keys) and `API/api_server.py` (OpenAI-compatible chat completions).
- Automation
  - `automated_ingestion_job/`: config, cleanup manager, orchestrator, scheduler, task.

## Requirements

- Python 3.11+
- pip packages (install once):

```bash
pip install -r requirements.txt
```

Create a `.env` in the project root:

```env
GEMINI_API_KEY=your-gemini-api-key
SECRET_KEY=your-secret-key-for-jwt
```

## One-Click Reinitialization (Recommended)

This will clean `chromadb_data/`, clone Clarity repos and docs, and ingest everything.

```bash
python automated_ingestion_job/update_data_task.py
```

To run on a schedule (APScheduler):

```bash
python automated_ingestion_job/scheduler.py
```

Cron schedule is configured in `automated_ingestion_job/config.json`.

## Manual Pipeline

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

## Run the API Servers

1) Auth server (users, API keys)

```bash
python -m uvicorn API.auth_server:app --reload --port 8001
```

2) RAG API server (OpenAI-compatible /v1/chat/completions)

```bash
python -m uvicorn API.api_server:app --reload --port 8100
```

3) Example client (register, login, create key, call API)

```bash
python API/client_example.py
```

## MCP Server for Cursor

Server entrypoint (registers tools `get_clarity_context` and `generate_clarity_code`):

```bash
python MCP_Server/server.py --port 3001
```

Client configuration (Cursor → Settings → MCP Tools):

```json
{
  "mcpServers": {
    "clarity-builder": {
      "command": "python",
      "args": ["MCP_Server/server.py"],
      "env": { "PYTHONPATH": ".", "GEMINI_API_KEY": "your-gemini-api-key" }
    }
  }
}
```

Alternatively, you can expose a simple HTTP MCP-style context endpoint:

```bash
python -m uvicorn API.mcp_api_server:app --reload --port 8200
```

## Project Layout (key paths)

- `clarity_official_docs/`: cloned documentation
- `temp_clarity_clone/projects/`: cloned example projects (.clar, Clarinet.toml)
- `chromadb_data/`: ChromaDB persistent data
- `ingest/`: ingestion scripts
- `rag/`: retrieval and generation strategies
- `tool/`: MCP tools
- `MCP_Server/`: MCP server and config
- `API/`: auth + RAG APIs and support modules
- `automated_ingestion_job/`: full reinit pipeline

## Troubleshooting

- Port in use: change ports in the commands above.
- Missing data: re-run the ingestion (manual or automated).
- Invalid API key: create a key via `API/client_example.py` or `POST /api-keys`.
- Fresh DB: delete `chromadb_data/` and re-run the pipeline.

## Notes

- Embedding model: `all-MiniLM-L6-v2` (local SentenceTransformer) for ingestion and retrieval.
- Collections: `clarity_code_samples` and `clarity_docs`.
- Gemini model (default): `models/gemini-2.5-flash`.
