<script lang="ts">
	import { CircleQuestionMark } from 'lucide-svelte';
	import { t } from '$lib/i18n';

	function handleAccordionClick(event: KeyboardEvent | MouseEvent) {
		// Find the radio input within the clicked accordion
		const radioInput = (event.currentTarget as HTMLElement).querySelector(
			'input[type="radio"]'
		) as HTMLInputElement;
		if (radioInput) {
			radioInput.checked = true;
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		// Handle Enter and Space key presses for accessibility
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			handleAccordionClick(event);
		}
	}
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
				<div class="rounded-xl bg-info/10 p-2">
					<CircleQuestionMark class="size-12 text-info" />
				</div>
			</div>
			<h1 class="mb-6 text-3xl font-bold md:text-4xl">{$t('faq.header.title')}</h1>
			<p class="mx-auto max-w-2xl text-xl leading-relaxed text-base-content/70">
				{$t('faq.header.subtitle')}
			</p>
		</div>

		<!-- FAQ Accordion -->
		<div class="space-y-3">
			<!-- Tech Stack Questions -->
			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why Svelte over React section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.whySvelte.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whySvelte.answer')}
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why DaisyUI for styling section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.whyDaisyui.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whyDaisyui.answer')}
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why Go for the backend section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.whyGo.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whyGo.answer')}
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why PostgreSQL section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.whyPostgresql.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whyPostgresql.answer')}
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why Railway for deployment section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.whyRailway.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whyRailway.answer')}
					</p>
				</div>
			</div>

			<!-- Architecture Questions -->
			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why custom authentication instead of Auth0 or Clerk section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					{$t('faq.questions.whyCustomAuth.question')}
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whyCustomAuth.answer')}
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle What is the organization-member-user model section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					{$t('faq.questions.whatIsOrgModel.question')}
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whatIsOrgModel.answer')}
					</p>
				</div>
			</div>

			<!-- Development Questions -->
			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why use a Makefile for development section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.whyMakefile.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.whyMakefile.answer')}
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle How do I get started with development section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.howToStart.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						<!-- eslint-disable-next-line svelte/no-at-html-tags -->
						{@html $t('faq.questions.howToStart.answer')}
					</p>
				</div>
			</div>

			<!-- Customization Questions -->
			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Can I use this without the organization features section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					{$t('faq.questions.withoutOrgs.question')}
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.withoutOrgs.answer')}
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle How do I add new OAuth providers section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">{$t('faq.questions.addOauth.question')}</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.addOauth.answer')}
					</p>
				</div>
			</div>

			<!-- Deployment Questions -->
			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Can I deploy this to other platforms besides Railway section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					{$t('faq.questions.otherPlatforms.question')}
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.otherPlatforms.answer')}
					</p>
				</div>
			</div>

			<!-- Billing & Business Questions -->
			<div
				class="collapse-arrow collapse border border-base-300 bg-base-200"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Is billing/payments functionality included section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					{$t('faq.questions.billingIncluded.question')}
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						{$t('faq.questions.billingIncluded.answer')}
					</p>
				</div>
			</div>
		</div>
	</div>
</div>
