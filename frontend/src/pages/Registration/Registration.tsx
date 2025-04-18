import React from "react";
import type {FC} from "react";
import {RegistrationProps} from "./Registration.props";
import Input from "@/ui/Input/Input";
import {InputStyleType} from "@/ui/Input/Input.props";
import Button from "@/ui/Button/Button";
import {ButtonStyleType} from "@/ui/Button/Button.props";

const Registration: FC<RegistrationProps> = ({}) => {
    return (
        <div className={"flex-1 flex h-screen items-center justify-center"}>
            <div className={"flex flex-col items-center gap-2"}>
                <div>
                    <h1 className={"text-2xl font-semibold text-base-darkBlue tracking-tight leading-snug"}>
                        Регистрация
                    </h1>
                </div>
                <Input style={InputStyleType.login} placeholder={"Почта"}/>
                <Input style={InputStyleType.login} placeholder={"Логин"}/>
                <Input style={InputStyleType.login} placeholder={"Пароль"}/>
                <Input style={InputStyleType.login} placeholder={"Повторите Пароль"}/>
                <div className={"flex items-center justify-between gap-2"}>
                    <Button styleType={ButtonStyleType.submit} disabled>Зарегистрироваться</Button>
                </div>
            </div>
        </div>
    );
};

export default Registration;

