import { performAuthenticationCheck } from '$lib/server/auth';
import { type Handle, type HandleServerError, type Redirect, type HttpError } from '@sveltejs/kit';

function isSvelteKitError(error: unknown): error is Redirect | HttpError {
	// Check if it's a SvelteKit error by checking for known properties
	// Redirect errors have a 'location' property, HttpError has 'status'
	if (error && typeof error === 'object') {
		const err = error as any;
		return (
			err instanceof Error && ('location' in err || (err.status && typeof err.status === 'number'))
		);
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
		status: (error as any)?.status ?? 500
	};
};
