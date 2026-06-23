import { cloneElement, isValidElement, type ReactElement, type ReactNode } from 'react'
import { joinClassNames } from './types'

type FieldControlProps = {
  id?: string
  required?: boolean
  'aria-describedby'?: string
  'aria-invalid'?: boolean | 'true' | 'false'
}

type FieldProps = {
  id: string
  label: string
  hint?: ReactNode
  error?: ReactNode
  required?: boolean
  className?: string
  children: ReactElement<FieldControlProps>
}

export function Field({ id, label, hint, error, required = false, className, children }: FieldProps) {
  const hintID = hint ? `${id}-hint` : undefined
  const errorID = error ? `${id}-error` : undefined
  const describedBy = [children.props['aria-describedby'], hintID, errorID].filter(Boolean).join(' ') || undefined
  const control = isValidElement<FieldControlProps>(children)
    ? cloneElement(children, {
      id,
      required,
      'aria-describedby': describedBy,
      'aria-invalid': error ? 'true' : undefined,
    })
    : children

  return (
    <div className={joinClassNames('ui-field', className)}>
      <label className="ui-field-label" htmlFor={id}>
        {label}
      </label>
      {control}
      {hint ? <p id={hintID} className="ui-field-hint">{hint}</p> : null}
      {error ? <p id={errorID} className="ui-field-error" role="alert">{error}</p> : null}
    </div>
  )
}
