import { FC } from "react";
import { BoolInputProps } from "./BoolInput.props";
import clsx from "clsx";



export const BoolInput = ({
  value,
  onChange,
  disabled = false,
  helperText,
  error = false,
  Icon,
}: BoolInputProps) => {
  const toggle = () => {
    if (!disabled) {
      onChange(!value);
    }
  };

  return (
    <div className={"flex flex-col w-full gap-1"}>
      {helperText && (
        <p
          className={clsx(
            "text-xs ml-1",
            error ? "text-red-400" : "text-base-darkBlue"
          )}
        >
          {helperText}
        </p>
      )}
      <div
        className={clsx(
          "flex items-center px-4 py-2 rounded-xl transition-all duration-300 ease-in-out",
          "bg-base-darkBlue",
          {
            "border-red-400 animate-shake": error,
            "bg-base-dark": disabled,
            "cursor-not-allowed opacity-60": disabled,
            "focus-within:ring-2 focus-within:ring-base-lightBlue":
              !disabled && !error,
          }
        )}
      >
        {Icon && <Icon className={"text-base-grayBlue w-4 h-4 mr-2"} />}
        <button
          type="button"
          className={clsx(
            "relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none",
            value ? "bg-base-lightBlue" : "bg-base-grayBlue",
            disabled && "opacity-60 cursor-not-allowed"
          )}
          disabled={disabled}
          onClick={toggle}
        >
          <span
            className={clsx(
              "inline-block h-5 w-5 transform rounded-full bg-white transition-transform",
              value ? "translate-x-6" : "translate-x-0"
            )}
          />
        </button>
        <span className="ml-2 text-sm text-white">
          {value ? "Вкл" : "Выкл"}
        </span>
      </div>
    </div>
  );
};

export default BoolInput;