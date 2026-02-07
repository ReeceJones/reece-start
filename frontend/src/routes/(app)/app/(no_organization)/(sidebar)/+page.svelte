<script lang="ts">
	import { ChevronRight, Plus, Users } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card';
	import * as Empty from '$lib/components/ui/empty';
	import * as m from '$lib/paraglide/messages';
	import { buttonVariants } from '$lib/components/ui/button';

	const { data }: PageProps = $props();
	const { organizations } = $derived(data);
</script>

{#if organizations.data.length === 0}
	<Empty.Root>
		<Empty.Header>
			<Empty.Media variant="icon">
				<Users class="size-6" />
			</Empty.Media>
			<Empty.Title>{m.no_organization__no_organizations__title()}</Empty.Title>
			<Empty.Description>{m.no_organization__no_organizations__description()}</Empty.Description>
		</Empty.Header>
		<Empty.Content>
			<a href="/app/create-organization/basic-information" class={buttonVariants()}
				><Plus /> {m.no_organization__create_organization()}</a
			>
		</Empty.Content>
	</Empty.Root>
{:else}
	<Card.Root>
		<Card.Header>
			<Card.Title class="flex items-end gap-2">
				<Users class="size-6" />
				{m.no_organization__organizations()}
			</Card.Title>
			<Card.Description>
				{m.no_organization__select_organization()}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="space-y-2">
				<div class="divider"></div>
				<ul class="my-4 grid grid-cols-3 gap-4">
					{#each organizations.data as organization (organization.id)}
						<li>
							<a
								href={`/app/${organization.id}`}
								class="border-base-300 bg-base-100 flex w-full items-start justify-between rounded-lg border p-3 overflow-ellipsis"
							>
								<div class="flex items-center gap-3">
									{#if organization.meta.logoDistributionUrl}
										<img
											src={organization.meta.logoDistributionUrl}
											alt="Organization logo"
											class="size-12 rounded-lg"
										/>
									{/if}
									<div>
										<p class="line-clamp-1 text-lg font-medium overflow-ellipsis">
											{organization.attributes.name}
										</p>
										{#if organization.attributes.description}
											<p class="line-clamp-1 text-sm overflow-ellipsis text-gray-500">
												{organization.attributes.description}
											</p>
										{:else}
											<p class="text-transparent">{m.no_organization__no_description()}</p>
										{/if}
									</div>
								</div>
								<ChevronRight />
							</a>
						</li>
					{/each}
				</ul>
			</div>
		</Card.Content>
		<Card.Footer>
			<Card.Action>
				<a href="/app/create-organization/basic-information" class={buttonVariants()}
					><Plus /> {m.no_organization__create_organization()}</a
				>
			</Card.Action>
		</Card.Footer>
	</Card.Root>
{/if}
