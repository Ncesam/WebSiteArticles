import { FC, useState } from "react"
import { Pie, Line } from "react-chartjs-2";
import { Chart as ChartJS, LineElement, Legend, LinearScale, ArcElement, CategoryScale, Tooltip, Title, PointElement, ChartData, ChartOptions, Filler } from "chart.js";
import Card from "@/ui/Card/Card";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import { Articles } from "@/types/article";

ChartJS.register(
    Legend, LineElement,
    LinearScale, CategoryScale,
    Tooltip, Title, PointElement,
    Filler, ArcElement,
)

const generateRandomArray = (arrayLength: number) => {
    let result = []
    for (let i = 0; i < arrayLength; i++) {
        result.push(Math.round(Math.random() * 100))
    }
    return result
}
const DashBoard: FC = () => {
    const [articlesQueue, setArticlesQueue] = useState<Articles[]>([])
    const yData = generateRandomArray(10);
    const xData = generateRandomArray(10);
    const [articlesData, setArticlesData] = useState<ChartData<'line'>>({
        labels: xData,
        datasets: [{
            data: yData,
            tension: 0.2,
            borderColor: "#082D3A",
            borderWidth: 2
        }],
    }

    )
    const [categoryData, setCategoryData] = useState<ChartData<'pie'>>({
        labels: ['Спорт', "Игры", "Дети"],
        datasets: [{
            data: generateRandomArray(3),
            backgroundColor: [
                "#267491", "#082D3A", "#A7B8B5"
            ],
            offset: 10,
        }]
    })
    const LineOptions: ChartOptions<'line'> = {
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: false
            }
        },
        scales: {
            x: {
                min: 0,
                max: 20,
                ticks: {
                    autoSkip: true,
                    maxRotation: 0,
                    color: "#082D3A",
                },
                border: {
                    display: false
                },
                grid: {
                    display: false
                },
            },
            y: {
                min: 0,
                max: 110,
                ticks: {
                    color: "#082D3A",
                },
                border: {
                    color: "#267491",

                },
                grid: {
                    color: "#267491",

                },
            }
        },
        animation: {
            duration: 400,
            easing: "linear",
        },
    }
    const PieOptions: ChartOptions<"pie"> = {
        plugins: {
            legend: {
                labels: {
                    color: "#082D3A"
                },
                position: "right",
                align: "start"
            },
        },
    }
    return (
        <div className={"w-full h-full flex justify-center items-center"}>
            <div className={"flex w-full h-full flex-col items-center gap-3"}>
                <div className={"w-full h-1/2 flex justify-center items-center gap-3"}>
                    <div className={"w-1/2"}>
                        <Card title="Статистика отправки">
                            <Line data={articlesData} options={LineOptions}></Line>
                        </Card>
                    </div>
                    <div className={"w-1/5"}>
                        <Card>
                            <div className={"flex flex-col items-start gap-3"}>
                                <div className={"text-base-darkBlue font-semibold text-lg"}>
                                    Сейчас в очереди: {articlesQueue.length}
                                </div>
                                <Button styleType={ButtonStyleType.submit}>Написать Статью</Button>
                                <Button styleType={ButtonStyleType.submit}>Посмотреть очередь статей</Button>
                            </div>
                        </Card>
                    </div>
                </div>

                <div className={"flex w-full justify-center items-center gap-2"}>
                    <div className={"h-full"}>
                        <Card title="Статистика Категорий">
                            <Pie data={categoryData} options={PieOptions}></Pie>
                        </Card>
                    </div>
                    <div className={"h-full"}>
                        <div>
                            <Card title="Настройки" subtitle="Настройка автопостинга">
                                <div className={"grid grid-rows-2 grid-cols-2 gap-4"}>
                                    <Button styleType={ButtonStyleType.submit}>Темы</Button>
                                    <Button styleType={ButtonStyleType.submit}>Промпты</Button>
                                    <Button styleType={ButtonStyleType.submit}>Автопостинг</Button>
                                    <Button styleType={ButtonStyleType.submit}>Промпты</Button>
                                </div>
                            </Card>
                        </div>
                    </div>

                </div>
            </div>
        </div >
    )
}

export default DashBoard;

