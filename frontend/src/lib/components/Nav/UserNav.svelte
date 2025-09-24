<script lang="ts">
	import { EyeOff, LogOut, Settings, User } from 'lucide-svelte';
	import { getSelfUserResponseSchema } from '$lib/schemas/user';
	import type { z } from 'zod';
	import { page } from '$app/state';
	import { getIsImpersonatingUser } from '$lib/auth';

	const { user }: { user: z.infer<typeof getSelfUserResponseSchema> } = $props();
	const organizationId = $derived(page.params.organizationId);
	const profileHref = $derived(organizationId ? `/app/${organizationId}/profile` : '/app/profile');
	const isImpersonatingUser = $derived(getIsImpersonatingUser());
</script>

<ul class="menu menu-vertical w-full">
	<li class="w-full">
		<div class="dropdown dropdown-top dropdown-start w-full p-0">
			<div tabindex="0" role="button" class="flex w-full gap-2 px-3 py-1.5">
				{#if user.data.meta.logoDistributionUrl}
					<img src={user.data.meta.logoDistributionUrl} alt="User logo" class="size-6 rounded-sm" />
				{:else}
					<User class="size-5" />
				{/if}
				{user.data.attributes.name ?? 'Profile'}
			</div>
			<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
			<ul
				tabindex="0"
				class="dropdown-content menu bg-base-100 rounded-box z-1 ml-0 w-52 -translate-y-1.5 p-2 shadow-sm"
			>
				<li>
					<a href={profileHref}>
						<Settings class="size-4" />
						Settings
					</a>
				</li>
				<li>
					<button class="text-error flex items-center gap-2" type="submit" form="signout-form">
						<LogOut class="size-4" />
						Logout
					</button>
				</li>
				{#if isImpersonatingUser}
					<li>
						<button
							class="text-error flex items-center gap-2"
							type="submit"
							form="stop-impersonation-form"
						>
							<EyeOff class="size-4" />
							Stop Impersonation
						</button>
					</li>
				{/if}
			</ul>
			<form
				action="/app?/signout"
				method="POST"
				enctype="multipart/form-data"
				id="signout-form"
			></form>
			<form
				action="/app?/stopImpersonation"
				method="POST"
				enctype="multipart/form-data"
				id="stop-impersonation-form"
			></form>
		</div>
	</li>
</ul>
