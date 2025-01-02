<script lang="ts">
  import type { SurveyResponse, Features, FormField } from '../types/Survey';

  export let currentStep: number;
  export let onSubmit: (data: SurveyResponse) => Promise<void>;
  export let isSubmitting = false;
  export let formFields: FormField[] = [];

  const FIELDS_PER_STEP = 5;
  const TOTAL_STEPS = Math.ceil(formFields.length / FIELDS_PER_STEP);

  type FormDataType = Partial<SurveyResponse> & {
    features: Partial<Features>;
    [key: string]: any; // Allow string indexing for dynamic field access
  };

  type FieldValue = string | number | boolean | undefined;

  // Initialize formData with undefined values for radio buttons
  let formData: FormDataType = {
    features: {
      offline: undefined,
      collaboration: undefined,
      assetManagement: undefined,
      pdfHandling: undefined,
      versionControl: undefined,
      workflows: undefined,
    },
    betaInterest: false,
    role: '',
    otherRole: '',
    cmsUsage: '',
    otherCmsUsage: '',
    email: '',
    biggestFrustrations: '',
    specificProblems: '',
    usageFrequency: '',
    primaryPurpose: '',
    platforms: '',
    cmsPreference: '',
    wishedFeatures: '',
    workflowImportance: '',
    teamSize: '',
    collaborationFrequency: '',
    pricingSensitivity: '',
    pricingModel: '',
    integrations: '',
    integrationImportance: '',
    contentTypes: '',
    customFormats: '',
  };

  let errors: Record<string, string> = {};
  let currentFields: FormField[] = [];

  // Add feature labels for display
  const featureLabels = {
    offline: 'Offline Capabilities',
    collaboration: 'Collaboration',
    assetManagement: 'Asset Management',
    pdfHandling: 'PDF Handling',
    versionControl: 'Version Control',
    workflows: 'Workflows',
  } as const;

  // Reactive statement to update current fields when step changes
  $: {
    const fieldsPerStep = 5;
    const startIndex = currentStep * fieldsPerStep;
    currentFields = formFields.slice(startIndex, startIndex + fieldsPerStep);
  }

  // Validate current step fields
  function validateCurrentStep(): boolean {
    errors = {};
    let isValid = true;

    currentFields.forEach((field) => {
      const fieldName = field.name;

      // Handle nested feature fields
      if (fieldName.startsWith('features.')) {
        const [, featureName] = fieldName.split('.') as [string, keyof Features];
        if (
          field.required &&
          (!formData.features || formData.features[featureName] === undefined)
        ) {
          errors[fieldName] = `${field.label} is required`;
          isValid = false;
        }
      } else {
        // Handle regular fields
        const value = formData[fieldName as keyof SurveyResponse];
        if (field.required && (value === undefined || value === '')) {
          errors[fieldName] = `${field.label} is required`;
          isValid = false;
        }
      }

      // Check dependent fields
      if (field.dependsOn && field.dependsOnValue) {
        const parentValue = formData[field.dependsOn as keyof SurveyResponse];
        const fieldValue = formData[fieldName as keyof SurveyResponse];
        if (parentValue === field.dependsOnValue && !fieldValue) {
          errors[fieldName] =
            `${field.label} is required when ${field.dependsOn} is ${field.dependsOnValue}`;
          isValid = false;
        }
      }
    });

    return isValid;
  }

  // Handle field changes with special handling for features
  function handleChange(event: Event, field: FormField): void {
    const target = event.target as HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement;
    const value = field.type === 'checkbox' ? (target as HTMLInputElement).checked : target.value;

    if (field.name.startsWith('features.')) {
      const [, featureName] = field.name.split('.') as [string, keyof Features];
      formData = {
        ...formData,
        features: {
          ...formData.features,
          [featureName]:
            field.type === 'checkbox' ? Boolean(value) : parseInt(value as string) || 0,
        },
      };
    } else {
      formData = {
        ...formData,
        [field.name]: value,
      } as FormDataType;
    }

    // Clear error when field is updated
    if (errors[field.name]) {
      errors = {
        ...errors,
        [field.name]: '',
      };
    }
  }

  let submissionStatus = '';
  let showThankYou = false;

  // Handle form submission
  async function handleSubmit(event: Event): Promise<void> {
    event.preventDefault();
    console.log('Form submission started', { currentStep, isLastStep: isLastStep() });

    if (!validateCurrentStep()) {
      console.log('Validation failed');
      return;
    }

    if (!isLastStep()) {
      console.log('Moving to next step');
      currentStep++;
      return;
    }

    if (isSubmitting) {
      console.log('Already submitting');
      return;
    }

    try {
      isSubmitting = true;
      submissionStatus = 'submitting';
      console.log('Submitting form data:', formData);

      await onSubmit(formData as SurveyResponse);
      submissionStatus = 'success';
      showThankYou = true;
      console.log('Submission successful');
    } catch (error) {
      console.error('Form submission error:', error);
      submissionStatus = 'error';
      throw error; // Propagate error to parent
    } finally {
      isSubmitting = false;
    }
  }

  function isFeatureField(fieldName: string): boolean {
    return fieldName.startsWith('features.');
  }

  function getFieldValue(field: FormField): FieldValue {
    if (isFeatureField(field.name)) {
      const [, featureName] = field.name.split('.') as [string, keyof Features];
      return formData.features[featureName];
    }
    return formData[field.name];
  }

  function setFieldValue(field: FormField, value: FieldValue): void {
    if (isFeatureField(field.name)) {
      const [, featureName] = field.name.split('.') as [string, keyof Features];
      formData.features = {
        ...formData.features,
        [featureName]: value,
      };
    } else {
      formData = {
        ...formData,
        [field.name]: value,
      };
    }
  }

  function isLastStep(): boolean {
    return currentStep === TOTAL_STEPS - 1;
  }

  // Update radio button template
  function isSelected(field: FormField, option: string): boolean {
    if (field.name.startsWith('features.')) {
      const [, featureName] = field.name.split('.') as [string, keyof Features];
      return formData.features[featureName]?.toString() === option;
    }
    return formData[field.name] === option;
  }
