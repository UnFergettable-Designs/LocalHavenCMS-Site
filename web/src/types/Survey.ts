import type { Component } from 'svelte';

export interface Features {
  offline?: number;
  collaboration?: number;
  assetManagement?: number;
  pdfHandling?: number;
  versionControl?: number;
  workflows?: number;
}

export interface SurveyResponse {
  id: string;
  role: string;
  otherRole?: string;
  cmsUsage: string;
  otherCmsUsage?: string;
  features: Features;
  betaInterest: boolean;
  email?: string;
  createdAt: string;
  biggestFrustrations: string;
  specificProblems: string;
  usageFrequency: string;
  primaryPurpose: string;
  platforms: string;
  cmsPreference: string;
  wishedFeatures: string;
  workflowImportance: string;
  teamSize: string;
  collaborationFrequency: string;
  pricingSensitivity: string;
  pricingModel: string;
  integrations: string;
  integrationImportance: string;
  contentTypes: string;
  customFormats: string;
  feedbackSuggestions?: string;
  excitementFactors?: string;
  collaborationChallenges?: string;
  offlineWorkFrequency?: string;
  offlineWorkarounds?: string;
  currentChangeConflictHandling?: string;
  versionControlChallenges?: string;
}

export interface FormField {
  id: string;
  name: string;
  label: string;
  type: 'text' | 'email' | 'select' | 'textarea' | 'radio' | 'checkbox';
  required?: boolean;
  options?: string[];
  dependsOn?: keyof SurveyResponse;
  dependsOnValue?: string | boolean;
  placeholder?: string;
  validation?: {
    pattern?: string;
    min?: number;
    max?: number;
    minLength?: number;
    maxLength?: number;
  };
}

export interface ChartData {
  labels: string[];
  datasets: Array<{
    label?: string;
    data: number[];
    backgroundColor?: string | string[];
  }>;
}

export interface MetricsData {
  totalResponses: number;
  betaInterestCount: number;
  roleDistribution: Record<string, number>;
  cmsUsageDistribution: Record<string, number>;
  featureScores: Record<keyof Features, number>;
}

export interface ChartConfiguration {
  labels: string[];
  datasets: Array<{
    label?: string;
    data: number[];
    backgroundColor: string | string[];
  }>;
}

export interface DashboardMetrics {
  totalResponses: number;
  betaInterestCount: number;
  featureScores: Record<keyof Features, number>;
  distributions: {
    roles: Record<string, number>;
    cmsUsage: Record<string, number>;
    teamSizes: Record<string, number>;
    pricing: Record<string, number>;
  };
}

export interface Feature {
  icon: Component;
  title: string;
  description: string;
}

export type FieldValue = string | number | boolean | undefined;

export interface SurveyFormData {
  [key: string]: FieldValue | Features | Record<string, any>;
  features: {
    offline: number;
    collaboration: number;
    assetManagement: number;
    pdfHandling: number;
    versionControl: number;
    workflows: number;
  };
  role: string;
  otherRole?: string;
  cmsUsage: string;
  otherCmsUsage?: string;
  betaInterest: boolean;
  email?: string;
  biggestFrustrations: string;
  specificProblems: string;
  usageFrequency: string;
  primaryPurpose: string;
  platforms: string;
  cmsPreference: string;
  wishedFeatures: string;
  workflowImportance: string;
  teamSize: string;
  collaborationFrequency: string;
  pricingSensitivity: string;
  pricingModel: string;
  integrations: string;
  integrationImportance: string;
  contentTypes: string;
  customFormats: string;
  feedbackSuggestions?: string;
  excitementFactors?: string;
  collaborationChallenges?: string;
  offlineWorkFrequency?: string;
  offlineWorkarounds?: string;
  currentChangeConflictHandling?: string;
  versionControlChallenges?: string;
}
