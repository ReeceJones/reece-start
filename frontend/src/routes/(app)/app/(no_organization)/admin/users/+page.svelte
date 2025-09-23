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

<div class="card bg-base-200 border-base-300 border shadow-sm">
	<div class="card-body">
		<h2 class="card-title">Users</h2>
		<form class="flex gap-2" method="GET" id="search-form">
			<input type="hidden" name="page[cursor]" value={users.links.next ?? ''} id="page-cursor" />

			<input
				type="text"
				class="input input-bordered max-w-96 flex-1"
				placeholder="Search users..."
				name="search"
				defaultValue={defaultSearch}
			/>
			<button class="btn btn-primary">
				<Search class="size-4" />
				Search
			</button>
		</form>
		<div>
			<table class="table">
				<thead>
					<tr>
						<th>Name</th>
						<th>Email</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{#each users.data as user}
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
									<div tabindex="0" role="button" class="btn btn-ghost btn-sm btn-square">
										<EllipsisVertical class="size-4" />
									</div>
									<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
									<ul
										tabindex="0"
										class="dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm"
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
												Impersonate
											</button>

											<form method="POST" action="?/impersonate" id={`impersonate-form-${user.id}`}>
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
						<button class="btn btn-ghost btn-neutral btn-sm" onclick={loadPreviousPage}>
							<ChevronLeft class="size-4" />
							Previous
						</button>
					{/if}
					{#if users.links.next}
						<button class="btn btn-ghost btn-neutral btn-sm" onclick={loadNextPage}>
							Next
							<ChevronRight class="size-4" />
						</button>
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>
