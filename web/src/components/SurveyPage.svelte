<script lang="ts">
  import type { FormField, SurveyResponse } from '../types/Survey';
  import SurveyForm from './SurveyForm.svelte';
  import { config } from '../config';

  export let formFields: FormField[] = [];
  let isSubmitting = false;
  let currentStep = 0;
  let errorMessage = '';
  let formData: Partial<SurveyResponse> = {};

  async function handleSubmit(data: SurveyResponse): Promise<void> {
    isSubmitting = true;
    errorMessage = '';
    formData = data;

    try {
      const response = await fetch(`${config.apiUrl}/survey`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error('Failed to submit survey');
      }

      currentStep = 6; // Show success message
    } catch (error) {
      console.error('Survey submission error:', error);
      errorMessage = 'Failed to submit survey. Please try again.';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="survey-container">
  <h1>Help Shape the Future of LocalHaven CMS</h1>

  {#if currentStep < 6}
    <div class="step-container">
      {#each [0, 1, 2, 3, 4, 5] as step}
        <div class="step-indicator {step <= currentStep ? 'step-active' : 'step-inactive'}"></div>
      {/each}
    </div>

    <SurveyForm {formFields} bind:currentStep {isSubmitting} onSubmit={handleSubmit} />

    {#if errorMessage}
      <div class="error-message">{errorMessage}</div>
    {/if}
  {:else}
    <div class="success-message">
      <h2>Thank You!</h2>
      <p>Your feedback will help us build a better LocalHaven CMS.</p>
      {#if formData.betaInterest}
        <p class="beta-text">We'll be in touch about the beta program soon!</p>
      {/if}
    </div>
  {/if}
</div>

<style>
  .survey-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem 1rem;
  }

  .step-container {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    margin: 2rem 0;
  }

  .step-indicator {
    width: 1rem;
    height: 1rem;
    border-radius: 50%;
    background-color: var(--color-background-light);
  }

  .step-active {
    background-color: var(--color-primary);
  }

  .success-message {
    text-align: center;
    padding: 2rem;
  }

  .error-message {
    color: var(--color-error);
    text-align: center;
    margin-top: 1rem;
  }

  .beta-text {
    color: var(--color-primary);
    margin-top: 1rem;
  }
</style>
