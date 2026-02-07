<script lang="ts">
	import { CircleQuestionMark } from 'lucide-svelte';
	import * as Accordion from '$lib/components/ui/accordion';
	import * as m from '$lib/paraglide/messages';

	const faqs = $derived([
		{
			question: m.faq__questions__why_svelte__question(),
			answer: m.faq__questions__why_svelte__answer()
		},
		{
			question: m.faq__questions__why_daisyui__question(),
			answer: m.faq__questions__why_daisyui__answer()
		},
		{
			question: m.faq__questions__why_go__question(),
			answer: m.faq__questions__why_go__answer()
		},
		{
			question: m.faq__questions__why_postgresql__question(),
			answer: m.faq__questions__why_postgresql__answer()
		},
		{
			question: m.faq__questions__why_railway__question(),
			answer: m.faq__questions__why_railway__answer()
		},
		{
			question: m.faq__questions__why_custom_auth__question(),
			answer: m.faq__questions__why_custom_auth__answer()
		},
		{
			question: m.faq__questions__what_is_org_model__question(),
			answer: m.faq__questions__what_is_org_model__answer()
		},
		{
			question: m.faq__questions__why_makefile__question(),
			answer: m.faq__questions__why_makefile__answer()
		},
		{
			question: m.faq__questions__how_to_start__question(),
			answer: m.faq__questions__how_to_start__answer(),
			html: true
		},
		{
			question: m.faq__questions__without_orgs__question(),
			answer: m.faq__questions__without_orgs__answer()
		},
		{
			question: m.faq__questions__add_oauth__question(),
			answer: m.faq__questions__add_oauth__answer()
		},
		{
			question: m.faq__questions__other_platforms__question(),
			answer: m.faq__questions__other_platforms__answer()
		},
		{
			question: m.faq__questions__billing_included__question(),
			answer: m.faq__questions__billing_included__answer()
		}
	]);
</script>

<svelte:head>
	<title>{m.faq__meta__title()}</title>
	<meta name="description" content={m.faq__meta__description()} />
</svelte:head>

<!-- FAQ Header -->
<div class="container mx-auto px-4">
	<div class="mx-auto max-w-4xl py-24">
		<div class="mb-12 text-center">
			<div class="mb-6 flex justify-center">
				<div class="bg-info/10 rounded-xl p-2">
					<CircleQuestionMark class="text-info size-12" />
				</div>
			</div>
			<h1 class="mb-6 text-2xl font-bold md:text-4xl">{m.faq__header__title()}</h1>
			<p class="text-base-content/70 mx-auto max-w-2xl text-xl leading-relaxed">
				{m.faq__header__subtitle()}
			</p>
		</div>

		<!-- FAQ Accordion -->

		<Accordion.Root type="multiple">
			{#each faqs as faq (faq.question)}
				<Accordion.Item value={faq.question}>
					<Accordion.Trigger class="text-base font-semibold">{faq.question}</Accordion.Trigger>
					<Accordion.Content class="text-primary/90">
						{faq.answer}
					</Accordion.Content>
				</Accordion.Item>
			{/each}
		</Accordion.Root>
	</div>
</div>
