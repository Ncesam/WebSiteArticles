import { Dispatch } from "react"


export interface BoolInputProps {
  value: boolean;
  onChange: (value: boolean) => void;
  disabled?: boolean;
  helperText?: string;
  error?: boolean;
  Icon?: React.ComponentType<{ className?: string }>;
};