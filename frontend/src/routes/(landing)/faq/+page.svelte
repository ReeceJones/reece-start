<script lang="ts">
	import { CircleQuestionMark } from 'lucide-svelte';
	import * as Accordion from '$lib/components/ui/accordion';
	import { t } from '$lib/i18n';

	const faqs = $derived([
		{
			question: $t('faq.questions.whySvelte.question'),
			answer: $t('faq.questions.whySvelte.answer')
		},
		{
			question: $t('faq.questions.whyDaisyui.question'),
			answer: $t('faq.questions.whyDaisyui.answer')
		},
		{
			question: $t('faq.questions.whyGo.question'),
			answer: $t('faq.questions.whyGo.answer')
		},
		{
			question: $t('faq.questions.whyPostgresql.question'),
			answer: $t('faq.questions.whyPostgresql.answer')
		},
		{
			question: $t('faq.questions.whyRailway.question'),
			answer: $t('faq.questions.whyRailway.answer')
		},
		{
			question: $t('faq.questions.whyCustomAuth.question'),
			answer: $t('faq.questions.whyCustomAuth.answer')
		},
		{
			question: $t('faq.questions.whatIsOrgModel.question'),
			answer: $t('faq.questions.whatIsOrgModel.answer')
		},
		{
			question: $t('faq.questions.whyMakefile.question'),
			answer: $t('faq.questions.whyMakefile.answer')
		},
		{
			question: $t('faq.questions.howToStart.question'),
			answer: $t('faq.questions.howToStart.answer'),
			html: true
		},
		{
			question: $t('faq.questions.withoutOrgs.question'),
			answer: $t('faq.questions.withoutOrgs.answer')
		},
		{
			question: $t('faq.questions.addOauth.question'),
			answer: $t('faq.questions.addOauth.answer')
		},
		{
			question: $t('faq.questions.otherPlatforms.question'),
			answer: $t('faq.questions.otherPlatforms.answer')
		},
		{
			question: $t('faq.questions.billingIncluded.question'),
			answer: $t('faq.questions.billingIncluded.answer')
		}
	]);
</script>

<svelte:head>
	<title>{$t('faq.meta.title')}</title>
	<meta name="description" content={$t('faq.meta.description')} />
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
			<h1 class="mb-6 text-2xl font-bold md:text-4xl">{$t('faq.header.title')}</h1>
			<p class="text-base-content/70 mx-auto max-w-2xl text-xl leading-relaxed">
				{$t('faq.header.subtitle')}
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
