<script lang="ts">
  import { onMount } from 'svelte';
  import type {
    SurveyResponse,
    FeatureAverages,
    RoleDistribution,
    UsageDistribution,
  } from '../types/Survey';
  import { config } from '../config';
  import { Bar } from 'svelte-chartjs';
  import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
  } from 'chart.js';

  ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

  let surveys: SurveyResponse[] = [];
  let isLoading = true;
  let error = '';

  let featureAverages: FeatureAverages = {};
  let roleDistribution: RoleDistribution = {};
  let usageDistribution: UsageDistribution = {};
  let betaInterestCount = 0;
  let totalResponses = 0;

  const featureLabels = {
    offline: 'Offline Capabilities',
    collaboration: 'Real-time Collaboration',
    assetManagement: 'Asset Management',
    pdfHandling: 'PDF Handling',
    versionControl: 'Version Control',
    workflows: 'Approval Workflows',
  };

  onMount(async () => {
    try {
      const response = await fetch(`${config.apiUrl}/surveys`);
      if (!response.ok) throw new Error('Failed to fetch survey data');

      surveys = await response.json();
      calculateMetrics();
    } catch (e) {
      error = e.message;
    } finally {
      isLoading = false;
    }
  });

  function calculateMetrics() {
    totalResponses = surveys.length;
    betaInterestCount = surveys.filter((s) => s.betaInterest).length;

    // Calculate feature averages
    const features = Object.keys(featureLabels);
    features.forEach((feature) => {
      const sum = surveys.reduce((acc, survey) => acc + survey.features[feature], 0);
      featureAverages[feature] = Number((sum / totalResponses).toFixed(2));
    });

    // Calculate role distribution
    surveys.forEach((survey) => {
      const role = survey.role === 'other' ? survey.otherRole || 'Other' : survey.role;
      roleDistribution[role] = (roleDistribution[role] || 0) + 1;
    });

    // Calculate usage distribution
    surveys.forEach((survey) => {
      const usage = survey.cmsUsage === 'other' ? survey.otherCmsUsage || 'Other' : survey.cmsUsage;
      usageDistribution[usage] = (usageDistribution[usage] || 0) + 1;
    });
  }

  $: featureChartData = {
    labels: Object.values(featureLabels),
    datasets: [
      {
        label: 'Average Rating',
        data: Object.keys(featureLabels).map((key) => featureAverages[key]),
        backgroundColor: '#047857',
      },
    ],
  };
</script>

<div class="dashboard">
  <h1>Survey Analytics Dashboard</h1>

  {#if isLoading}
    <div class="loading">Loading survey data...</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else}
    <div class="metrics-grid">
      <div class="metric-card">
        <h3>Total Responses</h3>
        <p class="metric-value">{totalResponses}</p>
      </div>
      <div class="metric-card">
        <h3>Beta Interest</h3>
        <p class="metric-value">{betaInterestCount}</p>
        <p class="metric-subtitle">
          ({Math.round((betaInterestCount / totalResponses) * 100)}% of total)
        </p>
      </div>
    </div>

    <div class="chart-container">
      <h2>Feature Importance Ratings</h2>
      <Bar data={featureChartData} options={{ responsive: true, maintainAspectRatio: false }} />
    </div>

    <div class="tables-grid">
      <div class="table-container">
        <h2>Role Distribution</h2>
        <table>
          <thead>
            <tr>
              <th>Role</th>
              <th>Count</th>
            </tr>
          </thead>
          <tbody>
            {#each Object.entries(roleDistribution) as [role, count]}
              <tr>
                <td>{role}</td>
                <td>{count}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <div class="table-container">
        <h2>CMS Usage</h2>
        <table>
          <thead>
            <tr>
              <th>Frequency</th>
              <th>Count</th>
            </tr>
          </thead>
          <tbody>
            {#each Object.entries(usageDistribution) as [usage, count]}
              <tr>
                <td>{usage}</td>
                <td>{count}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>

<style>
  .dashboard {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  h1 {
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: 2rem;
    color: #111827;
  }

  h2 {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 1rem;
    color: #374151;
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
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .metric-value {
    font-size: 2rem;
    font-weight: 700;
    color: #047857;
  }

  .metric-subtitle {
    font-size: 0.875rem;
    color: #6b7280;
  }

  .chart-container {
    background: white;
    padding: 1.5rem;
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    margin-bottom: 2rem;
    height: 400px;
  }

  .tables-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1rem;
  }

  .table-container {
    background: white;
    padding: 1.5rem;
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid #e5e7eb;
  }

  th {
    font-weight: 600;
    color: #374151;
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: #6b7280;
  }

  .error {
    text-align: center;
    padding: 2rem;
    color: #dc2626;
    background: #fee2e2;
    border-radius: 0.5rem;
  }
</style>
