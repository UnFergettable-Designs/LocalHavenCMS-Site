<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '../stores/auth';
  import { config } from '../config';
  import type { SurveyResponse, MetricsData, Features } from '../types/Survey';
  import AnalyticsDashboard from './AnalyticsDashboard.svelte';

  let surveyResults: SurveyResponse[] = [];
  let loading = true;
  let error = '';
  let deletingIds: Set<string> = new Set();

  // Add new metrics
  let totalResponses = 0;
  let averageFeatureScores: Record<keyof Features, number> = {
    offline: 0,
    collaboration: 0,
    assetManagement: 0,
    pdfHandling: 0,
    versionControl: 0,
    workflows: 0,
  };
  let roleDistribution: Record<string, number> = {};
  let cmsUsageDistribution: Record<string, number> = {};
  let usageFrequencyDistribution: Record<string, number> = {};
  let betaInterestCount = 0;
  let metrics: MetricsData | null = null;

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
        throw new Error('Failed to fetch results');
      }

      surveyResults = await response.json();
      calculateMetrics();
    } catch (e) {
      console.error('Error:', e);
      error = 'Failed to fetch survey results';
    } finally {
      loading = false;
    }
  });

  function calculateMetrics() {
    totalResponses = surveyResults.length;
    betaInterestCount = surveyResults.filter((s) => s.betaInterest).length;

    // Reset distributions
    roleDistribution = {};
    cmsUsageDistribution = {};
    usageFrequencyDistribution = {};
    averageFeatureScores = {
      offline: 0,
      collaboration: 0,
      assetManagement: 0,
      pdfHandling: 0,
      versionControl: 0,
      workflows: 0,
    };

    // Calculate distributions and averages
    surveyResults.forEach((survey) => {
      // Role distribution
      const role = survey.role === 'Other' ? survey.otherRole || 'Other' : survey.role;
      roleDistribution[role] = (roleDistribution[role] || 0) + 1;

      // CMS usage distribution
      const cmsUsage =
        survey.cmsUsage === 'Other' ? survey.otherCmsUsage || 'Other' : survey.cmsUsage;
      cmsUsageDistribution[cmsUsage] = (cmsUsageDistribution[cmsUsage] || 0) + 1;

      // Usage frequency distribution
      usageFrequencyDistribution[survey.usageFrequency] =
        (usageFrequencyDistribution[survey.usageFrequency] || 0) + 1;

      // Feature scores
      Object.entries(survey.features).forEach(([feature, score]) => {
        averageFeatureScores[feature] += score;
      });
    });

    // Calculate averages for feature scores
    Object.keys(averageFeatureScores).forEach((feature) => {
      averageFeatureScores[feature] = Number(
        (averageFeatureScores[feature] / totalResponses).toFixed(2)
      );
    });
  }

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
      calculateMetrics(); // Recalculate metrics after deletion
    } catch (e) {
      error = 'Failed to delete result';
      console.error('Error:', e);
    } finally {
      deletingIds.delete(id);
      deletingIds = deletingIds; // trigger reactivity
    }
  }

  async function fetchResults() {
    try {
      loading = true;
      error = '';

      const response = await fetch(`${config.apiUrl}/results`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          auth.logout();
          window.location.href = '/login';
          return;
        }
        throw new Error(`Server responded with status: ${response.status}`);
      }

      surveyResults = await response.json();
    } catch (err: any) {
      error = err.message || 'Failed to fetch results';
      console.error('Error fetching results:', err);
    } finally {
      loading = false;
    }
  }
</script>

