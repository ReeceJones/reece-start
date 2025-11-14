/**
 * Factory function for mocking SvelteKit's enhance function from $app/forms
 *
 * This mock allows testing components that use the enhance action on forms.
 * It properly handles the callback that can return a submit handler function.
 *
 * Usage in test files:
 * ```ts
 * vi.mock('$app/forms', async () => {
 *   const { createMockEnhance } = await import('$lib/test-utils');
 *   return createMockEnhance();
 * });
 * ```
 */
export async function createMockEnhance() {
	const { vi } = await import('vitest');
	return {
		enhance: vi.fn(
			(
				callback?: () =>
					| void
					| (({ update, result }: { update: () => void; result?: unknown }) => void | Promise<void>)
			) => {
				// Return a Svelte action function that can be used with use:
				return (node: HTMLFormElement) => {
					if (callback) {
						const submitHandler = callback();
						// If callback returns a function, set up form submission handler
						if (typeof submitHandler === 'function') {
							node.addEventListener('submit', async (e) => {
								e.preventDefault();
								// Call the submit handler with mock update function
								await submitHandler({ update: vi.fn() });
							});
						}
					}
					// Return cleanup function for the action
					return {
						destroy: vi.fn()
					};
				};
			}
		)
	};
}
