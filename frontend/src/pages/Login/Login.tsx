import type { FC } from "react";
import React, { useState } from "react";
import { LoginProps } from "./Login.props";
import { AxiosResponse } from "axios";
import Input from "@/ui/Input/Input";
import { InputStyleType } from "@/ui/Input/Input.props";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import { UserService } from "@/http/User";
import { useNavigate } from "react-router-dom";
import { PANEL_ROUTE } from "@/utils/consts";
import { useStore } from "@/hooks/store";
import { IUser } from "@/types/user";

const Login: FC<LoginProps> = ({ }) => {
    const [errors, setErrors] = useState<{ nickname?: string; password?: string }>({});
    const {userStore} = useStore();
    const [nickname, setNickname] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const navigate = useNavigate();
    const validate = (data: any) => {
        const fieldErrors: { nickname?: string, password?: string } = {};
        switch (data.message) {
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
    const [isLoading, setIsLoading] = useState<boolean>();
    const login = async () => {
        setIsLoading(true)
        const [ok, data] = await UserService.login(nickname, password)
        setIsLoading(false)
        if (!ok) {
            validate(data)
            return
        }
        userStore.SetIsAuth(true)
        const user: IUser = {
            email:  data?.email,
            nickname: nickname,
            password: password
        }
        userStore.SetUser(user);
        navigate(PANEL_ROUTE);
        console.log(userStore);
    }
    return (
        <div className={"flex-1 flex h-screen items-center justify-center"}>
            <div className={"flex flex-col items-center gap-2"}>
                <div>
                    <h1 className={"text-2xl font-semibold text-base-darkBlue tracking-tight leading-snug"}>
                        Вход
                    </h1>
                </div>
                <Input style={InputStyleType.login} helperText={errors.nickname} onChange={(e) => setNickname(e.target.value)} error={errors.nickname ? true : undefined} placeholder={"Логин"} />
                <Input style={InputStyleType.login} helperText={errors.password} onChange={(e) => setPassword(e.target.value)} error={errors.password ? true : undefined} placeholder={"Пароль"} />
                <div className={"flex items-center justify-between gap-2"}>
                    <Button styleType={ButtonStyleType.submit} loading={isLoading} onClick={login}>Войти</Button>
                </div>
            </div>
        </div>
    );
};

export default Login;

