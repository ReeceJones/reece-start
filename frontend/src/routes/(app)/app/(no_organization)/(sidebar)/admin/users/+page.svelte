<script lang="ts">
	import {
		EllipsisVertical,
		HatGlasses,
		Search,
		User,
		ChevronLeft,
		ChevronRight
	} from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { page } from '$app/state';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import { t } from '$lib/i18n';

	const defaultSearch = $derived(page.url.searchParams.get('search') ?? '');

	const { data }: PageProps = $props();
	const users = $derived(data.users);

	function loadPreviousPage() {
		const searchForm = document.getElementById('search-form') as HTMLFormElement;
		const pageCursor = document.getElementById('page-cursor') as HTMLInputElement;

		pageCursor.value = users.links.prev ?? '';
		searchForm.submit();
	}

	function loadNextPage() {
		const searchForm = document.getElementById('search-form') as HTMLFormElement;
		const pageCursor = document.getElementById('page-cursor') as HTMLInputElement;

		pageCursor.value = users.links.next ?? '';
		searchForm.submit();
	}
</script>

<Card>
	<CardBody>
		<CardTitle>{$t('noOrganization.admin.users.title')}</CardTitle>
		<form class="flex gap-2" method="GET" id="search-form">
			<input type="hidden" name="page[cursor]" value={users.links.next ?? ''} id="page-cursor" />

			<input
				type="text"
				class="input-bordered input max-w-96 flex-1"
				placeholder={$t('noOrganization.admin.users.searchPlaceholder')}
				name="search"
				defaultValue={defaultSearch}
			/>
			<button class="btn btn-primary">
				<Search class="size-4" />
				{$t('noOrganization.admin.users.search')}
			</button>
		</form>
		<div>
			<table class="table">
				<thead>
					<tr>
						<th>{$t('noOrganization.admin.users.name')}</th>
						<th>{$t('noOrganization.admin.users.email')}</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{#each users.data as user (user.id)}
						<tr>
							<td>
								<div class="flex items-start gap-3">
									{#if user.meta.logoDistributionUrl}
										<img
											src={user.meta.logoDistributionUrl}
											alt={user.attributes.name}
											class="rounded-box size-8"
										/>
									{:else}
										<User class="rounded-box bg-base-300 size-8" />
									{/if}
									<span class="font-semibold">
										{user.attributes.name}
									</span>
								</div>
							</td>
							<td>
								<a href={`mailto:${user.attributes.email}`} class="link">
									{user.attributes.email}
								</a>
							</td>
							<td class="flex justify-end">
								<div class="dropdown dropdown-end">
									<div tabindex="0" role="button" class="btn btn-square btn-ghost btn-sm">
										<EllipsisVertical class="size-4" />
									</div>
									<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
									<ul
										tabindex="0"
										class="dropdown-content menu z-1 rounded-box bg-base-100 w-52 p-2 shadow-sm"
									>
										<li>
											<button
												onclick={(e) => {
													e.preventDefault();
													// get the child form and submit it
													const form = document.getElementById(
														`impersonate-form-${user.id}`
													) as HTMLFormElement;
													form.submit();
												}}
											>
												<HatGlasses class="size-4" />
												{$t('noOrganization.admin.users.impersonate')}
											</button>

											<form
												method="POST"
												action="?/impersonate"
												id={`impersonate-form-${user.id}`}
												hidden
											>
												<input type="hidden" name="impersonatedUserId" value={user.id} />
											</form>
										</li>
									</ul>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			{#if users.links.prev || users.links.next}
				<div class="mt-4 flex justify-center gap-2">
					{#if users.links.prev}
						<button class="btn btn-ghost btn-sm btn-neutral" onclick={loadPreviousPage}>
							<ChevronLeft class="size-4" />
							{$t('noOrganization.admin.users.previous')}
						</button>
					{/if}
					{#if users.links.next}
						<button class="btn btn-ghost btn-sm btn-neutral" onclick={loadNextPage}>
							{$t('noOrganization.admin.users.next')}
							<ChevronRight class="size-4" />
						</button>
					{/if}
				</div>
			{/if}
		</div>
	</CardBody>
</Card>
