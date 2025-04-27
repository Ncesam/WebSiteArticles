import Editor from "@/components/Editor/Editor";
import Card from "@/ui/Card/Card";
import { FC, useState } from "react";


const AddArticle: FC = () => {
    const [title, setTitle] = useState('')
    const [cover, setCover] = useState<File | null>(null)
    const [data, setData] = useState<any>(null)
    return (
        <div className={"w-1/2"}>
            <Card title="Статья">
                ...
            </Card>
        </div>
    )
}
export default AddArticle;