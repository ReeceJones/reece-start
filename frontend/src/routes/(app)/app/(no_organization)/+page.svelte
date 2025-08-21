<script lang="ts">
	import { ChevronRight, Plus, Users } from 'lucide-svelte';
	import type { PageProps } from './$types';

	const { data }: PageProps = $props();
	const { organizations } = data;
</script>

<div class="mx-auto max-w-[600px] p-4">
	<div class="card bg-base-200 shadow-sm">
		<div class="card-body">
			<h2 class="card-title"><Users /> Organizations</h2>
			<p>Select an organization to continue to the app.</p>

			<div class="space-y-2">
				<div class="divider"></div>
				<ul class="my-4 space-y-3">
					{#if organizations.data.length === 0}
						<li class="pb-3">
							You are not a member of any organizations. Create or join an organization to get
							started.
						</li>
					{/if}
					{#each organizations.data as organization}
						<li class="w-full">
							<a
								href={`/app/${organization.id}`}
								class="bg-base-100 border-base-300 flex w-full justify-between rounded-lg border p-3"
							>
								<div>
									<p class="text-lg font-medium">{organization.attributes.name}</p>
									{#if organization.attributes.description}
										<p class="text-sm text-gray-500">{organization.attributes.description}</p>
									{/if}
								</div>
								<ChevronRight />
							</a>
						</li>
					{/each}
				</ul>
				<div class="card-actions justify-end">
					<a href="/app/create-organization" class="btn btn-primary"><Plus /> Create Organization</a
					>
				</div>
			</div>
		</div>
	</div>
</div>
