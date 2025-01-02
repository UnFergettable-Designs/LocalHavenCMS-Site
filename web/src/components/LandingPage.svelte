<script lang="ts">
  import { Clock, Users, Cloud, Code, ArrowRight, CheckCircle } from 'lucide-svelte';
  import type { FormField, SurveyResponse } from '../types/Survey';
  import type { ComponentType } from 'svelte';
  import LocalHavenLogo from '../assets/LocalHavenCMS.webp';
  import SurveyForm from './SurveyForm.svelte';
  import { config } from '../config';

  export let formFields: FormField[];

  let currentStep = 0;
  let isSubmitting = false;
  let errorMessage = '';
  let formData: Partial<SurveyResponse> = {};

  interface Feature {
    icon: ComponentType;
    title: string;
    description: string;
  }

  const features: Feature[] = [
    {
      icon: Clock,
      title: 'Work Anywhere',
      description: 'True offline-first architecture ensures uninterrupted work',
    },
    {
      icon: Users,
      title: 'Real-Time Collaboration',
      description: 'Seamless team coordination across devices',
    },
    {
      icon: Cloud,
      title: 'Asset Management',
      description: 'Efficient handling of all your digital content',
    },
    {
      icon: Code,
      title: 'Developer Friendly',
      description: 'Built with modern tech stack for maximum flexibility',
    },
  ];

  async function handleSubmit(data: SurveyResponse): Promise<void> {
    try {
      isSubmitting = true;
      errorMessage = '';
      formData = data;

      if (!config.apiUrl) {
        throw new Error('API URL is not configured');
      }

      const response = await fetch(`${config.apiUrl}/survey`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        mode: 'cors',
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(
          (await response.text()) || `Server responded with status: ${response.status}`
        );
      }

      currentStep = 5; // Move to success step
    } catch (error: unknown) {
      errorMessage =
        error instanceof Error ? error.message : 'Failed to submit survey. Please try again.';
      throw error;
    } finally {
      isSubmitting = false;
    }
  }

  function scrollToSurvey(): void {
    document.getElementById('survey')?.scrollIntoView({ behavior: 'smooth' });
  }
</script>

