import { z } from 'zod';

export class ApiError extends Error {
	constructor(
		message: string,
		public code: number
	) {
		super(message);
	}
}

export async function post<T extends z.ZodTypeAny, K extends z.ZodTypeAny>(
	path: string,
	data: z.infer<T>,
	options: {
		fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;
		responseSchema: K;
		requestSchema: T;
	}
): Promise<z.infer<K>> {
	try {
		options.requestSchema.parse(data);
	} catch (error) {
		console.warn('Invalid request data', error, data);
	}

	const headers = {
		'Content-Type': 'application/json'
	};

	const startTime = Date.now();
	console.log(`[API Request] POST ${path}`);

	const response = await options.fetch(path, {
		method: 'POST',
		body: JSON.stringify(data),
		headers
	});

	const duration = Date.now() - startTime;
	console.log(`[API Response] POST ${path} ${response.status} (${duration}ms)`);

	if (!response.ok) {
		const json = await response.json().catch(() => ({}));
		throw new ApiError(
			json.message ?? `Request failed with invalid status code: ${response.status}`,
			response.status
		);
	}

	const parsedData = await response.json();

	return options.responseSchema.parse(parsedData);
}

export async function put<T extends z.ZodTypeAny, K extends z.ZodTypeAny>(
	path: string,
	data: z.infer<T>,
	options: {
		fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;
		responseSchema: K;
		requestSchema: T;
	}
): Promise<z.infer<K>> {
	try {
		options.requestSchema.parse(data);
	} catch (error) {
		console.warn('Invalid request data', error, data);
	}

	const headers = {
		'Content-Type': 'application/json'
	};

	const startTime = Date.now();
	console.log(`[API Request] PUT ${path}`);

	const response = await options.fetch(path, {
		method: 'PUT',
		body: JSON.stringify(data),
		headers
	});

	const duration = Date.now() - startTime;
	console.log(`[API Response] PUT ${path} ${response.status} (${duration}ms)`);

	if (!response.ok) {
		const json = await response.json().catch(() => ({}));
		throw new ApiError(
			json.message ?? `Request failed with invalid status code: ${response.status}`,
			response.status
		);
	}

	const parsedData = await response.json();

	return options.responseSchema.parse(parsedData);
}

export async function patch<T extends z.ZodTypeAny, K extends z.ZodTypeAny>(
	path: string,
	data: z.infer<T>,
	options: {
		fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;
		responseSchema: K;
		requestSchema: T;
	}
): Promise<z.infer<K>> {
	try {
		options.requestSchema.parse(data);
	} catch (error) {
		console.warn('Invalid request data', error, data);
	}

	const headers = {
		'Content-Type': 'application/json'
	};

	const startTime = Date.now();
	console.log(`[API Request] PATCH ${path}`);

	const response = await options.fetch(path, {
		method: 'PATCH',
		body: JSON.stringify(data),
		headers
	});

	const duration = Date.now() - startTime;
	console.log(`[API Response] PATCH ${path} ${response.status} (${duration}ms)`);

	if (!response.ok) {
		const json = await response.json().catch(() => ({}));
		throw new ApiError(
			json.message ?? `Request failed with invalid status code: ${response.status}`,
			response.status
		);
	}

	const parsedData = await response.json();

	return options.responseSchema.parse(parsedData);
}

export async function get<T extends z.ZodTypeAny, K extends z.ZodTypeAny>(
	path: string,
	options: {
		fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;
		responseSchema: T;
		paramsSchema?: K;
		params?: z.infer<K>;
	}
): Promise<z.infer<T>> {
	const headers = {
		'Content-Type': 'application/json'
	};

	// parse out the params string from the path
	let paramsString = path.includes('?') ? path.slice(path.indexOf('?')) : '';
	const url = path.includes('?') ? path.slice(0, path.indexOf('?')) : path;

	if (options.params && options.paramsSchema) {
		options.paramsSchema.parse(options.params);
		// convert nested objects to foo[bar]=baz format
		paramsString = paramsString
			? `${paramsString}&${convertParamsToQueryString(options.params)}`
			: `?${convertParamsToQueryString(options.params)}`;
	}

	const fullPath = `${url}${paramsString}`;
	const startTime = Date.now();
	console.log(`[API Request] GET ${fullPath}`);

	const response = await options.fetch(fullPath, {
		headers
	});

	const duration = Date.now() - startTime;
	console.log(`[API Response] GET ${fullPath} ${response.status} (${duration}ms)`);

	if (!response.ok) {
		const json = await response.json().catch(() => ({}));
		throw new ApiError(
			json.message ?? `Request failed with invalid status code: ${response.status}`,
			response.status
		);
	}

	const parsedData = await response.json();

	return options.responseSchema.parse(parsedData);
}

export async function del(
	path: string,
	options: {
		fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;
	}
): Promise<void> {
	const headers = {
		'Content-Type': 'application/json'
	};

	const startTime = Date.now();
	console.log(`[API Request] DELETE ${path}`);

	const response = await options.fetch(path, {
		method: 'DELETE',
		headers
	});

	const duration = Date.now() - startTime;
	console.log(`[API Response] DELETE ${path} ${response.status} (${duration}ms)`);

	if (!response.ok) {
		const json = await response.json().catch(() => ({}));
		throw new ApiError(
			json.message ?? `Request failed with invalid status code: ${response.status}`,
			response.status
		);
	}
}

function convertParamsToQueryString(params: Record<string, unknown>, prefix = ''): string {
	const queryString = new URLSearchParams();

	function addParam(key: string, value: unknown) {
		if (value === null || value === undefined) {
			return;
		}

		if (Array.isArray(value)) {
			value.forEach((item, index) => {
				addParam(`${key}[${index}]`, item);
			});
		} else if (typeof value === 'object') {
			for (const [subKey, subValue] of Object.entries(value)) {
				addParam(`${key}[${subKey}]`, subValue);
			}
		} else {
			queryString.set(key, String(value));
		}
	}

	for (const [key, value] of Object.entries(params)) {
		const paramKey = prefix ? `${prefix}[${key}]` : key;
		addParam(paramKey, value);
	}

	return queryString.toString();
}
