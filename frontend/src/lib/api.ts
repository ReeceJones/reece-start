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

	let paramsString = '';

	if (options.params && options.paramsSchema) {
		options.paramsSchema.parse(options.params);
		paramsString = `?${new URLSearchParams(options.params).toString()}`;
	}

	const response = await options.fetch(`${path}${paramsString}`, {
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
