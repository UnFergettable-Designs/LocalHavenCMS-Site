<script lang="ts">
  import { CheckCircle, Laptop, Users, Folder, Code } from 'lucide-svelte';
  import type { FormField, SurveyResponse, Feature } from '../types/Survey';
  import LocalHavenLogo from '../assets/LocalHavenCMS.webp';
  import SurveyForm from './SurveyForm.svelte';
  import { config } from '../config';

  export let formFields: FormField[] = [];
  export let isSubmitting = false;
  export let currentStep = 0;
  let errorMessage = '';
  let formData: Partial<SurveyResponse> = {};

  const features: Feature[] = [
    {
      icon: Laptop,
      title: 'Work Anywhere',
      description: 'True offline-first architecture ensures uninterrupted work',
    },
    {
      icon: Users,
      title: 'Real-Time Collaboration',
      description: 'Seamless team coordination across devices',
    },
    {
      icon: Folder,
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

    <div class="container">
      <div class="features-grid">
        {#each features as feature}
          <div class="feature-card">
            <svelte:component this={feature.icon} size={32} class="feature-icon" />
            <h3>{feature.title}</h3>
            <p>{feature.description}</p>
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

<!-- Main Landing Page Layout -->
<section class="landing-page">
  <!-- Existing content / survey embed -->
  <!-- The SurveyForm or whichever component you display next -->
  <!-- Example: <SurveyForm ... /> or other content -->

  {#if errorMessage}
    <div class="error">{errorMessage}</div>
  {/if}
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

  svg > path {
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
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .feature-icon {
    margin-bottom: 1rem;
    color: #047857; /* green-600 */
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

  .landing-page {
    max-width: 800px;
    margin: 0 auto;
    padding: 1rem;
  }
  .value-prop {
    background: #f9fafb;
    padding: 1.5rem;
    border-radius: 0.5rem;
    margin-bottom: 2rem;
  }
  .value-prop h1 {
    margin-bottom: 1rem;
    font-size: 1.5rem;
    color: #111827;
  }
  .value-prop p {
    line-height: 1.6;
    color: #374151;
  }
  .error {
    color: #b91c1c;
    margin-top: 1rem;
  }

  .hero {
    padding: 4rem 2rem;
    text-align: center;
    background: linear-gradient(180deg, #047857 0%, #065f46 100%);
    color: white;
  }

  .lead {
    font-size: 1.5rem;
    max-width: 36rem;
    margin: 1rem auto 3rem;
    opacity: 0.9;
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .feature-card {
    background: rgba(255, 255, 255, 0.1);
    padding: 2rem;
    border-radius: 0.5rem;
    transition: transform 0.2s;
  }

  .feature-card:hover {
    transform: translateY(-4px);
  }

  .icon {
    font-size: 2rem;
    margin-bottom: 1rem;
    display: block;
  }
</style>
