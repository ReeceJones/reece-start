<script lang="ts">
	import { ChevronRight, Plus, Users } from 'lucide-svelte';
	import type { PageProps } from './$types';

	const { data }: PageProps = $props();
	const { organizations } = data;
</script>

<div class="card bg-base-200 border-base-300 border shadow-sm">
	<div class="card-body">
		<h2 class="card-title"><Users /> Organizations</h2>
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
				{#each organizations.data as organization}
					<li>
						<a
							href={`/app/${organization.id}`}
							class="bg-base-100 border-base-300 flex w-full items-start justify-between overflow-ellipsis rounded-lg border p-3"
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
									<p class="line-clamp-1 overflow-ellipsis text-lg font-medium">
										{organization.attributes.name}
									</p>
									{#if organization.attributes.description}
										<p class="line-clamp-1 overflow-ellipsis text-sm text-gray-500">
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
			<div class="card-actions justify-end">
				<a href="/app/create-organization" class="btn btn-primary"><Plus /> Create Organization</a>
			</div>
		</div>
	</div>
</div>
