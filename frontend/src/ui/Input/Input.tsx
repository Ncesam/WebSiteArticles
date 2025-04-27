import React from "react";
import type { FC } from "react";
import { InputProps } from "./Input.props";
import { clsx } from "clsx";

const Input: FC<InputProps> = ({
    style,
    Icon,
    onChange,
    placeholder,
    type = "text",
    value,
    helperText,
    error,
    disabled
}) => {
    return (
        <div className={"flex flex-col w-full gap-1"}>
            <div className={
                clsx(
                    "flex items-center px-4 py-2 rounded-xl transition-all duration-300 ease-in-out",
                    "bg-base-darkBlue",
                    {
                        "border-red-400 animate-shake": error,
                        "bg-base-dark": disabled,
                        "cursor-not-allowed opacity-60": disabled,
                        "focus-within:ring-2 focus-within:ring-base-lightBlue":
                            !disabled && !error,
                    })
            }>
                {Icon && <Icon className={"text-base-grayBlue w-4 h-4 mr-2"} />}
                <input type={type} className={style} onChange={onChange} placeholder={placeholder} value={value} />
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

export default Input;

