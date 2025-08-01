import { render, screen } from '@testing-library/react'
import { Badge } from '@/components/ui/badge'

describe('Badge Component', () => {
  it('renders correctly with default props', () => {
    render(<Badge>Default Badge</Badge>)
    
    const badge = screen.getByText('Default Badge')
    expect(badge).toBeInTheDocument()
    expect(badge).toHaveClass('bg-primary')
  })

  it('renders with different variants', () => {
    const { rerender } = render(<Badge variant="secondary">Secondary</Badge>)
    expect(screen.getByText('Secondary')).toHaveClass('bg-secondary')

    rerender(<Badge variant="destructive">Destructive</Badge>)
    expect(screen.getByText('Destructive')).toHaveClass('bg-destructive')

    rerender(<Badge variant="outline">Outline</Badge>)
    expect(screen.getByText('Outline')).toHaveClass('border')

    rerender(<Badge variant="success">Success</Badge>)
    expect(screen.getByText('Success')).toHaveClass('bg-green-500')

    rerender(<Badge variant="warning">Warning</Badge>)
    expect(screen.getByText('Warning')).toHaveClass('bg-yellow-500')

    rerender(<Badge variant="info">Info</Badge>)
    expect(screen.getByText('Info')).toHaveClass('bg-blue-500')
  })

  it('applies custom className', () => {
    render(<Badge className="custom-class">Custom Badge</Badge>)
    
    const badge = screen.getByText('Custom Badge')
    expect(badge).toHaveClass('custom-class')
  })

  it('renders with different content types', () => {
    const { rerender } = render(<Badge>Text Content</Badge>)
    expect(screen.getByText('Text Content')).toBeInTheDocument()

    rerender(
      <Badge>
        <span>HTML Content</span>
      </Badge>
    )
    expect(screen.getByText('HTML Content')).toBeInTheDocument()

    rerender(<Badge>{123}</Badge>)
    expect(screen.getByText('123')).toBeInTheDocument()
  })

  it('has correct accessibility attributes', () => {
    render(<Badge>Accessible Badge</Badge>)
    
    const badge = screen.getByText('Accessible Badge')
    
    // Badge should be inline and not focusable by default
    expect(badge.tagName).toBe('DIV')
    expect(badge).not.toHaveAttribute('tabIndex')
  })

  it('can be made interactive', () => {
    const handleClick = jest.fn()
    render(
      <Badge 
        onClick={handleClick}
        className="cursor-pointer"
        role="button"
        tabIndex={0}
      >
        Clickable Badge
      </Badge>
    )
    
    const badge = screen.getByRole('button', { name: 'Clickable Badge' })
    expect(badge).toBeInTheDocument()
    expect(badge).toHaveClass('cursor-pointer')
  })

  it('renders with icons', () => {
    render(
      <Badge>
        <span>✓</span>
        Success Badge
      </Badge>
    )
    
    expect(screen.getByText('Success Badge')).toBeInTheDocument()
    expect(screen.getByText('✓')).toBeInTheDocument()
  })
})