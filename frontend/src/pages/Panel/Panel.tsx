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
import { PopupToast } from '@/components/Popup/Popup';
import { PopupStyle } from '@/components/Popup/Popup.props';



const MAX_FILE_SIZE_MB = 5;
const Panel: FC = () => {
  const [configs, setConfigs] = useState<BotConfig[] | undefined>();
  const [selected, setSelected] = useState<string>("");
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const [showToast, setShowToast] = useState<boolean>(false);
  const [toastText, setToastText] = useState<{ title: string; text: string; type: PopupStyle }>({
    title: "",
    text: "",
    type: PopupStyle.Info,
  });

  const navigate = useNavigate();

  useEffect(() => {
    const fetchConfigs = async () => {
      const [ok, data] = await ConfigService.getConfigs();
      if (data.Configs) setConfigs(data.Configs);
      else setConfigs(undefined);
    };
    fetchConfigs();
  }, []);

  const handleFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > MAX_FILE_SIZE_MB * 1024 * 1024) {
      resetToast("Ошибка файла", `Файл не должен превышать ${MAX_FILE_SIZE_MB} МБ`, PopupStyle.Error);
      return;
    }

    setFile(file);
  };

  const resetToast = (title: string, text: string, type: PopupStyle) => {
    setToastText({ title, text, type });
    setShowToast(true);
  };

  const startConfig = async () => {
    if (!selected) return;

    const selectedConfig = configs?.find((config) => config.Id === selected);
    if (!selectedConfig) {
      resetToast("Ошибка", "Выбранный конфиг не найден", PopupStyle.Error);
      return;
    }

    if (!file) {
      resetToast("Ошибка запуска", "Выберите файл до 5мб", PopupStyle.Error);
      return;
    }

    setLoading(true);
    try {
      await ConfigService.startConfig({
        configId: selectedConfig.Id,
        data: file,
        prompt: selectedConfig.prompt,
        userId: 0,
      });
      resetToast("Бот запущен", `Bot config ID ${selected}`, PopupStyle.Success);
      setFile(null);
    } catch (err) {
      console.error("Ошибка запуска бота:", err);
      resetToast("Ошибка запуска", "Не удалось запустить бота", PopupStyle.Error);
    } finally {
      setLoading(false);
    }
  };
    return (
        <div className={"w-full h-full flex justify-center items-center"}>
            {showToast && (
                <PopupToast title={toastText.title} style={toastText.type} onClose={() => setShowToast(false)}>
                {toastText.text}
                </PopupToast>
            )}
            <div className={"flex w-full h-full flex-col items-center gap-3"}>
                <div className={"flex w-full justify-center items-center gap-2"}>
                    <div className={"h-full"}>
                        <div>
                            <Card title={"Автопостинг"}>
                                <div className={"flex flex-col gap-4"}>
                                    <div className={"flex justify-between items-center"}>
                                        <span className={"font-bold text-xl"}>Выбери Конфиг</span>
                                        <SelectMenu
                                            value={selected}
                                            onChange={setSelected}
                                            options={configs?.map((config) => ({ label: config.name, value: config.Id }))}
                                        />
                                    </div>
                                </div>
                                <div className={"flex justify-between items-center gap-5"}>
                                    <Button styleType={ButtonStyleType.submit} onClick={() => navigate(ADD_CONFIG_ROUTE)}>Создать конфиг</Button>
                                    {!loading && !file && (
                                        <FileInput
                                            accept=".xls,.xlsx,application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
                                            onChange={handleFile}
                                        />
                                    )}
                                    {loading && !file && <span className="loader" />}

                                    {file && (
                                    <Button styleType={ButtonStyleType.submit} onClick={() => setFile(null)}>
                                        Сбросить
                                    </Button>
                                    )}

                                    <Button
                                        styleType={ButtonStyleType.submit}
                                        disabled={!selected || loading}
                                        onClick={startConfig}
                                    >
                                        Запустить
                                    </Button>
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