import {FC, useEffect} from "react";
import React from "react";
import {HeaderProps} from "./Header.props";
import Button from "@/ui/Button/Button";
import {ButtonStyleType} from "@/ui/Button/Button.props";
import {useStore} from "@/hooks/store";

const Header: FC<HeaderProps> = () => {
    const { userStore } = useStore();
    useEffect(() => {
    }, [userStore]);
    return (
        <div className="flex items-center justify-end m-4">
            {userStore.IsAuth ? (
                <div>
                    <Button styleType={ButtonStyleType.submit}>Выйти</Button>
                </div>
            ) : (
                <div className="flex gap-4">
                    <Button styleType={ButtonStyleType.submit}>Войти</Button>
                    <Button styleType={ButtonStyleType.submit}>Зарегистрироваться</Button>
                </div>
            )}
        </div>
    );
};

export default Header;

