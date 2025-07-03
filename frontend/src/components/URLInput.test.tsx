import { render, screen, fireEvent } from '@testing-library/react';
import { URLInput } from './URLInput';

describe('URLInput', () => {
  const mockOnSubmit = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders correctly', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={false} />);
    
    expect(screen.getByLabelText('URL to verify')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('https://example.com')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Verify OGP' })).toBeInTheDocument();
  });

  it('calls onSubmit with URL when form is submitted', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={false} />);
    
    const input = screen.getByPlaceholderText('https://example.com');
    const button = screen.getByRole('button', { name: 'Verify OGP' });
    
    fireEvent.change(input, { target: { value: 'https://example.com' } });
    fireEvent.click(button);
    
    expect(mockOnSubmit).toHaveBeenCalledWith('https://example.com');
  });

  it('trims whitespace from URL', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={false} />);
    
    const input = screen.getByPlaceholderText('https://example.com');
    const button = screen.getByRole('button', { name: 'Verify OGP' });
    
    fireEvent.change(input, { target: { value: '  https://example.com  ' } });
    fireEvent.click(button);
    
    expect(mockOnSubmit).toHaveBeenCalledWith('https://example.com');
  });

  it('does not submit empty URL', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={false} />);
    
    const button = screen.getByRole('button', { name: 'Verify OGP' });
    
    fireEvent.click(button);
    
    expect(mockOnSubmit).not.toHaveBeenCalled();
  });

  it('does not submit whitespace-only URL', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={false} />);
    
    const input = screen.getByPlaceholderText('https://example.com');
    const button = screen.getByRole('button', { name: 'Verify OGP' });
    
    fireEvent.change(input, { target: { value: '   ' } });
    fireEvent.click(button);
    
    expect(mockOnSubmit).not.toHaveBeenCalled();
  });

  it('shows loading state', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={true} />);
    
    const input = screen.getByPlaceholderText('https://example.com');
    const button = screen.getByRole('button', { name: 'Verifying...' });
    
    expect(input).toBeDisabled();
    expect(button).toBeDisabled();
  });

  it('disables button when URL is empty', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={false} />);
    
    const button = screen.getByRole('button', { name: 'Verify OGP' });
    
    expect(button).toBeDisabled();
  });

  it('enables button when URL is entered', () => {
    render(<URLInput onSubmit={mockOnSubmit} loading={false} />);
    
    const input = screen.getByPlaceholderText('https://example.com');
    const button = screen.getByRole('button', { name: 'Verify OGP' });
    
    fireEvent.change(input, { target: { value: 'https://example.com' } });
    
    expect(button).not.toBeDisabled();
  });
});