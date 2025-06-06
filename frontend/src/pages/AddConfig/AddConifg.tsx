import { FC, useState } from "react";
import { useNavigate } from "react-router-dom";

import { ConfigService } from "@/http/Config";
import { FormBotConfig } from "@/types/Config";

import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import Card from "@/ui/Card/Card";
import Input from "@/ui/Input/Input";
import { InputStyleType, InputType } from "@/ui/Input/Input.props";
import TextArea from "@/ui/TextArea/TextArea";

import { PANEL_ROUTE } from "@/utils/consts";

const AddConfig: FC = () => {
  const [name, setName] = useState<string>("");
  const [prompt, setPrompt] = useState<string>("");
  const [delay, setDelay] = useState<number>(10);
  const [emailAccount, setEmailAccount] = useState<string>("");
  const [passwordAccount, setPasswordAccount] = useState<string>("");

  const navigate = useNavigate();

  const handleCreate = async () => {
    const config: FormBotConfig = {
      name,
      prompt,
      delay,
      email: emailAccount,
      password: passwordAccount,
    };

    const [ok, data] = await ConfigService.addConfig(config);
    if (!ok) {
      console.error("Ошибка создания конфига:", data);
      return;
    }

    navigate(PANEL_ROUTE);
  };

  return (
    <div className="w-full flex-1 flex items-center justify-center p-4">
      <div className="w-full max-w-6xl grid gap-6 grid-cols-1 md:grid-cols-2">
        <Card title="Название конфига">
          <div className="flex flex-col gap-4">
            <Input
              style={InputStyleType.login}
              placeholder="Введите название"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />

            <Input
              style={InputStyleType.login}
              type={InputType.text}
              placeholder="Время между постами"
              helperText="В часах"
              value={delay}
              onChange={(e) => setDelay(Number(e.target.value))}
            />

            <Input
              style={InputStyleType.login}
              placeholder="Email"
              helperText="Данные аккаунта"
              value={emailAccount}
              onChange={(e) => setEmailAccount(e.target.value)}
            />

            <Input
              style={InputStyleType.login}
              placeholder="Пароль"
              value={passwordAccount}
              onChange={(e) => setPasswordAccount(e.target.value)}
            />

            <Button styleType={ButtonStyleType.submit} onClick={handleCreate}>
              Создать
            </Button>
          </div>
        </Card>

        <Card
          title="Промпт"
          subtitle="Введите промпт. Используйте {name} как плейсхолдер для переменных из файла."
        >
          <TextArea
            placeholder="Введите текст..."
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
          />
        </Card>
      </div>
    </div>
  );
};

export default AddConfig;
