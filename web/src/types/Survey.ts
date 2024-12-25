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
