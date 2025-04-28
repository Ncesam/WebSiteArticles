import { BotConfig } from "@/types/article";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import Card from "@/ui/Card/Card";
import SelectMenu from "@/ui/SelectMenu/SelectMenu";
import { FC, useState } from "react";


const Settings: FC = () => {
    const [status, setStatus] = useState<boolean>(false);
    const [configs, setConfigs] = useState<BotConfig[]>(
        [{label: "876", delay: 4, id: 4, promptId: 5, themeId: 5, userId: 5}]
    );
    return (
        <div className={"flex justify-center items-center gap-3"}>
            <Card title={"Автопостинг"}>
                <div className={"flex flex-col gap-4"}>
                    <div className={"flex justify-between items-between"}>
                        <span className={"font-bold text-xl inline-block"}>Выбери Конфиг</span>
                        <SelectMenu options={configs?.map((config) => ({label: config.label, value: config.id.toString()}))}></SelectMenu>
                    </div>
                    <span className="">Статус: {status ? "Активен" : "Не активен"}</span>
                </div>
                <div>
                    <Button styleType={ButtonStyleType.submit}>Создать конфиг</Button>
                </div>
            </Card>
            <Card title={"Общие настройки"}>
                <div></div>
            </Card>
        </div>
    )
}

export default Settings;