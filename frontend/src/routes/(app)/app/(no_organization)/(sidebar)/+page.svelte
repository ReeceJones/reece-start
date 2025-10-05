<script lang="ts">
	import { ChevronRight, Plus, Users } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import CardActions from '$lib/components/Card/CardActions.svelte';

	const { data }: PageProps = $props();
	const { organizations } = data;
</script>

<Card>
	<CardBody>
		<CardTitle><Users /> Organizations</CardTitle>
		<p>Select an organization to continue to the app.</p>

		<div class="space-y-2">
			<div class="divider"></div>
			<ul class="my-4 grid grid-cols-3 gap-4">
				{#if organizations.data.length === 0}
					<li class="col-span-full pb-3">
						You are not a member of any organizations. Create or join an organization to get
						started.
					</li>
				{/if}
				{#each organizations.data as organization (organization.id)}
					<li>
						<a
							href={`/app/${organization.id}`}
							class="flex w-full items-start justify-between rounded-lg border border-base-300 bg-base-100 p-3 overflow-ellipsis"
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
										<p class="text-transparent">No description</p>
									{/if}
								</div>
							</div>
							<ChevronRight />
						</a>
					</li>
				{/each}
			</ul>
			<CardActions>
				<a href="/app/create-organization/basic-information" class="btn btn-primary"
					><Plus /> Create Organization</a
				>
			</CardActions>
		</div>
	</CardBody>
</Card>
