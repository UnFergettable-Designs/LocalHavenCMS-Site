<script lang="ts">
  import { onMount } from 'svelte';
  import { config } from '../config';
  import Chart from 'chart.js/auto';
  import type { ChartData, ChartOptions, ChartType, Chart as ChartInstance } from 'chart.js';
  import type { SurveyResponse, Features, DashboardMetrics } from '../types/Survey';

  interface TypedChartData extends ChartData {
    labels: string[];
    datasets: Array<{
      label?: string;
      data: number[];
      backgroundColor: string | string[];
    }>;
  }

  interface DashboardCharts {
    feature?: ChartInstance<'bar', number[], string>;
    role?: ChartInstance<'doughnut', number[], string>;
    cms?: ChartInstance<'doughnut', number[], string>;
    teamSize?: ChartInstance<'doughnut', number[], string>;
    pricing?: ChartInstance<'doughnut', number[], string>;
  }

  export let surveyResults: SurveyResponse[] = [];

  let metrics: DashboardMetrics = {
    totalResponses: 0,
    betaInterestCount: 0,
    featureScores: {
      offline: 0,
      collaboration: 0,
      assetManagement: 0,
      pdfHandling: 0,
      versionControl: 0,
      workflows: 0,
    },
    distributions: {
      roles: {},
      cmsUsage: {},
      teamSizes: {},
      pricing: {},
    },
  };

  let charts: DashboardCharts = {};

  const chartOptions: ChartOptions<'bar'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom',
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 5,
      },
    },
  };

  const pieOptions: ChartOptions<'doughnut'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom',
      },
    },
  };

  let featureScoresChart: TypedChartData = {
    labels: [],
    datasets: [
      {
        label: 'Feature Scores',
        data: [],
        backgroundColor: '#059669',
      },
    ],
  };

  let roleDistributionChart: TypedChartData = {
    labels: [],
    datasets: [
      {
        data: [],
        backgroundColor: ['#059669', '#0891b2', '#6366f1', '#8b5cf6', '#ec4899'],
      },
    ],
  };

  onMount(() => {
    initializeCharts();
    return () => {
      Object.values(charts).forEach((chart) => {
        if (chart) chart.destroy();
      });
    };
  });

  function initializeCharts(): void {
    const getContext = (id: string): CanvasRenderingContext2D | null => {
      const canvas = document.getElementById(id) as HTMLCanvasElement;
      return canvas?.getContext('2d');
    };

    const featureCtx = getContext('featureChart');
    const roleCtx = getContext('roleChart');
    const cmsCtx = getContext('cmsChart');
    const teamSizeCtx = getContext('teamSizeChart');
    const pricingCtx = getContext('pricingChart');

    if (featureCtx) {
      charts.feature = new Chart(featureCtx, {
        type: 'bar',
        data: featureScoresChart,
        options: chartOptions,
      });
    }

    if (roleCtx) {
      charts.role = new Chart(roleCtx, {
        type: 'doughnut',
        data: roleDistributionChart,
        options: pieOptions,
      });
    }

    // ... similar for other charts
  }

  function calculateMetrics(): void {
    metrics.totalResponses = surveyResults.length;
    metrics.betaInterestCount = surveyResults.filter((s) => s.betaInterest).length;

    for (const feature of Object.keys(metrics.featureScores) as Array<keyof Features>) {
      metrics.featureScores[feature] =
        surveyResults.reduce((sum, response) => sum + (response.features[feature] || 0), 0) /
        metrics.totalResponses;
    }

    metrics.distributions = {
      roles: calculateFrequencyDistribution(surveyResults.map((r) => r.role)),
      cmsUsage: calculateFrequencyDistribution(surveyResults.map((r) => r.cmsUsage)),
      teamSizes: calculateFrequencyDistribution(surveyResults.map((r) => r.teamSize)),
      pricing: calculateFrequencyDistribution(surveyResults.map((r) => r.pricingModel)),
    };
  }

  function calculateFrequencyDistribution(items: (string | undefined)[]): Record<string, number> {
    return items.reduce(
      (acc, item) => {
        if (item) {
          acc[item] = (acc[item] || 0) + 1;
        }
        return acc;
      },
      {} as Record<string, number>
    );
  }

  $: if (surveyResults.length > 0) {
    calculateMetrics();
    updateCharts();
  }
</script>

<div class="dashboard">
  {#if metrics}
    <div class="metrics-summary">
      <div class="metric-card">
        <h4>Total Responses</h4>
        <p class="metric-value">{metrics.totalResponses}</p>
      </div>
      <div class="metric-card">
        <h4>Beta Interest</h4>
        <p class="metric-value">{metrics.betaInterestCount} / {metrics.totalResponses}</p>
      </div>
    </div>
  {/if}

  <div class="chart-grid">
    <div class="chart-container">
      <h3>Feature Importance Ratings</h3>
      <div class="chart">
        <canvas id="featureChart"></canvas>
      </div>
    </div>

    <div class="chart-container">
      <h3>Role Distribution</h3>
      <div class="chart">
        <canvas id="roleChart"></canvas>
      </div>
    </div>

    <div class="chart-container">
      <h3>CMS Usage Distribution</h3>
      <div class="chart">
        <canvas id="cmsChart"></canvas>
      </div>
    </div>

    <div class="chart-container">
      <h3>Team Size Distribution</h3>
      <div class="chart">
        <canvas id="teamSizeChart"></canvas>
      </div>
    </div>

    <div class="chart-container">
      <h3>Pricing Preferences</h3>
      <div class="chart">
        <canvas id="pricingChart"></canvas>
      </div>
    </div>
  </div>
</div>

<style>
  .dashboard {
    padding: 1rem;
  }

  .chart-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
  }

  .chart-container {
    background: white;
    padding: 1.5rem;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  h3 {
    color: #374151;
    margin-bottom: 1rem;
    text-align: center;
  }

  .chart {
    height: 300px;
    position: relative;
  }

  .metrics-summary {
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
</style>
