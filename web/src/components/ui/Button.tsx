import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { joinClassNames } from './types'

export type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'danger' | 'link'

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant
  isWorking?: boolean
  workingLabel?: string
  children: ReactNode
}

export function Button({ variant = 'primary', isWorking = false, workingLabel = '处理中', className, disabled, children, type = 'button', ...props }: ButtonProps) {
  const label = isWorking ? workingLabel : children
  return (
    <button
      {...props}
      type={type}
      className={joinClassNames('ui-button', `ui-button-${variant}`, className)}
      disabled={disabled || isWorking}
      aria-busy={isWorking ? 'true' : undefined}
    >
      {label}
    </button>
  )
}
