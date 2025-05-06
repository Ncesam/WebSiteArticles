import { ChangeEventHandler } from "react";

export interface TextAreaProps {
    helperText?: string;
    error?: boolean;
    disabled?: boolean;
    style?: string;
    placeholder?: string;
    onChange?: ChangeEventHandler<HTMLTextAreaElement>;
    value?: string
}