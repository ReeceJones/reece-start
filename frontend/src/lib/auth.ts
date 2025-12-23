import type { UserScope } from './schemas/jwt';
import { getContext, setContext } from 'svelte';

export const scopesKey = 'membership-scopes';
const isImpersonatingUserKey = 'is-impersonating-user';

export function setScopes(scopes: () => UserScope[]) {
	setContext(scopesKey, scopes);
}

export function getScopes() {
	return getContext<() => UserScope[]>(scopesKey)?.() ?? [];
}

export function hasScope(scope: UserScope): boolean {
	return getScopes().includes(scope) ?? false;
}

export function setIsImpersonatingUser(isImpersonatingUser: () => boolean) {
	setContext(isImpersonatingUserKey, isImpersonatingUser);
}

export function getIsImpersonatingUser() {
	return getContext<() => boolean>(isImpersonatingUserKey)?.() ?? false;
}
