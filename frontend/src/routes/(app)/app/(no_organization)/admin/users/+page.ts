import { get } from '$lib';
import { getUsersParamsSchema, getUsersResponseSchema } from '$lib/schemas/user';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, fetch }) => {
	const search = url.searchParams.get('search') ?? undefined;
	const cursor = url.searchParams.get('page[cursor]') ?? undefined;

	const apiUrl = cursor ? `/api${cursor}` : `/api/users`;

	console.log('apiUrl', apiUrl);

	const users = await get(apiUrl, {
		fetch,
		responseSchema: getUsersResponseSchema,
		paramsSchema: getUsersParamsSchema,
		params: {
			search,
			page: {
				cursor,
				size: 20
			}
		}
	});

	console.log('users', users);

	return {
		users,
		// Pass along the cursor to help with client-side state management
		currentCursor: cursor
	};
};
