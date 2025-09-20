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

	const response = await options.fetch(path, {
		method: 'POST',
		body: JSON.stringify(data),
		headers
	});

	if (!response.ok) {
		const json = await response.json();
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

	const response = await options.fetch(path, {
		method: 'PUT',
		body: JSON.stringify(data),
		headers
	});

	if (!response.ok) {
		throw new ApiError(
			`Request failed with invalid status code: ${response.status}`,
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

	const response = await options.fetch(path, {
		method: 'PATCH',
		body: JSON.stringify(data),
		headers
	});

	if (!response.ok) {
		throw new ApiError(
			`Request failed with invalid status code: ${response.status}`,
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

	const response = await options.fetch(`${url}${paramsString}`, {
		headers
	});

	if (!response.ok) {
		console.error('Request failed with invalid status code:', response.status, response.statusText);
		console.error('Response:', await response.json());
		throw new ApiError(
			`Request failed with invalid status code: ${response.status}`,
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

	const response = await options.fetch(path, {
		method: 'DELETE',
		headers
	});

	if (!response.ok) {
		throw new ApiError(
			`Request failed with invalid status code: ${response.status}`,
			response.status
		);
	}
}

function convertParamsToQueryString(params: Record<string, any>, prefix = ''): string {
	const queryString = new URLSearchParams();

	function addParam(key: string, value: any) {
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
