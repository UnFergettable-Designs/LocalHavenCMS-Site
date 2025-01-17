<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '../stores/auth';
  import { config } from '../config';
  import type { SurveyResponse } from '../types/Survey';

  let surveyResults: SurveyResponse[] = [];
  let loading = true;
  let error = '';
  let deletingIds: Set<string> = new Set();

  onMount(async () => {
    const token = localStorage.getItem('token');
    console.log('Token from localStorage:', token ? 'exists' : 'missing');

    if (!token) {
      console.log('No token found, redirecting to login');
      window.location.href = '/login';
      return;
    }

    try {
      // First verify the token
      console.log('Verifying token...');
      const verifyResponse = await fetch(`${config.apiUrl}/verify`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!verifyResponse.ok) {
        console.log('Token verification failed');
        auth.logout();
        window.location.href = '/login';
        return;
      }

      // Then fetch survey results
      console.log('Token verified, fetching results...');
      const response = await fetch(`${config.apiUrl}/results`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          console.log('Unauthorized when fetching results');
          auth.logout();
          window.location.href = '/login';
          return;
        }
        throw new Error('Failed to fetch survey results');
      }

      surveyResults = await response.json();
      console.log('Results fetched successfully');
    } catch (e) {
      error = 'Failed to load survey results';
      console.error('Error:', e);
    } finally {
      loading = false;
    }
  });

  async function deleteResult(id: string) {
    try {
      deletingIds.add(id);
      deletingIds = deletingIds; // trigger reactivity

      const token = localStorage.getItem('token');
      const response = await fetch(`${config.apiUrl}/results/${id}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (!response.ok) {
        throw new Error('Failed to delete result');
      }

      surveyResults = surveyResults.filter((result) => result.id !== id);
    } catch (e) {
      error = 'Failed to delete result';
      console.error('Error:', e);
    } finally {
      deletingIds.delete(id);
      deletingIds = deletingIds; // trigger reactivity
    }
  }
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
          {#if result.email}
            <p><strong>Email:</strong> {result.email}</p>
          {/if}
          <button on:click={() => deleteResult(result.id)} disabled={deletingIds.has(result.id)}>
            {#if deletingIds.has(result.id)}
              Deleting...
            {:else}
              Delete
            {/if}
          </button>
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

  button:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
</style>
