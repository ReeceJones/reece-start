<script lang="ts">
	import { DoorOpen, Rocket, Menu, XIcon } from 'lucide-svelte';
	import * as NavigationMenu from '$lib/components/ui/navigation-menu';
	import * as Sheet from '$lib/components/ui/sheet';
	import { Button } from '$lib/components/ui/button';
	import * as m from '$lib/paraglide/messages';
	import { env } from '$env/dynamic/public';
	import { cn } from '$lib/utils';
	import { buttonVariants } from '../ui/button';

	const { isLoggedIn }: { isLoggedIn: boolean } = $props();

	const isSignInDisabled = $derived(env.PUBLIC_DISABLE_SIGNIN === 'true');
</script>

<div class="container mx-auto mt-1 flex">
	<div class="flex-1">
		<NavigationMenu.Root>
			<NavigationMenu.List>
				<NavigationMenu.Item>
					<a
						class={cn(
							buttonVariants({
								variant: 'ghost',
								size: 'lg'
							}),
							'text-lg tracking-tight'
						)}
						href="/"
					>
						<Rocket class="size-6" />
						reece-start
					</a>
				</NavigationMenu.Item>
			</NavigationMenu.List>
		</NavigationMenu.Root>
	</div>

	<!-- Desktop Navigation -->
	<NavigationMenu.Root class="hidden lg:flex">
		<NavigationMenu.List>
			<NavigationMenu.Item>
				<a
					href="/faq"
					class={buttonVariants({
						variant: 'ghost'
					})}>{m.footer__faq()}</a
				>
			</NavigationMenu.Item>
			<NavigationMenu.Item>
				<a
					href="/pricing"
					class={buttonVariants({
						variant: 'ghost'
					})}>{m.footer__pricing()}</a
				>
			</NavigationMenu.Item>

			{#if isLoggedIn}
				<a
					href="/app"
					class={buttonVariants({
						variant: 'default'
					})}
				>
					<DoorOpen class="size-5" />
					{m.dashboard()}
				</a>
			{:else if !isSignInDisabled}
				<a
					href="/signin"
					class={buttonVariants({
						variant: 'outline'
					})}
				>
					{m.sign_in()}
				</a>
				<a href="/signup" class={buttonVariants({ variant: 'default' })}> {m.get_started()} </a>
			{/if}
		</NavigationMenu.List>
	</NavigationMenu.Root>

	<!-- Mobile Hamburger Menu -->
	<div class="lg:hidden">
		<Sheet.Root>
			<Sheet.Trigger>
				<Button variant="ghost" aria-label="Toggle mobile menu">
					<Menu class="size-6" />
				</Button>
			</Sheet.Trigger>
			<Sheet.Content showClose={false}>
				<div class="space-y-3 px-4 py-6">
					<a
						href="/faq"
						class={cn(buttonVariants({ variant: 'ghost' }), 'w-full justify-end text-lg')}
					>
						{m.footer__faq()}
					</a>
					<a
						href="/pricing"
						class={cn(buttonVariants({ variant: 'ghost' }), 'w-full justify-end text-lg')}
					>
						{m.footer__pricing()}
					</a>
					{#if isLoggedIn}
						<a
							href="/app"
							class={cn(buttonVariants({ variant: 'ghost' }), 'w-full justify-end text-lg')}
						>
							<DoorOpen class="mr-2 inline size-5" />
							{m.dashboard()}
						</a>
					{:else}
						<a
							href="/signin"
							class={cn(buttonVariants({ variant: 'ghost' }), 'w-full justify-end text-lg')}
						>
							{m.sign_in()}
						</a>
						<a
							href="/signup"
							class={cn(buttonVariants({ variant: 'ghost' }), 'w-full justify-end text-lg')}
						>
							{m.get_started()}
						</a>
					{/if}
					<Sheet.Close class="w-full">
						<Button variant="secondary" class="w-full justify-end text-lg">
							<XIcon class="size-6" />
							Close
						</Button>
					</Sheet.Close>
				</div>
			</Sheet.Content>
		</Sheet.Root>
	</div>
</div>
