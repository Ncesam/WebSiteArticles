import React, { FC, Suspense } from "react";
import { useStore } from "@/hooks/store";
import { Navigate, Route, Routes } from "react-router-dom";
import { LOGIN_ROUTE, PANEL_ROUTE } from "@/utils/consts";
import { authRoutes, primaryRoutes } from "@/routes";
import { observer} from "mobx-react";


const AppRoutes: FC = observer(() => {
    const { userStore } = useStore()
    return (
        <Suspense fallback={<div className="flex w-full h-full justify-center items-center"><span className="loader"></span></div>}>
            {userStore.isLoading ? <div className="flex w-full h-full justify-center items-center"><span className="loader"></span></div> :
                <Routes>
                    {!userStore.IsAuth ? (
                        <>
                            {authRoutes.map(({ path, component }) => (
                                <Route key={path} path={path} element={React.createElement(component)} />
                            ))}
                            <Route path="*" element={<Navigate to={LOGIN_ROUTE} replace />} />
                        </>
                    ) : (
                        <>
                            {primaryRoutes.map(({ path, component }) => (
                                <Route key={path} path={path} element={React.createElement(component)} />
                            ))}
                            <Route path="*" element={<Navigate to={PANEL_ROUTE} replace />} />
                        </>
                    )}
                </Routes>}
        </Suspense>
    )
})


export default AppRoutes;