<div class="analytics-container">
  <h1>Survey Results</h1>

  {#if error}
    <div class="error">{error}</div>
  {/if}

  {#if loading}
    <div class="loading">Loading...</div>
  {:else}
    <div class="metrics-grid">
      <div class="metric-card">
        <h3>Total Responses</h3>
        <p class="metric-value">{totalResponses}</p>
      </div>
      <div class="metric-card">
        <h3>Beta Interest</h3>
        <p class="metric-value beta-interest">{betaInterestCount} / {totalResponses}</p>
      </div>
    </div>

    <div class="charts-section">
      <div class="chart-container">
        <h3>Feature Importance</h3>
        <div class="feature-scores">
          {#each Object.entries(averageFeatureScores) as [feature, score]}
            <div class="feature-score">
              <span class="feature-name">{feature}</span>
              <div class="score-bar" style="width: {score * 20}%">{score}</div>
            </div>
          {/each}
        </div>
      </div>

      <div class="chart-container">
        <h3>Role Distribution</h3>
        <div class="distribution-list">
          {#each Object.entries(roleDistribution) as [role, count]}
            <div class="distribution-item">
              <span>{role}</span>
              <span class="count">{count}</span>
            </div>
          {/each}
        </div>
      </div>

      <div class="chart-container">
        <h3>CMS Usage</h3>
        <div class="distribution-list">
          {#each Object.entries(cmsUsageDistribution) as [cms, count]}
            <div class="distribution-item">
              <span>{cms}</span>
              <span class="count">{count}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>

    <div class="results-grid">
      {#each surveyResults as result}
        <div class="result-card">
          <button
            class="delete-button"
            on:click={() => deleteResult(result.id)}
            disabled={deletingIds.has(result.id)}
          >
            {#if deletingIds.has(result.id)}
              Deleting...
            {:else}
              Delete
            {/if}
          </button>
          <h3>Response {result.id}</h3>
          <ul>
            <li><strong>Role:</strong> {result.role}</li>
            <li><strong>CMS Usage:</strong> {result.cmsUsage}</li>
            <li><strong>Usage Frequency:</strong> {result.usageFrequency}</li>
            <li><strong>Primary Purpose:</strong> {result.primaryPurpose}</li>
            <li><strong>Beta Interest:</strong> {result.betaInterest ? 'Yes' : 'No'}</li>
            {#if result.email}
              <li><strong>Email:</strong> {result.email}</li>
            {/if}
            <li>
              <strong>Feature Ratings:</strong>
              <ul class="feature-list">
                {#each Object.entries(result.features) as [feature, rating]}
                  <li>{feature}: {rating}/5</li>
                {/each}
              </ul>
            </li>
          </ul>
        </div>
      {/each}
    </div>
  {/if}
  <AnalyticsDashboard {surveyResults} />
</div>

<style>
  .analytics-container {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .metrics-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .metric-card {
    background: white;
    padding: 1.5rem;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    text-align: center;
  }

  .metric-value {
    font-size: 2rem;
    font-weight: 600;
    color: #059669;
  }

  .charts-section {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    margin-bottom: 2rem;
  }

  .chart-container {
    background: white;
    padding: 1.5rem;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .feature-scores {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .feature-score {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .feature-name {
    min-width: 120px;
  }

  .score-bar {
    background: #059669;
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    min-width: 2rem;
    text-align: right;
  }

  .distribution-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .distribution-item {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem;
    background: #f3f4f6;
    border-radius: 0.25rem;
  }

  .count {
    font-weight: 600;
    color: #059669;
  }

  .results-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
  }

  .result-card {
    position: relative;
    background: white;
    padding: 1.5rem;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .delete-button {
    position: absolute;
    top: 1rem;
    right: 1rem;
    padding: 0.5rem 1rem;
    background: #dc2626;
    color: white;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }

  .delete-button:disabled {
    background: #f87171;
    cursor: not-allowed;
  }

  .feature-list {
    margin-top: 0.5rem;
    padding-left: 1rem;
  }

  .error {
    color: #dc2626;
    padding: 1rem;
    background-color: #fee2e2;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: #6b7280;
  }

  h1 {
    color: #111827;
    margin-bottom: 2rem;
  }

  h3 {
    color: #374151;
    margin-bottom: 1rem;
  }

  ul {
    list-style: none;
    padding: 0;
  }

  li {
    margin-bottom: 0.5rem;
  }

  strong {
    color: #374151;
  }
</style>
