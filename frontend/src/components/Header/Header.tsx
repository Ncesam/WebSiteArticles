import { FC, useEffect } from "react";
import React from "react";
import { HeaderProps } from "./Header.props";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import { useStore } from "@/hooks/store";
import { useNavigate } from "react-router-dom";
import { LOGIN_ROUTE, PANEL_ROUTE, REGISTER_ROUTE } from "@/utils/consts";
import { observer } from "mobx-react";
import { UserService } from "@/http/User";

const Header: FC<HeaderProps> = observer(() => {
    const { userStore } = useStore();
    const navigate = useNavigate();
    const logout = async () => {
        
    }
    return (
        <div className="w-full m-4">
            <div className={"flex items-center justify-end gap-4 mr-4"}>
                <Button styleType={ButtonStyleType.submit} onClick={() => navigate(PANEL_ROUTE)}>Панель Управления</Button>
                {userStore.IsAuth ? (
                    <div>
                        <Button styleType={ButtonStyleType.submit} onClick={async () => await UserService.logout()}>Выйти</Button>
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
})

export default Header;

