import type { HandleClientError } from '@sveltejs/kit';
import * as Sentry from '@sentry/sveltekit';
import { env } from '$env/dynamic/public';

interface ErrorWithStatus extends Error {
	status?: number;
}

function hasStatus(error: unknown): error is ErrorWithStatus {
	return error instanceof Error && 'status' in error;
}

function getErrorStatus(error: unknown): number {
	if (hasStatus(error)) {
		return error.status ?? 500;
	}
	return 500;
}

const defaultErrorHandler: HandleClientError = ({ error, event }) => {
	console.error('[Unhandled Error]', {
		error,
		url: event.url.toString(),
		status: getErrorStatus(error),
		message: error instanceof Error ? error.message : String(error),
		stack: error instanceof Error ? error.stack : undefined
	});

	return {
		message: error instanceof Error ? error.message : 'An unexpected error occurred',
		status: getErrorStatus(error)
	};
};

export const handleError = Sentry.handleErrorWithSentry(defaultErrorHandler);

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

	Sentry.init({
		dsn: env.PUBLIC_SENTRY_DSN,
		// Adds request headers and IP for users, for more info visit:
		// https://docs.sentry.io/platforms/javascript/guides/sveltekit/configuration/options/#sendDefaultPii
		sendDefaultPii: true
	});
}
