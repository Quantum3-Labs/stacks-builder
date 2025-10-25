import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';

const BACKEND_URL_ENV = 'BACKEND_URL';
const DEFAULT_BACKEND_BASE_URL = 'http://localhost:8080';
const RAG_RETRIEVE_PATH = '/api/v1/rag/retrieve';

const GetClarityContextArgsSchema = z.object({
	query: z
		.string()
		.min(1)
		.describe("What you're looking for"),
	n_results: z
		.number()
		.int()
		.min(1)
		.max(5)
		.optional()
		.describe('How many matches to return (1-5, defaults to 5).'),
});

type GetClarityContextArgs = z.infer<typeof GetClarityContextArgsSchema>;

type ToolOptions = {
	apiKey: string;
	baseUrl?: string;
};

type RetrieveContextResponse = {
	formatted_context?: string;
	warning?: string;
	error?: string;
};

function resolveBackendBaseUrl(candidate?: string): string {
	const explicit = candidate?.trim();
	if (explicit) {
		return explicit.replace(/\/+$/, '');
	}

	const fromEnv = process.env[BACKEND_URL_ENV]?.trim();
	if (fromEnv) {
		return fromEnv.replace(/\/+$/, '');
	}

	return DEFAULT_BACKEND_BASE_URL;
}

export function registerGetClarityContext(
	server: McpServer,
	options: ToolOptions,
) {
	const apiKey = options.apiKey.trim();

	if (!apiKey) {
		throw new Error(
			'API key is required to register get_clarity_context tool.',
		);
	}

	const backendBaseUrl = resolveBackendBaseUrl(options.baseUrl);

	server.registerTool(
		'get_clarity_context',
		{
			title: 'Retrieve Clarity Context',
			description:
				'Fetches relevant Clarity code and documentation snippets from the backend RAG service.',
			inputSchema: GetClarityContextArgsSchema.shape,
		},
		async ({ query, n_results }: GetClarityContextArgs) => {
			const cappedResults =
				typeof n_results === 'number'
					? Math.min(Math.max(Math.trunc(n_results), 1), 5)
					: 5;

			const payload = {
				query: query.trim(),
				n_results: cappedResults,
			};

			try {
				const response = await fetch(
					`${backendBaseUrl}${RAG_RETRIEVE_PATH}`,
					{
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
							'x-api-key': apiKey,
						},
						body: JSON.stringify(payload),
					},
				);

				const rawBody = await response.text();
				const parsedBody: RetrieveContextResponse | undefined =
					rawBody.length > 0 ? JSON.parse(rawBody) : undefined;

				if (!response.ok) {
					const backendError =
						parsedBody?.error ??
						`Backend returned status ${response.status}`;
					throw new Error(backendError);
				}

				if (!parsedBody) {
					throw new Error('Backend response was empty.');
				}

				const formattedContext =
					parsedBody.formatted_context?.trim() ?? '';
				const warning = parsedBody.warning?.trim();

				const summaryPieces = [
					formattedContext.length > 0
						? 'Formatted context retrieved successfully.'
						: 'Backend returned an empty formatted context.',
				];

				if (warning && warning.length > 0) {
					summaryPieces.push(`Warning: ${warning}`);
				}

				return {
					content: [
						{
							type: 'text' as const,
							text: summaryPieces.join(' '),
						},
						{
							type: 'text' as const,
							text:
								formattedContext.length > 0
									? `Formatted context:\n${formattedContext}`
									: 'No formatted context returned.',
						},
					],
				};
			} catch (error) {
				const message =
					error instanceof Error ? error.message : String(error);
				throw new Error(
					`Failed to retrieve Clarity context: ${message}`,
				);
			}
		},
	);
}
