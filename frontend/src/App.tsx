import { FC } from "react";
import { BrowserRouter } from "react-router-dom";
import AppRoutes from "@/AppRoutes";
import Header from "@/components/Header/Header";
import { useAuth, useAutoLogin } from "@/hooks/auth";


const App: FC = () => {
    return (
        <BrowserRouter>
            <AuthWrapper />
        </BrowserRouter>
    );
};

const AuthWrapper: FC = () => {
    useAuth();
    useAutoLogin();

    return (
        <div className={"w-full h-screen flex flex-col items-center gap-2"}>
            <Header />
            <AppRoutes />
        </div>
    );
};



export default App;