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
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Link } from '$lib/components/ui/link';
	import { t } from '$lib/i18n';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';

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

<Card.Root>
	<Card.Header>
		<Card.Title>{$t('noOrganization.admin.users.title')}</Card.Title>
	</Card.Header>
	<Card.Content>
		<form class="flex gap-2" method="GET" id="search-form">
			<input type="hidden" name="page[cursor]" value={users.links.next ?? ''} id="page-cursor" />

			<Input
				type="search"
				class="max-w-96 flex-1"
				placeholder={$t('noOrganization.admin.users.searchPlaceholder')}
				name="search"
				defaultValue={defaultSearch}
			/>
			<Button type="submit">
				<Search class="size-4" />
				{$t('noOrganization.admin.users.search')}
			</Button>
		</form>
		<div>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>{$t('noOrganization.admin.users.name')}</Table.Head>
						<Table.Head>{$t('noOrganization.admin.users.email')}</Table.Head>
						<Table.Head class="text-right"></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each users.data as user (user.id)}
						<Table.Row>
							<Table.Cell>
								<div class="flex items-start gap-3">
									{#if user.meta.logoDistributionUrl}
										<img
											src={user.meta.logoDistributionUrl}
											alt={user.attributes.name}
											class="size-8 rounded-md"
										/>
									{:else}
										<div class="flex size-8 items-center justify-center rounded-md bg-muted">
											<User class="size-4" />
										</div>
									{/if}
									<span class="font-semibold">
										{user.attributes.name}
									</span>
								</div>
							</Table.Cell>
							<Table.Cell>
								<Link href={`mailto:${user.attributes.email}`}>
									{user.attributes.email}
								</Link>
							</Table.Cell>
							<Table.Cell class="text-right">
								<DropdownMenu.Root>
									<DropdownMenu.Trigger>
										{#snippet child({ props })}
											<Button variant="ghost" size="sm" {...props}>
												<EllipsisVertical class="size-4" />
											</Button>
										{/snippet}
									</DropdownMenu.Trigger>
									<DropdownMenu.Content align="end" class="w-52">
										<DropdownMenu.Item
											onSelect={(e) => {
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
										</DropdownMenu.Item>

										<form
											method="POST"
											action="?/impersonate"
											id={`impersonate-form-${user.id}`}
											hidden
										>
											<input type="hidden" name="impersonatedUserId" value={user.id} />
										</form>
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
			{#if users.links.prev || users.links.next}
				<div class="mt-4 flex justify-center gap-2">
					{#if users.links.prev}
						<Button size="sm" variant="ghost" onclick={loadPreviousPage}>
							<ChevronLeft class="size-4" />
							{$t('noOrganization.admin.users.previous')}
						</Button>
					{/if}
					{#if users.links.next}
						<Button size="sm" variant="ghost" onclick={loadNextPage}>
							{$t('noOrganization.admin.users.next')}
							<ChevronRight class="size-4" />
						</Button>
					{/if}
				</div>
			{/if}
		</div>
	</Card.Content>
</Card.Root>
