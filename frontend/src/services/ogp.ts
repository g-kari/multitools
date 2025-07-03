import type { OGPRequest, OGPResponse } from '@/types/ogp';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export class OGPService {
  private static instance: OGPService;

  private constructor() {}

  public static getInstance(): OGPService {
    if (!OGPService.instance) {
      OGPService.instance = new OGPService();
    }
    return OGPService.instance;
  }

  public async verifyOGP(request: OGPRequest): Promise<OGPResponse> {
    const response = await fetch(`${API_BASE_URL}/api/v1/ogp/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }

    return response.json();
  }

  public async healthCheck(): Promise<{ status: string; timestamp: string }> {
    const response = await fetch(`${API_BASE_URL}/health`);
    
    if (!response.ok) {
      throw new Error(`Health check failed: ${response.status}`);
    }

    return response.json();
  }
}