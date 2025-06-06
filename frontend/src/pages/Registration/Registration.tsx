import { FC, useState } from "react";
import { useNavigate } from "react-router-dom";

import { RegistrationProps } from "./Registration.props";
import Input from "@/ui/Input/Input";
import { InputStyleType, InputType } from "@/ui/Input/Input.props";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import { UserService } from "@/http/User";
import { LOGIN_ROUTE } from "@/utils/consts";

const Registration: FC<RegistrationProps> = () => {
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [nickname, setNickname] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const [errors, setErrors] = useState<{
    email?: string;
    nickname?: string;
    password?: string;
    confirmPassword?: string;
  }>({});

  const validate = (message: string): boolean => {
    const fieldErrors: typeof errors = {};

    switch (message) {
      case "Fill fields":
        fieldErrors.email = "Заполните почту";
        fieldErrors.nickname = "Заполните логин";
        fieldErrors.password = "Заполните пароль";
        break;
      case "Password not validate":
        fieldErrors.password = "Длина пароля должна быть от 8 символов";
        break;
      case "email already exists":
        fieldErrors.email = "Пользователь с такой почтой уже существует";
        break;
      case "nickname already exists":
        fieldErrors.nickname = "Пользователь с таким именем уже существует";
        break;
      default:
        break;
    }

    setErrors(fieldErrors);
    return Object.keys(fieldErrors).length === 0;
  };

  const handleRegister = async () => {
    if (password !== confirmPassword) {
      setErrors((prev) => ({
        ...prev,
        confirmPassword: "Пароли не совпадают",
      }));
      return;
    }

    const [ok, err] = await UserService.register(email, nickname, password);
    if (!ok) {
      validate(err.message);
      return;
    }

    navigate(LOGIN_ROUTE);
  };

  return (
    <div className="flex-1 flex items-center justify-center">
      <div className="w-full flex-col flex items-center gap-4">
        <h1 className={"text-2xl font-semibold text-base-darkBlue tracking-tight leading-snug"}>Регистрация</h1>

        <Input
          style={InputStyleType.login}
          placeholder="Почта"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          helperText={errors.email}
          error={!!errors.email}
        />

        <Input
          style={InputStyleType.login}
          placeholder="Логин"
          value={nickname}
          onChange={(e) => setNickname(e.target.value)}
          helperText={errors.nickname}
          error={!!errors.nickname}
        />

        <Input
          style={InputStyleType.login}
          placeholder="Пароль"
          type={InputType.password}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          helperText={errors.password}
          error={!!errors.password}
        />

        <Input
          style={InputStyleType.login}
          placeholder="Повторите пароль"
          type={InputType.password}
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          helperText={errors.confirmPassword}
          error={!!errors.confirmPassword}
        />
        <Button
        styleType={ButtonStyleType.submit}
        onClick={handleRegister}
        disabled={!email || !nickname || !password || !confirmPassword}
        >
        Зарегистрироваться
        </Button>
        
      </div>
    </div>
  );
};

export default Registration;
