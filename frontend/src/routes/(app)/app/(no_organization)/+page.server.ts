import { authenticate } from '$lib/server/auth';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	authenticate();
};
