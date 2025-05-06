import React from "react";
import type {FC} from "react";
import {CardProps, styleMap} from "./Card.props";
import clsx from "clsx";

const Card: FC<CardProps> = ({style = "default", title, subtitle, children}) => {
    return (
        <div
            className={clsx(
                "w-full h-full rounded-2xl p-6 bg-base-lightBlue/75 text-white shadow-md transition duration-300 delay-100 ease-in hover:shadow-lg hover:shadow-base-darkBlue",
                styleMap[style]
            )}
        >
            {title && (
                <h2 className="text-lg font-semibold text-base-darkBlue mb-1 animate-fade-in">
                    {title}
                </h2>
            )}
            {subtitle && (
                <p className="text-sm text-base-grayBlue mb-4 animate-fade-in">
                    {subtitle}
                </p>
            )}
            <div className="w-full h-full animate-fade-in">{children}</div>
        </div>
    );
};

export default Card;

