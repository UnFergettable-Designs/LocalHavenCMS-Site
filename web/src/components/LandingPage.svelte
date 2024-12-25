<script lang="ts">
  import { Clock, Users, Cloud, Code, ArrowRight, CheckCircle } from 'lucide-svelte';
  import LocalHavenLogo from '../assets/LocalHavenCMS.webp';
  import { config } from '../config';

  interface FormData {
    role: string;
    otherRole: string;
    cmsUsage: string;
    otherCmsUsage: string;
    features: {
      [key: string]: number;
      offline: number;
      collaboration: number;
      assetManagement: number;
      pdfHandling: number;
      versionControl: number;
      workflows: number;
    };
    betaInterest: boolean | null;
    email: string;
  }

  let currentStep = 0;
  let formData: FormData = {
    role: '',
    otherRole: '',
    cmsUsage: '',
    otherCmsUsage: '',
    features: {
      offline: 0,
      collaboration: 0,
      assetManagement: 0,
      pdfHandling: 0,
      versionControl: 0,
      workflows: 0,
    },
    betaInterest: null,
    email: '',
  };

  const features = [
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

  const featureLabels = {
    offline: 'Offline Capabilities',
    collaboration: 'Real-time Collaboration',
    assetManagement: 'Asset Management',
    pdfHandling: 'PDF Handling',
    versionControl: 'Version Control',
    workflows: 'Approval Workflows',
  };

  async function handleSubmit() {
    try {
      console.log('API URL:', config.apiUrl);
      console.log('Submitting form data:', formData);

      const response = await fetch(`${config.apiUrl}/survey`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
        credentials: 'include',
        mode: 'cors',
        body: JSON.stringify({
          role: formData.role,
          otherRole: formData.otherRole,
          cmsUsage: formData.cmsUsage,
          otherCmsUsage: formData.otherCmsUsage,
          features: formData.features,
          betaInterest: formData.betaInterest ?? false, // Ensure boolean
          email: formData.email,
        }),
      });

      console.log('Response status:', response.status);

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Error response:', errorText);
        throw new Error(errorText || 'Failed to submit survey');
      }

      const result = await response.json();
      console.log('Success response:', result);
      currentStep = 3;
    } catch (error: any) {
      console.error('Error details:', error);
      alert(`Failed to submit survey: ${error.message}`);
    }
  }

  function scrollToSurvey() {
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

          {#if currentStep < 3}
            <div class="step-container">
              {#each [0, 1, 2] as step}
                <div
                  class="step-indicator {step <= currentStep ? 'step-active' : 'step-inactive'}"
                ></div>
              {/each}
            </div>

            {#if currentStep === 0}
              <div class="space-y-4">
                <div>
                  <label class="label"
                    >What is your primary role?
                    <select class="role-select" bind:value={formData.role}>
                      <option value="marketer">Marketer</option>
                      <option value="developer">Developer</option>
                      <option value="designer">Designer</option>
                      <option value="content">Content</option>
                      <option value="other">Other</option>
                    </select>
                  </label>
                  {#if formData.role === 'other'}
                    <div class="mt-2">
                      <input
                        type="text"
                        placeholder="Please specify"
                        class="input-text"
                        bind:value={formData.otherRole}
                      />
                    </div>
                  {/if}
                </div>

                <div>
                  <label class="label"
                    >How often do you use a CMS?
                    <select class="cms-usage-select" bind:value={formData.cmsUsage}>
                      <option value="">Select frequency...</option>
                      <option value="daily">Daily</option>
                      <option value="weekly">Weekly</option>
                      <option value="monthly">Monthly</option>
                      <option value="rarely">Rarely</option>
                      <option value="other">Other</option>
                    </select>
                  </label>
                  {#if formData.cmsUsage === 'other'}
                    <div class="mt-2">
                      <input
                        type="text"
                        placeholder="Please specify"
                        class="input-text"
                        bind:value={formData.otherCmsUsage}
                      />
                    </div>
                  {/if}
                </div>
              </div>
            {:else if currentStep === 1}
              <div class="space-y-4">
                {#each Object.entries(featureLabels) as [feature, label]}
                  <div>
                    <label class="label">
                      {label}
                      <div class="flex gap-2">
                        {#each [1, 2, 3, 4, 5] as value}
                          {@const currentFeature = formData.features[feature]}
                          <button
                            class="feature-button"
                            class:active={currentFeature === value}
                            on:click={() => (formData.features[feature] = value)}
                          >
                            {value}
                          </button>
                        {/each}
                      </div>
                    </label>
                  </div>
                {/each}
              </div>
            {:else}
              <div class="space-y-4">
                <div>
                  <label class="label"
                    >Would you like to join our beta program?
                    <div class="radio-group">
                      <input type="radio" value={true} bind:group={formData.betaInterest} /> Yes
                      <input type="radio" value={false} bind:group={formData.betaInterest} /> No
                    </div>
                  </label>
                </div>
                {#if formData.betaInterest}
                  <div>
                    <label class="label"
                      >Email Address
                      <input
                        type="email"
                        class="input-text"
                        bind:value={formData.email}
                        placeholder="Enter your email"
                      />
                    </label>
                  </div>
                {/if}
              </div>
            {/if}

            <div class="flex justify-between mt-6">
              {#if currentStep > 0}
                <button class="button-secondary" on:click={() => currentStep--}> Back </button>
              {/if}
              <button
                class="button-primary"
                on:click={() => {
                  if (currentStep === 2) {
                    handleSubmit();
                  } else {
                    currentStep++;
                  }
                }}
              >
                {currentStep === 2 ? 'Submit' : 'Next'}
                <ArrowRight size={16} class="ml-2" />
              </button>
            </div>
          {:else}
            <div class="text-center section-padding">
              <CheckCircle size={64} class="success-icon" />
              <h3 class="text-2xl font-semibold mb-2">Thank You!</h3>
              <p class="gray-text">Your feedback will help us build a better LocalHaven CMS.</p>
              {#if formData.betaInterest}
                <p class="mt-4 success-text">We'll be in touch about the beta program soon!</p>
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

  .button-secondary {
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    font-weight: 500;
    border: 1px solid #d1d5db;
    background-color: white;
    color: #374151;
    cursor: pointer;
  }

  .button-secondary:hover {
    background-color: #f3f4f6;
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

  .label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    margin-bottom: 0.5rem;
  }

  .role-select,
  .cms-usage-select,
  .input-text {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #d1d5db;
    border-radius: 0.25rem;
  }

  .feature-button {
    padding: 0.5rem;
    width: 2rem;
    border: 1px solid #d1d5db;
    border-radius: 0.25rem;
    cursor: pointer;
  }

  .feature-button.active {
    background-color: #059669;
    color: white;
  }

  .radio-group {
    display: flex;
    gap: 1rem;
  }

  .footer {
    text-align: center;
    padding: 1rem;
    background-color: #f9fafb;
    color: #6b7280;
  }

  .space-y-4 > * + * {
    margin-top: 1rem;
  }

  .mt-2 {
    margin-top: 0.5rem;
  }

  .mt-4 {
    margin-top: 1rem;
  }

  .mt-6 {
    margin-top: 1.5rem;
  }

  .mb-2 {
    margin-bottom: 0.5rem;
  }

  .mb-4 {
    margin-bottom: 1rem;
  }

  .ml-2 {
    margin-left: 0.5rem;
  }

  .flex {
    display: flex;
  }

  .gap-2 {
    gap: 0.5rem;
  }

  .justify-between {
    justify-content: space-between;
  }

  .text-center {
    text-align: center;
  }

  .text-2xl {
    font-size: 1.5rem;
  }

  .font-semibold {
    font-weight: 600;
  }

  .success-icon {
    color: #059669;
    margin: 0 auto 1rem auto;
  }

  .success-text {
    color: #059669;
  }

  .gray-text {
    color: #6b7280;
  }

  .section-padding {
    padding: 2rem 0;
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
</style>
