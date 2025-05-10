import { ConfigService } from "@/http/Config";
import { BotConfig } from "@/types/Config";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import Card from "@/ui/Card/Card";
import FileInput from "@/ui/FileInput/FileInput";
import SelectMenu from "@/ui/SelectMenu/SelectMenu";
import { ADD_CONFIG_ROUTE } from "@/utils/consts";
import { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";



const Panel: FC = () => {
    const [status, setStatus] = useState<boolean>(false);
    const [configs, setConfigs] = useState<BotConfig[]>(
        [{ id: 1, name: "52", delay: 43, prompt: "fgfd", userId: 5 }]
    );
    const navigate = useNavigate();
    const [selected, setSelected] = useState<string>();
    const [file, setFile] = useState<File | null>();
    const [helperText, setHelperText] = useState<{ item: string, text: string }>();
    const startConfig = async () => {
        console.log("start")
    }
    useEffect(() => {
        const fetchConfigs = async () => {
            
        }
    })
    return (
        <div className={"w-full h-full flex justify-center items-center"}>
            <div className={"flex w-full h-full flex-col items-center gap-3"}>
                <div className={"flex w-full justify-center items-center gap-2"}>
                    <div className={"h-full"}>
                        <div>
                            <Card title={"Автопостинг"}>
                                <div className={"flex flex-col gap-4"}>
                                    <div className={"flex justify-between items-center"}>
                                        <span className={"font-bold text-xl"}>Выбери Конфиг</span>
                                        <SelectMenu value={selected} onChange={(value) => setSelected(value)} options={configs?.map((config) => ({ label: config.name, value: config.id.toString() }))}></SelectMenu>
                                    </div>
                                    <span className={"font-semibold text-base-grayBlue  "}>Статус: {status ? "Активен" : "Не активен"}</span>
                                </div>
                                <div className={"flex justify-between items-center gap-5"}>
                                    <Button styleType={ButtonStyleType.submit} onClick={() => navigate(ADD_CONFIG_ROUTE)}>Создать конфиг</Button>
                                    <div>
                                        <FileInput accept=".xls,.xlsx,application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" onChange={(file) => setFile(file)} />
                                    </div>
                                    <div>
                                        {!selected ? <Button styleType={ButtonStyleType.submit} disabled onClick={startConfig}>Запустить</Button> :
                                            <Button styleType={ButtonStyleType.submit} onClick={startConfig}>Запустить</Button>}
                                    </div>
                                </div>
                            </Card>
                        </div>
                    </div>

                </div>
            </div>
        </div >
    )
}

export default Panel;