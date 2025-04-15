import {FC, useContext} from "react";
import {BrowserRouter, Routes} from "react-router-dom";
import AppRoutes from "./AppRoutes";


const App: FC = () => {
    return (
        <BrowserRouter>
            <AppRoutes/>
        </BrowserRouter>
    )
}


export default App;