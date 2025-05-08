import { lazy } from "react";
import { ADD_CONFIG_ROUTE, LOGIN_ROUTE, PANEL_ROUTE, REGISTER_ROUTE,  } from "./utils/consts";

const Registration = lazy(() => import("@/pages/Registration/Registration"));
const Login = lazy(() => import("@/pages/Login/Login"))
const Panel = lazy(() => import("@/pages/Panel/Panel"))
const AddConfig = lazy(() => import("@/pages/AddConfig/AddConifg"))

export const authRoutes = [
    {
        path: LOGIN_ROUTE,
        component: Login
    },
    {
        path: REGISTER_ROUTE,
        component: Registration
    },
]
export const primaryRoutes = [
    {
        path: PANEL_ROUTE,
        component: Panel
    },
    {
        path: ADD_CONFIG_ROUTE,
        component: AddConfig
    }
]