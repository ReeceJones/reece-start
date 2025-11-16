import type { HandleClientError } from '@sveltejs/kit';

export const handleError: HandleClientError = ({ error, event }) => {
	console.error('[Unhandled Error]', {
		error,
		url: event.url.toString(),
		method: event.request.method,
		status: error instanceof Error ? (error as any).status : 500,
		message: error instanceof Error ? error.message : String(error),
		stack: error instanceof Error ? error.stack : undefined
	});

	return {
		message: error instanceof Error ? error.message : 'An unexpected error occurred',
		status: error instanceof Error ? ((error as any).status ?? 500) : 500
	};
};

// Handle unhandled promise rejections
if (typeof window !== 'undefined') {
	window.addEventListener('error', (event) => {
		console.error('[Unhandled JavaScript Error]', {
			message: event.message,
			filename: event.filename,
			lineno: event.lineno,
			colno: event.colno,
			error: event.error,
			stack: event.error?.stack
		});
	});

	window.addEventListener('unhandledrejection', (event) => {
		console.error('[Unhandled Promise Rejection]', {
			reason: event.reason,
			promise: event.promise,
			error: event.reason instanceof Error ? event.reason : undefined,
			stack: event.reason instanceof Error ? event.reason.stack : undefined
		});
	});
}
