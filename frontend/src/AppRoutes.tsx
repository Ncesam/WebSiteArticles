import React, { FC, Suspense } from "react";
import { useStore } from "./hooks/store";
import { Navigate, Route, Routes } from "react-router-dom";
import RootStore from "./stores/rootStore";
import { LOGIN_ROUTE, PANEL_ROUTE } from "./utils/consts";
import { authRoutes, primaryRoutes } from "./routes";


const AppRoutes: FC = () => {
    const context: RootStore = useStore()
    return (
        <Suspense fallback={<div></div>}>
            <Routes>
                {context.userStore.IsAuth ? (
                    <>
                        {authRoutes.map(({ path, component }) => (
                            <Route key={path} path={path} element={React.createElement(component)} />
                        ))}
                        <Route path="*" element={<Navigate to={LOGIN_ROUTE} />} />
                    </>
                ) : (
                    <>
                        {primaryRoutes.map(({ path, component }) => (
                            <Route key={path} path={path} element={React.createElement(component)} />
                        ))}
                        <Route path="*" element={<Navigate to={PANEL_ROUTE} />} />
                    </>
                )}
            </Routes>
        </Suspense>
    )
}


export default AppRoutes;