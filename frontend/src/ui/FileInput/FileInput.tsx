import React, { FC, useRef, useState } from "react";
import { clsx } from "clsx";
import { FileInputProps } from "./FileInput.props";


const FileInput: FC<FileInputProps> = ({
    onChange,
    placeholder = "Выберите файл",
    helperText,
    error,
    disabled,
    style,
    accept
}) => {
    const inputRef = useRef<HTMLInputElement>(null);
    const [fileName, setFileName] = useState<string>("");

    const handleClick = () => {
        if (!disabled) inputRef.current?.click();
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        setFileName(file?.name || "");
        onChange?.(file || null);
    };

    return (
        <div className="flex flex-col w-full gap-1">
            <div
                onClick={handleClick}
                className={clsx(
                    "flex items-center px-4 py-2 rounded-xl transition-all duration-300 ease-in-out cursor-pointer",
                    "bg-base-darkBlue text-white",
                    {
                        "border-red-400 animate-shake": error,
                        "bg-base-dark": disabled,
                        "cursor-not-allowed opacity-60": disabled,
                        "focus-within:ring-2 focus-within:ring-base-lightBlue": !disabled && !error,
                    }
                )}
            >
                <span className={clsx("text-sm", style)}>
                    {fileName || placeholder}
                </span>
                <input
                    ref={inputRef}
                    type="file"
                    accept={accept}
                    className="hidden"
                    onChange={handleFileChange}
                    disabled={disabled}
                />
            </div>

            {helperText && (
                <p
                    className={clsx(
                        "text-xs ml-1",
                        error ? "text-red-400" : "text-base-lightBlue"
                    )}
                >
                    {helperText}
                </p>
            )}
        </div>
    );
};

export default FileInput;
