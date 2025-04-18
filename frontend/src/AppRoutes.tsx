import React, {FC, Suspense} from "react";
import {useStore} from "./hooks/store";
import {Navigate, Route, Routes} from "react-router-dom";
import RootStore from "./stores/rootStore";
import {DASHBOARD_ROUTE, LOGIN_ROUTE} from "./utils/consts";
import {authRoutes, primaryRoutes} from "./routes";


const AppRoutes: FC = () => {
    const context: RootStore = useStore()
    return (
        <Suspense>
            <Routes>
                {!context.userStore.IsAuth ? (
                    <>
                        {authRoutes.map(({path, component}) => (
                            <Route key={path} path={path} element={React.createElement(component)}/>
                        ))}
                        <Route path="*" element={<Navigate to={LOGIN_ROUTE} replace/>}/>
                    </>
                ) : (
                    <>
                        {primaryRoutes.map(({path, component}) => (
                            <Route key={path} path={path} element={React.createElement(component)}/>
                        ))}
                        <Route path="*" element={<Navigate to={DASHBOARD_ROUTE} replace/>}/>
                    </>
                )}
            </Routes>
        </Suspense>
    )
}


export default AppRoutes;