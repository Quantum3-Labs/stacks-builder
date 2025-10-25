# Stacks Builder Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that provides **AI-powered Clarity programming assistance** through Retrieval-Augmented Generation (RAG). Get intelligent code suggestions, context-aware completions, and instant access to Clarity documentation directly in your IDE.

## ✨ Features

- 🔍 **Smart Context Retrieval** - Search through 40+ Clarity code samples and official documentation
- 🤖 **AI Code Generation** - Generate Clarity code with LLM assistance (Gemini/OpenAI/Claude)
- ⚡ **RAG-Powered** - Combines vector similarity search with LLM generation for accurate results
- 🎯 **IDE Integration** - Works seamlessly with Cursor, Claude Desktop, and MCP-compatible editors
- 🔒 **Type-Safe** - Built with TypeScript and Zod validation
- 🌐 **Production Ready** - Backed by a robust Go backend with ChromaDB vector store

## 🚀 Quick Start

### Step 1: Get an API Key

You need an API key from the Stacks Builder backend to use this MCP server.

1. Visit the Swagger UI: **<https://stacks-builder.q3labs.io/swagger/index.html>**
2. Register a new account via `/api/v1/auth/register` endpoint
3. Login using `/api/v1/auth/login` endpoint
4. Generate your API key from `/api/v1/keys` endpoint
5. **Save this key** - you'll need it for the next step

### Step 2: Configure in Cursor

Add this configuration to your Cursor MCP settings file at `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "stacks-builder": {
      "command": "npx",
      "args": [
        "-y",
        "@q3labs/stacks-builder"
      ],
      "env": {
        "API_KEY": "your-api-key-here",
        "BACKEND_URL": "https://stacks-builder.q3labs.io"
      }
    }
  }
}
```

**Important:**

- Replace `your-api-key-here` with your actual API key from Step 1
- Use `https://stacks-builder.q3labs.io` for production (recommended)
- Use `http://localhost:8080` if you're running the backend locally

### Step 3: Restart Cursor

**Completely restart Cursor** (not just reload) for the changes to take effect.

## 🛠️ Available Tools

Once configured, you'll have access to these MCP tools in Cursor:

### `get_clarity_context`

Retrieves relevant Clarity code snippets and documentation from the RAG system based on your query.

**Parameters:**

- `query` (required) - What you're looking for
- `n_results` (optional) - Number of matches to return (1-5, default: 5)

**Example usage:**

```
"How do I create a data variable in Clarity?"
```

### `generate_clarity_code`

Generates Clarity code using backend RAG context combined with LLM generation.

**Parameters:**

- `query` (required) - Your code generation request
- `temperature` (optional) - Controls creativity (0-2, default varies by provider)
- `max_tokens` (optional) - Maximum tokens to generate

**Example usage:**

```
"Generate a Clarity contract for managing user profiles with CRUD operations"
```

## 🔧 Troubleshooting

### MCP Tools Not Showing Up

If the MCP server shows "No tools, prompts, or resources" after restarting Cursor:

1. **Try global installation:**

   ```bash
   npm install -g @q3labs/stacks-builder
   ```

2. **Update your config:**

   ```json
   {
     "mcpServers": {
       "stacks-builder": {
         "command": "stacks-builder",
         "args": [],
         "env": {
           "API_KEY": "your-api-key-here",
           "BACKEND_URL": "https://stacks-builder.q3labs.io"
         }
       }
     }
   }
   ```

### Node.js Version

This package requires **Node.js 22.0.0 or higher**. Check your version:

```bash
node --version
```

## 📚 How It Works

1. **You ask a question** about Clarity in your IDE
2. **Context retrieval** searches the ChromaDB vector store for relevant code examples and documentation
3. **LLM generation** combines the retrieved context with your prompt
4. **Smart response** returns context-aware code suggestions directly in Cursor

```
┌─────────────┐
│   Cursor    │  Your IDE
└──────┬──────┘
       │ MCP Protocol
┌──────▼─────────────┐
│   Stacks Builder   │  MCP Server (this package)
│   MCP Server       │ 
└─────────┬──────────┘
       │ HTTP/REST
┌──────▼──────┐
│  Backend    │  Go API + Python RAG
│  ChromaDB   │  Vector Store + LLM
└─────────────┘
```

## 🌟 Use Cases

- **Learning Clarity** - Get instant examples and documentation
- **Building Contracts** - Generate boilerplate code for common patterns
- **Debugging** - Find similar code examples to solve issues
- **Best Practices** - Learn from 40+ curated Clarity samples

## 📖 Documentation

For complete documentation including:

- Backend setup and local development
- Contributing guidelines
- Architecture diagrams
- Advanced configuration

Visit the **[main repository](https://github.com/Quantum3-Labs/stacks-builder#readme)**.

## 🤝 Support

- **Issues**: [GitHub Issues](https://github.com/Quantum3-Labs/stacks-builder/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Quantum3-Labs/stacks-builder/discussions)
- **Repository**: [Quantum3-Labs/stacks-builder](https://github.com/Quantum3-Labs/stacks-builder)

## 📄 License

MIT License - see [LICENSE](https://github.com/Quantum3-Labs/stacks-builder/blob/main/LICENSE) file for details.

---

Built with ❤️ by [Quantum3 Labs](https://github.com/Quantum3-Labs) for the Stacks blockchain ecosystem.
