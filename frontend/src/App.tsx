import React, { FC } from "react";
import { BrowserRouter } from "react-router-dom";
import Button from "./ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import AppRoutes from "@/AppRoutes";
import Header from "@/components/Header/Header";
import { useAuth, useAutoLogin } from "./hooks/auth";


const App: FC = () => {
    useAuth()
    useAutoLogin()
    return (
        <BrowserRouter>
            <div className={"w-full h-screen flex flex-col items-center gap-2"}>
                <Header />
                <AppRoutes />
            </div>
        </BrowserRouter>
    )
}


export default App;