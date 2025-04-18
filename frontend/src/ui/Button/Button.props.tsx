import {MouseEventHandler, ReactNode} from "react";

export interface ButtonProps {
    styleType: ButtonStyleType;
    onClick?: MouseEventHandler<HTMLButtonElement>;
    loading?: boolean;
    disabled?: boolean;
    children?: ReactNode | undefined;
}

export enum ButtonStyleType {
    submit = "flex items-center gap-2 px-6 py-2 rounded-xl bg-base-lightBlue text-white font-medium hover:bg-base-darkBlue outline-none",
    disabled = "bg-base-grayBlue text-white cursor-not-allowed opacity-60 hover:scale-100",
}

