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
import * as XLSX from "xlsx";



const Panel: FC = () => {
    const [status, setStatus] = useState<boolean>(false);
    const [configs, setConfigs] = useState<BotConfig[]>();
    const [loading, setLoading] = useState<boolean>();
    const navigate = useNavigate();
    const [selected, setSelected] = useState<string>();
    const [file, setFile] = useState<File | null>();
    const [helperText, setHelperText] = useState<{ item: string, text: string }>();
    const handleFile = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        setLoading(true);
        const reader = new FileReader();

        reader.onload = (event) => {
            const data = new Uint8Array(event.target?.result as ArrayBuffer);
            const workbook = XLSX.read(data, { type: 'array' });

            const firstSheetName = workbook.SheetNames[0];
            const worksheet = workbook.Sheets[firstSheetName];


            const rows: any[][] = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

            if (rows.length === 0) return;

            const headers = rows[0] as string[];

            const dataObjects = rows.slice(1).map((row) => {
                const obj: Record<string, any> = {};
                headers.forEach((header, index) => {
                    obj[header] = row[index];
                });
                return obj;
            });

            const jsonString = JSON.stringify(dataObjects);

            console.log('JSON строка:', jsonString);
        };

        reader.readAsArrayBuffer(file);
        setLoading(false);
    };
    const startConfig = async () => {
        console.log("start")
    }
    useEffect(() => {
        const fetchConfigs = async () => {
            const [ok, data] = await ConfigService.getConfigs();
            if (!ok) {
                console.log("error fetch configs: %s", data.message)
                return
            }
            setConfigs(data.Configs);
        }
        fetchConfigs();
    }, [])
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
                                        <SelectMenu value={selected} onChange={(value) => setSelected(value)} options={configs?.map((config) => ({ label: config.name, value: config.Id }))}></SelectMenu>
                                    </div>
                                    <span className={"font-semibold text-base-grayBlue  "}>Статус: {status ? "Активен" : "Не активен"}</span>
                                </div>
                                <div className={"flex justify-between items-center gap-5"}>
                                    <Button styleType={ButtonStyleType.submit} onClick={() => navigate(ADD_CONFIG_ROUTE)}>Создать конфиг</Button>
                                    <div>
                                        {!loading ? <FileInput accept=".xls,.xlsx,application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" onChange={handleFile} /> : <span className="loader"></span>}
                                    </div>
                                    <div>
                                        <Button styleType={ButtonStyleType.submit} disabled={!selected && loading} onClick={startConfig}>Запустить</Button>
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