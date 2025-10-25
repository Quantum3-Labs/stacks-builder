import contextlib
import logging
from collections.abc import AsyncIterator
from typing import Any

import click
import mcp.types as types
from mcp.server.lowlevel import Server
from mcp.server.streamable_http_manager import StreamableHTTPSessionManager
from starlette.applications import Starlette
from starlette.middleware.cors import CORSMiddleware
from starlette.routing import Mount
from starlette.types import Receive, Scope, Send
import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from starlette.applications import Starlette
from starlette.middleware import Middleware
import tool.tool_factory as tool_factory
import tool.get_clarity_context as get_clarity_context
import tool.generate_clarity_code as generate_clarity_code
import uvicorn

logger = logging.getLogger(__name__)
logging.getLogger("mcp.server.streamable_http").setLevel(logging.WARNING)

class _SuppressClosedResourceError(logging.Filter):
    def filter(self, record: logging.LogRecord) -> bool:
        msg = record.getMessage()
        return "ClosedResourceError" not in msg

logging.getLogger("mcp.server.streamable_http").addFilter(_SuppressClosedResourceError())


@click.command()
@click.option("--port", default=3001, help="Port to listen on for HTTP")
@click.option(
    "--log-level",
    default="INFO",
    help="Logging level (DEBUG, INFO, WARNING, ERROR, CRITICAL)",
)
@click.option(
    "--json-response",
    is_flag=True,
    default=False,
    help="Enable JSON responses instead of SSE streams",
)
def main(
    port: int,
    log_level: str,
    json_response: bool,
) -> int:
    # Configure logging
    logging.basicConfig(
        level=getattr(logging, log_level.upper()),
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )
    GET_CLARITY_CONTEXT_TOOL = "get_clarity_context"
    GENERATE_CLARITY_CODE_TOOL = "generate_clarity_code"
    # init Factory
    tool_factory.ToolFactory.register(
        GET_CLARITY_CONTEXT_TOOL, get_clarity_context.GetClarityContext
    )
    tool_factory.ToolFactory.register(
        GENERATE_CLARITY_CODE_TOOL, generate_clarity_code.GenerateClarityCode
    )

    app = Server("clarity-builder-mcp-server")

    @app.call_tool()
    async def call_tool(
        name: str, arguments: dict[str, Any]
    ) -> list[types.ContentBlock]:
        ctx = app.request_context
        if name != "":
            tool = tool_factory.ToolFactory.create(name)
            return tool.action(arguments)

    @app.list_tools()
    async def list_tools() -> list[types.Tool]:
        return [
            types.Tool(
                name="get_clarity_context",
                description=(
                    "Retrieves relevant Clarity code examples based on a query."
                ),
                inputSchema={
                    "type": "object",
                    "required": ["query"],
                    "properties": {
                        "query": {
                            "type": "string",
                            "description": "What you're looking for",
                        },
                    },
                },
            ),
            types.Tool(
                name="generate_clarity_code",
                description=(
                    "Generates complete Clarity code using Gemini with RAG context."
                ),
                inputSchema={
                    "type": "object",
                    "required": ["query"],
                    "properties": {
                        "query": {
                            "type": "string",
                            "description": "Your code generation request",
                        },
                    },
                },
            ),
        ]

    session_manager = StreamableHTTPSessionManager(
        app=app,
        event_store=None,
        json_response=json_response,
        stateless=True,
    )

    async def handle_streamable_http(
        scope: Scope, receive: Receive, send: Send
    ) -> None:
        await session_manager.handle_request(scope, receive, send)

    @contextlib.asynccontextmanager
    async def lifespan(app: Starlette) -> AsyncIterator[None]:
        """Context manager for session manager."""
        async with session_manager.run():
            logger.info("Application started with StreamableHTTP session manager!")
            try:
                yield
            finally:
                logger.info("Application shutting down...")

    middlewares = []
    starlette_app = Starlette(
        debug=True,
        routes=[
            Mount("/mcp", app=handle_streamable_http),
        ],
        lifespan=lifespan,
        middleware=middlewares,
    )

    starlette_app = CORSMiddleware(
        starlette_app,
        allow_origins=["*"],
        allow_methods=["GET", "POST", "DELETE"],
        expose_headers=["Mcp-Session-Id"],
    )

    uvicorn.run(starlette_app, host="127.0.0.1", port=port)


if __name__ == "__main__":
    main()


