import { FC } from "react";
import { useNavigate } from "react-router-dom";
import { observer } from "mobx-react";

import { HeaderProps } from "./Header.props";
import Button from "@/ui/Button/Button";
import { ButtonStyleType } from "@/ui/Button/Button.props";
import { useStore } from "@/hooks/store";
import { LOGIN_ROUTE, PANEL_ROUTE, REGISTER_ROUTE } from "@/utils/consts";
import { UserService } from "@/http/User";

const Header: FC<HeaderProps> = observer(() => {
  const { userStore } = useStore();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await UserService.logout();
    userStore.clear();
  };

    return (
        <div className="w-full m-4">
            <div className={"flex items-center justify-end"}>
                {userStore.IsAuth ? (
                    <div className="flex gap-4 items-center w-1/2 justify-end mr-4">
                        <Button styleType={ButtonStyleType.submit} onClick={() => navigate(PANEL_ROUTE)}>Панель Управления</Button>
                        <Button styleType={ButtonStyleType.submit} onClick={handleLogout}>Выйти</Button>
                    </div>
                ) : (
                    <div className="flex gap-4 items-center w-1/2 justify-end mr-4">
                        <Button styleType={ButtonStyleType.submit} onClick={() => navigate(LOGIN_ROUTE)}>Войти</Button>
                        <Button styleType={ButtonStyleType.submit} onClick={() => navigate(REGISTER_ROUTE)}>Зарегистрироваться</Button>
                    </div>
                )}
            </div>
        </div>
    );
})

export default Header;