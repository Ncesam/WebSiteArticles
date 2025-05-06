

export interface FileInputProps {
    onChange?: (file: File | null) => void;
    placeholder?: string;
    helperText?: string;
    error?: boolean;
    disabled?: boolean;
    style?: string;
    accept?: string;
}
