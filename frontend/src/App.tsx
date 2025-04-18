import React, {FC} from "react";
import {BrowserRouter} from "react-router-dom";
import Button from "./ui/Button/Button";
import {ButtonStyleType} from "@/ui/Button/Button.props";
import AppRoutes from "@/AppRoutes";
import Header from "@/components/Header/Header";


const App: FC = () => {
    return (
        <BrowserRouter>
            <div className={"h-screen flex flex-col items-center gap-2"}>
                <div className={"w-full"}>
                    <Header/>
                </div>
                <AppRoutes/>
            </div>
        </BrowserRouter>
    )
}


export default App;