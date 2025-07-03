import { renderHook, act } from '@testing-library/react';
import { useOGP } from './useOGP';
import { OGPService } from '@/services/ogp';
import type { OGPResponse } from '@/types/ogp';

// Mock the OGPService
jest.mock('@/services/ogp');

const mockOGPService = OGPService as jest.MockedClass<typeof OGPService>;

describe('useOGP', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockOGPService.getInstance.mockReturnValue({
      verifyOGP: jest.fn(),
      healthCheck: jest.fn(),
    } as any);
  });

  it('should initialize with default values', () => {
    const { result } = renderHook(() => useOGP());

    expect(result.current.data).toBeNull();
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('should handle successful OGP verification', async () => {
    const mockResponse: OGPResponse = {
      url: 'https://example.com',
      ogp_data: {
        title: 'Test Title',
        description: 'Test Description',
        image: 'https://example.com/image.jpg',
        url: 'https://example.com',
        type: 'website',
        site_name: 'Test Site',
        image_width: '1200',
        image_height: '630',
        image_alt: 'Test Image',
      },
      validation: {
        is_valid: true,
        warnings: [],
        errors: [],
        checks: {
          has_title: true,
          has_description: true,
          has_image: true,
          image_valid: true,
          url_valid: true,
        },
      },
      previews: {
        twitter: {
          platform: 'twitter',
          title: 'Test Title',
          description: 'Test Description',
          image: 'https://example.com/image.jpg',
          is_valid: true,
          warnings: [],
          title_length: 10,
          desc_length: 16,
          max_title_len: 70,
          max_desc_len: 200,
        },
        facebook: {
          platform: 'facebook',
          title: 'Test Title',
          description: 'Test Description',
          image: 'https://example.com/image.jpg',
          is_valid: true,
          warnings: [],
          title_length: 10,
          desc_length: 16,
          max_title_len: 100,
          max_desc_len: 300,
        },
        discord: {
          platform: 'discord',
          title: 'Test Title',
          description: 'Test Description',
          image: 'https://example.com/image.jpg',
          is_valid: true,
          warnings: [],
          title_length: 10,
          desc_length: 16,
          max_title_len: 256,
          max_desc_len: 2048,
        },
      },
      timestamp: '2024-01-01T00:00:00Z',
    };

    const mockVerifyOGP = jest.fn().mockResolvedValue(mockResponse);
    mockOGPService.getInstance.mockReturnValue({
      verifyOGP: mockVerifyOGP,
      healthCheck: jest.fn(),
    } as any);

    const { result } = renderHook(() => useOGP());

    await act(async () => {
      await result.current.verifyOGP({ url: 'https://example.com' });
    });

    expect(result.current.data).toEqual(mockResponse);
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(mockVerifyOGP).toHaveBeenCalledWith({ url: 'https://example.com' });
  });

  it('should handle OGP verification error', async () => {
    const mockError = new Error('Network error');
    const mockVerifyOGP = jest.fn().mockRejectedValue(mockError);
    mockOGPService.getInstance.mockReturnValue({
      verifyOGP: mockVerifyOGP,
      healthCheck: jest.fn(),
    } as any);

    const { result } = renderHook(() => useOGP());

    await act(async () => {
      await result.current.verifyOGP({ url: 'https://example.com' });
    });

    expect(result.current.data).toBeNull();
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBe('Network error');
  });

  it('should set loading state during verification', async () => {
    const mockVerifyOGP = jest.fn().mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));
    mockOGPService.getInstance.mockReturnValue({
      verifyOGP: mockVerifyOGP,
      healthCheck: jest.fn(),
    } as any);

    const { result } = renderHook(() => useOGP());

    act(() => {
      result.current.verifyOGP({ url: 'https://example.com' });
    });

    expect(result.current.loading).toBe(true);
  });

  it('should reset state', () => {
    const { result } = renderHook(() => useOGP());

    act(() => {
      result.current.reset();
    });

    expect(result.current.data).toBeNull();
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });
});