</script>

{#if showThankYou}
  <div class="thank-you-card">
    <h2>Thank You!</h2>
    <p>Your feedback is valuable in helping us build LocalHaven CMS.</p>
    {#if formData.betaInterest}
      <p>We'll be in touch about the beta program soon!</p>
    {/if}
  </div>
{:else}
  <form on:submit={handleSubmit}>
    {#each currentFields as field}
      <div class="form-field">
        <label for={field.id}>
          {field.label}
          {#if field.required}
            <span class="required">*</span>
          {/if}
        </label>

        {#if !field.dependsOn || formData[field.dependsOn] === field.dependsOnValue}
          {#if field.type === 'select'}
            <select
              id={field.id}
              name={field.name}
              bind:value={formData[field.name]}
              on:change={(e) => handleChange(e, field)}
              disabled={isSubmitting}
            >
              <option value="">Select an option</option>
              {#each field.options || [] as option}
                <option value={option}>{option}</option>
              {/each}ß
            </select>
          {:else if field.type === 'textarea'}
            <textarea
              id={field.id}
              name={field.name}
              bind:value={formData[field.name]}
              on:input={(e) => handleChange(e, field)}
              disabled={isSubmitting}
            ></textarea>
          {:else if field.type === 'checkbox'}
            <input
              type="checkbox"
              id={field.id}
              name={field.name}
              checked={(formData[field.name] as boolean) || false}
              on:change={(e) => handleChange(e, field)}
              disabled={isSubmitting}
            />
          {:else if field.type === 'radio'}
            <fieldset role="radiogroup">
              {#each field.options || [] as option}
                <div class="radio-option">
                  <input
                    type="radio"
                    id={field.id + option}
                    name={field.name}
                    value={option}
                    checked={isSelected(field, option)}
                    on:change={(e) => handleChange(e, field)}
                    disabled={isSubmitting}
                  />
                  <label for={field.id + option}>{option}</label>
                </div>
              {/each}
            </fieldset>
          {:else}
            <input
              type={field.type}
              id={field.id}
              name={field.name}
              bind:value={formData[field.name]}
              on:input={(e) => handleChange(e, field)}
              disabled={isSubmitting}
            />
          {/if}
        {/if}

        {#if errors[field.name]}
          <span class="error-message">{errors[field.name]}</span>
        {/if}
      </div>
    {/each}

    <div class="button-group">
      {#if currentStep > 0}
        <button type="button" on:click={() => currentStep--} disabled={isSubmitting}>
          Previous
        </button>
      {/if}

      <button type="submit" disabled={isSubmitting}>
        {#if isSubmitting}
          Submitting...
        {:else}
          {isLastStep() ? 'Submit Survey' : 'Next Step'}
        {/if}
      </button>
    </div>

    {#if submissionStatus === 'error'}
      <div class="error-message">Failed to submit survey. Please try again.</div>
    {/if}
  </form>
{/if}

<style>
  form,
  textarea,
  input {
    font-family: var(--font-family-base);
    font-size: var(--font-size-base);
  }

  .input-group {
    margin-bottom: 1.5rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
  }

  input[type='radio'],
  input[type='checkbox'] {
    margin-right: 0.5rem;
  }

  fieldset {
    border: none;
    padding: 0;
    margin: 0;
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  fieldset label {
    margin-left: 0.25rem;
    display: inline-flex;
    align-items: center;
  }

  input[type='radio'] {
    margin: 0;
  }

  .radio-option {
    align-items: center;
  }

  .error-message {
    color: #dc2626;
    font-size: 0.875rem;
    margin-top: 0.25rem;
  }

  .form-field {
    margin-bottom: 1.5rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }

  .required {
    color: #dc2626;
    margin-left: 0.25rem;
  }

  input[type='text'],
  input[type='email'],
  select,
  textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #d1d5db;
    border-radius: 0.375rem;
    background-color: white;
  }

  textarea {
    min-height: 100px;
    resize: vertical;
  }

  .button-group {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
  }

  button {
    padding: 0.5rem 1rem;
    background-color: #059669;
    color: white;
    border: none;
    border-radius: 0.375rem;
    font-weight: 500;
    cursor: pointer;
  }

  button:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .thank-you-card {
    background: white;
    padding: 2rem;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    text-align: center;
    max-width: 400px;
    margin: 2rem auto;
  }

  .thank-you-card h2 {
    color: #059669;
    margin-bottom: 1rem;
  }

  .thank-you-card p {
    color: #4b5563;
    line-height: 1.5;
    margin-bottom: 1rem;
  }
</style>
