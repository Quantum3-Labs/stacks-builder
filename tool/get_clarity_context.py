from typing import Any
import mcp.types as types
from rag import inference_base as base
from rag import inference_gemini as gemini
from tool import tool_factory


class GetClarityContext(tool_factory.ToolFactory):
    def action(self, arguments: dict[str, Any]) -> list[types.ContentBlock]:
        query = arguments.get("query")
        code_results = int(arguments.get("code_results", 5))
        docs_results = int(arguments.get("docs_results", 5))

        gemini_strategy = gemini.GeminiStrategy()
        context = base.InferenceContext(gemini_strategy)

        # Retrieve from BOTH collections (code + docs)
        retrieved_data = context.retrieve_context(
            query, code_results=code_results, docs_results=docs_results
        )

        # Generate answer using the retrieved context
        if getattr(gemini_strategy, "use_sdk", False) and getattr(
            gemini, "GEMINI_SDK_AVAILABLE", False
        ):
            answer, _ = gemini_strategy.answer_with_gemini_sdk(query, retrieved_data)
        else:
            answer, _ = gemini_strategy.answer_with_gemini_rest(query, retrieved_data)

        return [
            types.TextContent(
                type="text",
                text=answer,
            )
        ]


