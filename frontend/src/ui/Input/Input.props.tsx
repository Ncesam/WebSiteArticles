import {ChangeEventHandler, ReactNode, SVGProps} from "react";

export interface InputProps {
    type?: InputType;
    style: InputStyleType;
    placeholder?: string;
    value?: any;
    onChange?: ChangeEventHandler<HTMLInputElement>;
    disabled?: boolean;
    error?: boolean;
    helperText?: string;
    Icon?: any;
}


export enum InputType {
    text = "text",
    password = "password",
    email = "email",
}
export enum InputStyleType {
    login = "flex-1 bg-transparent text-base-grayBlue placeholder-base-grayBlue outline-none"

}