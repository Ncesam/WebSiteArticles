import React, { useState } from "react";
import type { FC } from "react";
import { RegistrationProps } from "./Registration.props";
import Input from "@/ui/Input/Input";
import { InputStyleType } from "@/ui/Input/Input.props";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import { UserService } from "@/http/User";
import { useNavigate } from "react-router-dom";
import { LOGIN_ROUTE } from "@/utils/consts";

const Registration: FC<RegistrationProps> = ({ }) => {
    const [errors, setErrors] = useState<{ email?: string, nickname?: string, password?: string }>();
    const [email, setEmail] = useState<string>("");
    const [nickname, setNickname] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [isValid, setIsValid] = useState<boolean>();
    const [isLoading, setIsLoading] = useState<boolean>();
    const navigate = useNavigate();
    const validate = (data: any) => {
        const fieldErrors: { email?: string, nickname?: string, password?: string } = {};
        switch (data.message) {
            case "Fill fields":
                fieldErrors["password"] = "Заполните";
                fieldErrors["nickname"] = "Заполните";
                fieldErrors["email"] = "Заполните";
                setErrors(fieldErrors);
                return false;
            case "Password not validate":
                fieldErrors["password"] = "Длина пароля от 8 символов";
                setErrors(fieldErrors);
                return false;
            case "email already exists":
                fieldErrors['email'] = "Пользователь с такой почтой существует"
                setErrors(fieldErrors);
                return false
            case "nickname already exists":
                fieldErrors['nickname'] = "Пользователь с таким именем существует"
                setErrors(fieldErrors);
                return false;
            default:
                setErrors({});
                return true;
        }
    };
    const register = async () => {
        const [ok, err] = await UserService.register(email, nickname, password);
        if (!ok) {
            validate(err)
            return;
        }
        navigate(LOGIN_ROUTE);
    }
    return (
        <div className={"flex-1 flex h-screen items-center justify-center"}>
            <div className={"flex flex-col items-center gap-2"}>
                <div>
                    <h1 className={"text-2xl font-semibold text-base-darkBlue tracking-tight leading-snug"}>
                        Регистрация
                    </h1>
                </div>
                <Input style={InputStyleType.login} helperText={errors?.email} error={errors?.email ? true : undefined} onChange={(e) => setEmail(e.target.value)} placeholder={"Почта"} />
                <Input style={InputStyleType.login} helperText={errors?.nickname} error={errors?.nickname ? true : undefined} onChange={(e) => setNickname(e.target.value)} placeholder={"Логин"} />
                <Input style={InputStyleType.login} helperText={errors?.password} error={errors?.password ? true : undefined} onChange={(e) => setPassword(e.target.value)} placeholder={"Пароль"} />
                <Input style={InputStyleType.login} onChange={(e) => password === e.target.value ? setIsValid(false) : setIsValid(true)} helperText={isValid ? "Пароли не совпадают" : undefined} placeholder={"Повторите Пароль"} />
                <div className={"flex items-center justify-between gap-2"}>
                    <Button styleType={ButtonStyleType.submit} onClick={register} disabled={isValid}>Зарегистрироваться</Button>
                </div>
            </div>
        </div>
    );
};

export default Registration;