<section>
  <div class="min-h-screen">
    <header class="header">
      <img src={LocalHavenLogo.src} alt="LocalHaven CMS Logo" class="logo" />
      <h1 class="visually-hidden">LocalHaven CMS</h1>
    </header>

    <div class="hero">
      <div class="hero-content">
        <h2 class="hero-title">Your Content. Your Control.</h2>
        <p class="hero-description">
          Experience the freedom of local-first content management. Collaborate, create, and manage
          your digital assets with ease—even offline.
        </p>
        <button class="button button-primary" on:click={scrollToSurvey}> Get Started </button>
      </div>
    </div>

    <div class="container">
      <div class="features-grid">
        {#each features as feature}
          <div class="feature-card">
            <div class="feature-icon">
              <svelte:component this={feature.icon} size={48} />
            </div>
            <h3 class="feature-title">{feature.title}</h3>
            <p class="feature-description">{feature.description}</p>
          </div>
        {/each}
      </div>
    </div>

    <div id="survey" class="container">
      <div class="survey-section">
        <div class="survey-content">
          <h2 class="survey-title">Help Shape the Future of LocalHaven CMS</h2>

          {#if currentStep < 6}
            <div class="step-container">
              {#each [0, 1, 2, 3, 4, 5] as step}
                <div
                  class="step-indicator {step <= currentStep ? 'step-active' : 'step-inactive'}"
                ></div>
              {/each}
            </div>

            <SurveyForm {formFields} bind:currentStep {isSubmitting} onSubmit={handleSubmit} />

            <!-- <div class="nav-buttons">
              {#if currentStep > 0}
                <button
                  class="button-secondary"
                  on:click={() => currentStep--}
                  disabled={isSubmitting}
                >
                  Back
                </button>
              {/if}
              <button
                class="button-primary"
                on:click={() => (currentStep < 4 ? currentStep++ : null)}
                disabled={isSubmitting}
              >
                {#if isSubmitting}
                  Submitting...
                {:else}
                  {currentStep === 4 ? 'Submit' : 'Next'}
                  <svelte:component this={ArrowRight} size={16} class="icon-right" />
                {/if}
              </button>
            </div> -->

            {#if errorMessage}
              <div class="error-message">{errorMessage}</div>
            {/if}
          {:else}
            <div class="success-message">
              <svelte:component this={CheckCircle} size={64} class="success-icon" />
              <h3 class="success-title">Thank You!</h3>
              <p class="success-text">Your feedback will help us build a better LocalHaven CMS.</p>
              {#if formData.betaInterest}
                <p class="beta-text">We'll be in touch about the beta program soon!</p>
              {/if}
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>

  <footer class="footer">
    <p>&copy; 2024 LocalHaven CMS. All rights reserved.</p>
  </footer>
</section>

<style>
  .min-h-screen {
    min-height: 100vh;
  }

  .header {
    display: flex;
    align-items: center;
    padding: 1rem;
    background-color: #111827;
    color: white;
  }

  .logo {
    height: 8rem;
    margin-right: 1rem;
  }

  .hero {
    background: linear-gradient(180deg, #065f46 0%, #047857 100%);
    color: white;
    padding: 4rem 1rem;
  }

  .hero-content {
    max-width: 48rem;
    margin: 0 auto;
    text-align: center;
  }

  .hero-title {
    font-size: 2.25rem;
    font-weight: 700;
    margin-bottom: 1.5rem;
  }

  .hero-description {
    font-size: 1.25rem;
    margin-bottom: 2rem;
  }

  .button {
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    font-weight: 500;
    cursor: pointer;
  }

  .button-primary {
    background-color: white;
    color: #047857;
    border: none;
  }

  .button-primary:hover {
    background-color: #f9fafb;
  }

  .success-message {
    text-align: center;
    padding: 2rem 0;
  }

  .success-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }

  .success-text {
    color: #6b7280;
  }

  .beta-text {
    margin-top: 1rem;
    color: #047857;
  }

  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 4rem 1rem;
  }

  .features-grid {
    display: grid;
    gap: 2rem;
  }

  @media (min-width: 768px) {
    .features-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (min-width: 1024px) {
    .features-grid {
      grid-template-columns: repeat(4, 1fr);
    }
  }

  .feature-card {
    background-color: white;
    border-radius: 0.5rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 1.5rem;
    text-align: center;
  }

  .feature-icon {
    display: flex;
    justify-content: center;
    margin-bottom: 1rem;
    color: #059669; /* green-600 */
  }

  .feature-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }

  .feature-description {
    color: #4b5563; /* gray-600 */
  }

  .survey-section {
    max-width: 42rem;
    margin: 0 auto;
    background-color: white;
    border-radius: 0.5rem;
    box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
  }

  .survey-content {
    padding: 1.5rem;
  }

  .survey-title {
    font-size: 1.5rem;
    font-weight: 700;
    text-align: center;
    margin-bottom: 1.5rem;
  }

  .step-container {
    display: flex;
    justify-content: space-between;
    margin-bottom: 1rem;
  }

  .step-indicator {
    height: 0.5rem;
    flex: 1;
    margin: 0 0.25rem;
    border-radius: 0.25rem;
  }

  .step-active {
    background-color: #059669; /* green-600 */
  }

  .step-inactive {
    background-color: #e5e7eb; /* gray-200 */
  }

  .visually-hidden {
    border: 0;
    clip: rect(0 0 0 0);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    width: 1px;
    white-space: nowrap;
  }

  .error-message {
    color: #dc2626;
    text-align: center;
    padding: 0.5rem;
    border-radius: 0.25rem;
    background-color: #fee2e2;
  }

  button:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
</style>
