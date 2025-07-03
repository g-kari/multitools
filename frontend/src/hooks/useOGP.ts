import { useState, useCallback } from 'react';
import { OGPService } from '@/services/ogp';
import type { OGPRequest, OGPResponse } from '@/types/ogp';

export interface UseOGPReturn {
  data: OGPResponse | null;
  loading: boolean;
  error: string | null;
  verifyOGP: (request: OGPRequest) => Promise<void>;
  reset: () => void;
}

export const useOGP = (): UseOGPReturn => {
  const [data, setData] = useState<OGPResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const ogpService = OGPService.getInstance();

  const verifyOGP = useCallback(async (request: OGPRequest) => {
    setLoading(true);
    setError(null);

    try {
      const response = await ogpService.verifyOGP(request);
      setData(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unknown error occurred');
    } finally {
      setLoading(false);
    }
  }, [ogpService]);

  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);

  return {
    data,
    loading,
    error,
    verifyOGP,
    reset,
  };
};