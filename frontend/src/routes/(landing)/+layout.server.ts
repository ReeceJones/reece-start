import { isLoggedIn } from '$lib/server/auth';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async () => {
	return {
		isLoggedIn: isLoggedIn()
	};
};
