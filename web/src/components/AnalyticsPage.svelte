<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '../stores/auth';
  import type { SurveyResponse } from '../types/Survey';

  let surveyResults: SurveyResponse[] = [];
  let loading = true;
  let error = '';

  onMount(async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('http://api.localhavencms.com/results', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          auth.logout();
          window.location.href = '/login';
          return;
        }
        throw new Error('Failed to fetch survey results');
      }

      surveyResults = await response.json();
    } catch (e) {
      error = 'Failed to load survey results';
      console.error(e);
    } finally {
      loading = false;
    }
  });
</script>

<div class="analytics-container">
  <h1>Survey Results</h1>

  {#if loading}
    <p>Loading...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else}
    <div class="results-grid">
      {#each surveyResults as result}
        <div class="result-card">
          <h3>Response #{result.id}</h3>
          <p><strong>Role:</strong> {result.role}</p>
          <p><strong>CMS Usage:</strong> {result.cmsUsage}</p>
          <h4>Feature Ratings:</h4>
          <ul>
            <li>Offline: {result.features.offline}/5</li>
            <li>Collaboration: {result.features.collaboration}/5</li>
            <li>Asset Management: {result.features.assetManagement}/5</li>
            <li>PDF Handling: {result.features.pdfHandling}/5</li>
            <li>Version Control: {result.features.versionControl}/5</li>
            <li>Workflows: {result.features.workflows}/5</li>
          </ul>
          {#if result.betaInterest}
            <p class="beta-interest">Interested in Beta Program</p>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .analytics-container {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .results-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
    margin-top: 2rem;
  }

  .result-card {
    background: white;
    padding: 1.5rem;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .beta-interest {
    color: #047857;
    font-weight: 500;
  }

  .error {
    color: #dc2626;
    padding: 1rem;
    background-color: #fee2e2;
    border-radius: 0.5rem;
  }

  h1 {
    color: #111827;
    margin-bottom: 2rem;
  }

  ul {
    list-style: none;
    padding: 0;
  }

  li {
    margin-bottom: 0.5rem;
  }
</style>
