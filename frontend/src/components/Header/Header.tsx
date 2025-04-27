import { FC, useEffect } from "react";
import React from "react";
import { HeaderProps } from "./Header.props";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import { useStore } from "@/hooks/store";
import { useNavigate } from "react-router-dom";
import { ADD_ARTICLE_ROUTE, DASHBOARD_ROUTE, LOGIN_ROUTE, REGISTER_ROUTE, SETTINGS_ROUTE } from "@/utils/consts";

const Header: FC<HeaderProps> = () => {
    const { userStore } = useStore();
    const navigate = useNavigate();
    useEffect(() => {
    }, [userStore]);
    return (
        <div className="w-full m-4">
            <div className={"flex items-center justify-end gap-4 mr-4"}>
            <Button styleType={ButtonStyleType.submit} onClick={() => navigate(DASHBOARD_ROUTE)}>DashBoard</Button>
            <Button styleType={ButtonStyleType.submit} onClick={() => navigate(ADD_ARTICLE_ROUTE)}>Написать статью</Button>
            <Button styleType={ButtonStyleType.submit} onClick={() => navigate(SETTINGS_ROUTE)}>Настройки</Button>
            {userStore.IsAuth ? (
                <div>
                    <Button styleType={ButtonStyleType.submit}>Выйти</Button>
                </div>
            ) : (
                <div className="flex gap-4 items-center">
                    <Button styleType={ButtonStyleType.submit} onClick={() => navigate(LOGIN_ROUTE)}>Войти</Button>
                    <Button styleType={ButtonStyleType.submit} onClick={() => navigate(REGISTER_ROUTE)}>Зарегистрироваться</Button>
                </div>
            )}
            </div>
        </div>
    );
};

export default Header;

