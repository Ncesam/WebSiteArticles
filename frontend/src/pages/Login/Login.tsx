import type {FC} from "react";
import React, {useState} from "react";
import {LoginProps} from "./Login.props";
import {AxiosResponse} from "axios";
import Input from "@/ui/Input/Input";
import {InputStyleType} from "@/ui/Input/Input.props";
import Button from "@/ui/Button/Button";
import {ButtonStyleType} from "@/ui/Button/Button.props";

const Login: FC<LoginProps> = ({}) => {
    const [errors, setErrors] = useState<{ nickname?: string; password?: string }>({});
    const [nickname, setNickname] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const validate = (response: AxiosResponse) => {
        const fieldErrors: { nickname?: string, password?: string } = {};
        switch (response.data.message) {
            case "Password is incorrect":
                fieldErrors["password"] = "Пароль неправильный";
                setErrors(fieldErrors);
                return false;
            case "User not found":
                fieldErrors["nickname"] = "Пользователь не найден";
                setErrors(fieldErrors);
                return false;
            default:
                setErrors({});
                return true;
        }
    };
    return (

        <div className={"flex-1 flex h-screen items-center justify-center"}>
            <div className={"flex flex-col items-center gap-2"}>
                <div>
                    <h1 className={"text-2xl font-semibold text-base-darkBlue tracking-tight leading-snug"}>
                        Вход
                    </h1>
                </div>
                <Input style={InputStyleType.login} placeholder={"Логин"}/>
                <Input style={InputStyleType.login} placeholder={"Пароль"}/>
                <div className={"flex items-center justify-between gap-2"}>
                    <Button styleType={ButtonStyleType.submit} disabled>Войти</Button>
                </div>
            </div>
        </div>
    );
};

export default Login;

