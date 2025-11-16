import { performAuthenticationCheck } from '$lib/server/auth';
import { type Handle, type HandleServerError, type Redirect, type HttpError } from '@sveltejs/kit';

interface ErrorWithLocation extends Error {
	location: string;
}

interface ErrorWithRequiredStatus extends Error {
	status: number;
}

function hasLocation(error: Error): error is ErrorWithLocation {
	if ('location' in error) {
		const location = (error as ErrorWithLocation).location;
		return typeof location === 'string';
	}
	return false;
}

function hasRequiredStatus(error: Error): error is ErrorWithRequiredStatus {
	if ('status' in error) {
		const status = (error as ErrorWithRequiredStatus).status;
		return typeof status === 'number';
	}
	return false;
}

function isSvelteKitError(error: unknown): error is Redirect | HttpError {
	// Check if it's a SvelteKit error by checking for known properties
	// Redirect errors have a 'location' property, HttpError has 'status'
	if (error instanceof Error) {
		return hasLocation(error) || hasRequiredStatus(error);
	}
	return false;
}

export const handle: Handle = async ({ event, resolve }) => {
	const startTime = Date.now();
	console.log(`[Server Request] ${event.request.method} ${event.url.pathname}${event.url.search}`);

	try {
		await performAuthenticationCheck(event);
		const response = await resolve(event);
		const duration = Date.now() - startTime;
		console.log(
			`[Server Response] ${event.request.method} ${event.url.pathname} ${response.status} (${duration}ms)`
		);
		return response;
	} catch (error) {
		// Let SvelteKit handle redirects and HTTP errors - don't log these as errors
		if (isSvelteKitError(error)) {
			throw error;
		}

		const duration = Date.now() - startTime;
		console.error(`[Server Error] ${event.request.method} ${event.url.pathname} (${duration}ms)`, {
			error,
			message: error instanceof Error ? error.message : String(error),
			stack: error instanceof Error ? error.stack : undefined
		});
		throw error;
	}
};

interface ErrorWithOptionalStatus extends Error {
	status?: number;
}

function hasOptionalStatus(error: Error): error is ErrorWithOptionalStatus {
	return 'status' in error;
}

function getErrorStatus(error: unknown): number {
	if (error instanceof Error && hasOptionalStatus(error)) {
		return typeof error.status === 'number' ? error.status : 500;
	}
	return 500;
}

export const handleError: HandleServerError = ({ error, event }) => {
	console.error('[Server Unhandled Error]', {
		error,
		url: event.url.toString(),
		method: event.request.method,
		message: error instanceof Error ? error.message : String(error),
		stack: error instanceof Error ? error.stack : undefined
	});

	return {
		message: error instanceof Error ? error.message : 'An unexpected error occurred',
		status: getErrorStatus(error)
	};
};
