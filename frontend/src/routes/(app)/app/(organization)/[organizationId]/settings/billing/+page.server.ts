import { authenticate } from '$lib/server/auth';

export const load = async () => {
	authenticate();
};

