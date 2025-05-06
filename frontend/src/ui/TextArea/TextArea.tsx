import clsx from "clsx";
import { FC, useEffect, useRef } from "react";
import { TextAreaProps } from "./TextArea.props";

const TextArea: FC<TextAreaProps> = ({
    style,
    onChange,
    placeholder,
    value,
    helperText,
    error,
    disabled,
}) => {
    const textareaRef = useRef<HTMLTextAreaElement>(null);
    return (
        <div className="h-full flex flex-col w-full gap-1">
            <div
                className={clsx(
                    "px-4 py-2 w-full h-5/6 rounded-xl transition-all duration-300 ease-in-out",
                    "bg-base-darkBlue",
                    {
                        "border-red-400 animate-shake": error,
                        "bg-base-dark": disabled,
                        "cursor-not-allowed opacity-60": disabled,
                        "focus-within:ring-2 focus-within:ring-base-lightBlue":
                            !disabled && !error,
                    }
                )}
            >


                <textarea
                    ref={textareaRef}
                    className={clsx(
                        "w-full h-full resize-none overflow-y-hidden bg-transparent text-white outline-none text-sm",
                        style
                    )}
                    placeholder={placeholder}
                    onChange={onChange}
                    value={value}
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

export default TextArea;