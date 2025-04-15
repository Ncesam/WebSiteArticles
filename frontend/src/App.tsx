import {FC, useContext} from "react";
import {BrowserRouter, Routes} from "react-router-dom";


const App: FC = () => {
    const context = useContext(Context)
    return (
        <BrowserRouter>
            <Routes>
                {}
                {}
            </Routes>
        </BrowserRouter>
    )
}


export default App;