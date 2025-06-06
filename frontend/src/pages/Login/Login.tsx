import { FC, useState } from "react";
import { useNavigate } from "react-router-dom";

import { LoginProps } from "./Login.props";
import Input from "@/ui/Input/Input";
import { InputStyleType } from "@/ui/Input/Input.props";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";

import { UserService } from "@/http/User";
import { useStore } from "@/hooks/store";
import { IUser } from "@/types/user";
import { PANEL_ROUTE } from "@/utils/consts";

const Login: FC<LoginProps> = () => {
  const { userStore } = useStore();
  const navigate = useNavigate();

  const [nickname, setNickname] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState<{ nickname?: string; password?: string }>({});

  const validate = (message: string): boolean => {
    const fieldErrors: { nickname?: string; password?: string } = {};

    if (message === "Password is invalid") {
      fieldErrors.password = "Неверный пароль";
    } else if (message === "User not found") {
      fieldErrors.nickname = "Пользователь не найден";
    }

    setErrors(fieldErrors);
    return Object.keys(fieldErrors).length === 0;
  };

  const handleLogin = async () => {
    const [ok, data] = await UserService.login(nickname, password);
    console.log(data)
    if (!ok) {
      validate(data.message);
      return;
    }

    userStore.SetIsAuth(true);

    const user: IUser = {
      email: data.email,
      nickname,
      password,
    };

    userStore.SetUser(user);
    navigate(PANEL_ROUTE);
  };

     return (
        <div className={"flex-1 flex w-screen h-screen items-center justify-center"}>
            <div className={"w-1/6 flex flex-col items-center justify-center gap-2"}>
                <div>
                    <h1 className={"text-2xl font-semibold text-base-darkBlue tracking-tight leading-snug"}>
                        Вход
                    </h1>
                </div>
                <Input style={InputStyleType.login} helperText={errors.nickname} onChange={(e) => setNickname(e.target.value)} error={errors.nickname ? true : undefined} placeholder={"Логин"} />
                <Input style={InputStyleType.login} helperText={errors.password} onChange={(e) => setPassword(e.target.value)} error={errors.password ? true : undefined} placeholder={"Пароль"} />
                <div className={"w-full flex justify-center"}>
                    <Button styleType={ButtonStyleType.submit} onClick={handleLogin}>Войти</Button>
                </div>
            </div>
        </div>
    );
};

export default Login;
