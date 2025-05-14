import { ChangeEventHandler } from "react";


export interface FileInputProps {
    onChange?: ChangeEventHandler<HTMLInputElement>;
    placeholder?: string;
    helperText?: string;
    error?: boolean;
    disabled?: boolean;
    style?: string;
    accept?: string;
}
