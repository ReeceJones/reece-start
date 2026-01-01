import type { ActionResult } from '@sveltejs/kit';

export interface FormResult {
	success: boolean;
	message: string;
	[key: string]: unknown;
}

type MaybePromise<T> = T | Promise<T>;

export type FormResultCallback = MaybePromise<
	(opts: {
		formData: FormData;
		formElement: HTMLFormElement;
		action: URL;
		result: ActionResult<FormResult, FormResult>;
		/**
		 * Call this to get the default behavior of a form submission response.
		 * @param options Set `reset: false` if you don't want the `<form>` values to be reset after a successful submission.
		 * @param invalidateAll Set `invalidateAll: false` if you don't want the action to call `invalidateAll` after submission.
		 */
		update: (options?: { reset?: boolean; invalidateAll?: boolean }) => Promise<void>;
	}) => MaybePromise<void>
>;

export function getFormResult(
	result: ActionResult<FormResult, FormResult>
): FormResult | undefined {
	if (result.type === 'success' || result.type === 'failure') {
		return result.data;
	}

	return undefined;
}
