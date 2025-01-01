export interface SurveyResponse {
  id: string;
  createdAt: string;
  role: string;
  otherRole?: string;
  cmsUsage: string;
  otherCmsUsage?: string;
  features: {
    offline: number;
    collaboration: number;
    assetManagement: number;
    pdfHandling: number;
    versionControl: number;
    workflows: number;
  };
  betaInterest: boolean;
  email?: string;
  // New fields
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
  feedbackSuggestions: string;
  excitementFactors: string;
}

export interface FeatureAverages {
  [key: string]: number;
}

export interface RoleDistribution {
  [key: string]: number;
}

export interface UsageDistribution {
  [key: string]: number;
}
