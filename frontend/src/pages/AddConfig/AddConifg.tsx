import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import Card from "@/ui/Card/Card";
import Input from "@/ui/Input/Input";
import { InputStyleType, InputType } from "@/ui/Input/Input.props";
import TextArea from "@/ui/TextArea/TextArea";
import { FC, useState } from "react";

const AddConfig: FC = () => {
    const [name, setName] = useState<string>();
    const [promptText, setPromptText] = useState<string>();
    const [delay, setDelay] = useState<number>();
    const [emailAccount, setEmailAccount] = useState<string>();
    const [passwordAccount, setPasswordAccount] = useState<string>();
    const click = async () => {
        
    }

    return (
        <div className={"w-full h-full flex justify-center items-center"}>
            <div className={"w-1/2 h-5/6 flex justify-center gap-3"}>
                <div className={"w-3/4 h-1/2"}>
                    <Card title="Название конфига" >
                        <div>
                            <Input style={InputStyleType.login} onChange={(e) => setName(e.target.value)} />
                        </div>
                        <div className={"flex justify-between items-end gap-2"}>
                            <div className={"flex gap-2 flex-col"}>
                                <Input style={InputStyleType.login} helperText="В часах" placeholder="Время между постами" type={InputType.text} value={delay} onChange={(e) => setDelay(Number(e.target.value))} />
                                <Input style={InputStyleType.login} placeholder="Email" helperText="Данные аккаунта" value={emailAccount} onChange={(e) => setEmailAccount(e.target.value)} />
                                <Input style={InputStyleType.login} placeholder="Пароль" value={passwordAccount} onChange={(e) => setPasswordAccount(e.target.value)} />
                            </div>
                            <Button styleType={ButtonStyleType.submit} onClick={click}>Создать</Button>
                        </div>
                    </Card>
                </div>
                <div className={" w-full h-full"}>
                    <Card title={"Промпт"} subtitle={"Введите промпт. Под переменные поставьте {name}. Учитывайте, что переменные берутся из ваших файлов."}>
                        <div className={"w-full h-5/6"}>
                            <TextArea onChange={(e) => setPromptText(e.target.value)} />
                        </div>
                    </Card>
                </div>
            </div>
        </div>
    )
}

export default AddConfig;