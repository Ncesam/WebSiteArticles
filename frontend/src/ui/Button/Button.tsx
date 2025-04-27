import React, {useEffect, useState} from "react";
import type {FC} from "react";
import {ButtonProps, ButtonStyleType} from "./Button.props";
import {clsx} from "clsx";

const Button: FC<ButtonProps> = ({styleType, onClick, loading = false, disabled = false, children}) => {
    const [isDisabled, setIsDisabled] = useState<boolean>(disabled || loading);

    useEffect(() => {
        setIsDisabled(disabled || loading);
    }, [disabled, loading]);
    return (
        <button
            disabled={isDisabled}
            onClick={onClick}
            className={clsx(
                "transition transform hover:scale-105 ease-in-out duration-200 rounded-xl px-4 py-2 font-medium",
                styleType,
                isDisabled && ButtonStyleType.disabled
            )}
        >

            {children}
        </button>
    );
};

export default Button;

